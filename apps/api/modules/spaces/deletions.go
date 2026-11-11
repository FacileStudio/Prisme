package spaces

import (
	"context"
	"encoding/json"
	"log/slog"
	"maps"
	"net/http"
	"slices"
	"strings"

	"github.com/FacileStudio/porte/directory"

	"gorm.io/gorm"
)

// Deletions receives the people Registre erased. It polls Registre's deletion
// list with the same client credentials as the snapshot sync and accepts
// Registre's signed push on directory.DeletionsPath; both hand each person to
// the eraser once it answers success.
type Deletions struct {
	receiver *directory.Deletions
}

// Eraser is the receiver's two callbacks: Prisme's erasure of the person,
// then porte's own rows eraser, pg.Store.EraseIdentity.
type Eraser struct {
	OnPersonDeleted func(ctx context.Context, provider, subject string) error
	Erase           func(ctx context.Context, provider, subject string) error
}

// NewDeletions builds the receiver, keeping its cursor in Prisme's database.
func NewDeletions(db *gorm.DB, client directory.Client, eraser Eraser, logger *slog.Logger) *Deletions {
	return &Deletions{receiver: &directory.Deletions{
		Client:          client,
		OnPersonDeleted: eraser.OnPersonDeleted,
		Erase:           eraser.Erase,
		State:           NewDeletionState(db, client.Issuer),
		Logger:          logger,
	}}
}

// Run polls until ctx ends. Start it in exactly one replica.
func (d *Deletions) Run(ctx context.Context) {
	d.receiver.Run(ctx)
}

// Handler verifies and applies Registre's signed user.deleted push.
func (d *Deletions) Handler() http.Handler {
	return d.receiver.Handler()
}

// deletionState keeps the receiver's cursor and handled subjects in one row.
type deletionState struct {
	db     *gorm.DB
	issuer string
}

// NewDeletionState is the receiver's state in Prisme's database. A subject
// it hands back as handled is one no identity holds any more: an account
// re-created for an erased subject makes the next notice or poll erase it
// instead of skipping it.
func NewDeletionState(db *gorm.DB, issuer string) directory.StateStore {
	return deletionState{db: db, issuer: issuer}
}

// Load returns the saved state, or nil before the first save, reopening the
// handled subjects an identity holds again.
func (s deletionState) Load(ctx context.Context) ([]byte, error) {
	var rows []struct{ State []byte }
	err := s.db.WithContext(ctx).Raw(`SELECT state FROM directory_deletion_state WHERE id = 1`).Scan(&rows).Error
	if err != nil || len(rows) == 0 {
		return nil, err
	}
	return s.reopen(ctx, rows[0].State)
}

// reopen drops from porte's handled subjects every one an identity of the
// issuer holds, keeping every other field of the state as porte wrote it. A
// state it cannot read goes back untouched, for porte to report.
func (s deletionState) reopen(ctx context.Context, raw []byte) ([]byte, error) {
	var state map[string]json.RawMessage
	var done map[string]json.RawMessage
	if json.Unmarshal(raw, &state) != nil || json.Unmarshal(state["done"], &done) != nil || len(done) == 0 {
		return raw, nil
	}
	held, err := s.held(ctx, slices.Collect(maps.Keys(done)))
	if err != nil || len(held) == 0 {
		return raw, err
	}
	for _, subject := range held {
		delete(done, subject)
	}
	if state["done"], err = json.Marshal(done); err != nil {
		return nil, err
	}
	return json.Marshal(state)
}

func (s deletionState) held(ctx context.Context, subjects []string) ([]string, error) {
	trimmed := strings.TrimSuffix(s.issuer, "/")
	var held []string
	err := s.db.WithContext(ctx).Table("porte_identities").
		Where("provider IN ? AND subject IN ?", []string{trimmed, trimmed + "/"}, subjects).
		Distinct().Pluck("subject", &held).Error
	return held, err
}

// Save replaces the saved state.
func (s deletionState) Save(ctx context.Context, state []byte) error {
	return s.db.WithContext(ctx).Exec(
		`INSERT INTO directory_deletion_state (id, state, updated_at) VALUES (1, ?, now())
		 ON CONFLICT (id) DO UPDATE SET state = EXCLUDED.state, updated_at = now()`, state).Error
}
