// Package erasure deletes a person Registre erased, from Synapse and from
// Prisme's tables, so the two never disagree.
//
// The Synapse half runs first and the database half last, in one
// transaction. A failure in either returns an error, the deletion notice is
// answered with it, and Registre or the next poll retries. Every step is
// idempotent, so a retry after a failure anywhere converges: the Synapse
// calls repeat harmlessly and the rows still name the person until the final
// transaction commits. docs/erasure.md lists what remains afterwards.
package erasure

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"time"

	"github.com/FacileStudio/Prisme/apps/api/internal/matrix"
	"github.com/FacileStudio/porte/kit"

	"gorm.io/gorm"
)

const (
	adminAttempts = 12
	adminRetry    = 10 * time.Second
)

// Deps is what an Eraser needs. Issuer is the Registre issuer identities are
// filed under, empty when SSO is off. Bot is the appservice sender's id.
type Deps struct {
	DB     *gorm.DB
	SQL    *sql.DB
	Matrix *matrix.Client
	Issuer string
	Bot    string
	Logger *slog.Logger
}

// Eraser erases one person at a time.
type Eraser struct {
	deps  Deps
	store *kit.Store
}

// New builds the eraser over porte_users through porte/kit, whose erase hook
// resolves the identity under either issuer spelling and deletes the user row
// after Prisme's own rows, in one transaction.
func New(ctx context.Context, deps Deps) (*Eraser, error) {
	store, err := kit.NewStore(ctx, kit.Store{
		DB: deps.SQL, Columns: kit.Columns{Table: "porte_users"},
		OnCreate: refuseCreate, EraseData: eraseRows, Logger: deps.Logger,
	})
	if err != nil {
		return nil, err
	}
	return &Eraser{deps: deps, store: store}, nil
}

// refuseCreate keeps this store to erasing: porte/pg's UserStore owns sign-in.
func refuseCreate(context.Context, *sql.Tx, int64) error {
	return errors.New("erasure: this store never creates accounts")
}

// ErasePerson is porte/directory's OnPersonDeleted. It answers nil only when
// the person was found and erased, or provably never held an account here:
// Registre counts any success as delivered and never asks again.
func (e *Eraser) ErasePerson(ctx context.Context, provider, subject string) error {
	userID, err := e.owner(ctx, provider, subject)
	if err != nil {
		return err
	}
	if userID == 0 {
		return e.proveUnknown(ctx, provider, subject)
	}
	if err := e.eraseMatrix(ctx, userID); err != nil {
		return err
	}
	return e.store.EraseAccount(ctx, provider, subject)
}

// Report logs, once at boot, whether the admin credential works and how many
// accounts carry no Registre subject. It names no one. Synapse may still be
// starting, so the admin check retries for about two minutes before it
// reports a failure.
func (e *Eraser) Report(ctx context.Context) {
	logger := e.deps.Logger
	if err := e.waitAdmin(ctx); err != nil {
		logger.Error("erasure: the synapse admin api refused the appservice token; erasures will fail until it works",
			slog.Any("error", err))
	} else {
		logger.Info("erasure: the synapse admin api accepts the appservice token")
	}
	if e.deps.Issuer == "" {
		return
	}
	if unlinked, err := e.unlinked(ctx, e.deps.Issuer); err == nil && unlinked > 0 {
		logger.Warn("erasure: accounts carry no registre subject; a deletion notice for an unknown subject fails until they do",
			slog.Int64("accounts", unlinked))
	}
}

func (e *Eraser) waitAdmin(ctx context.Context) error {
	var err error
	for range adminAttempts {
		if err = e.deps.Matrix.AdminReady(ctx); err == nil {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(adminRetry):
		}
	}
	return err
}
