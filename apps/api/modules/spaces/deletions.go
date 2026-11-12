package spaces

import (
	"context"
	"log/slog"
	"net/http"

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

// Eraser is the receiver's callbacks: Prisme's erasure of the person, then
// porte's own rows eraser, pg.Store.EraseIdentity, and pg.Store.HoldsIdentity,
// which makes a re-pushed notice erase an account re-created for a handled
// subject instead of skipping it.
type Eraser struct {
	OnPersonDeleted func(ctx context.Context, provider, subject string) error
	Erase           func(ctx context.Context, provider, subject string) error
	Holds           func(ctx context.Context, provider, subject string) (bool, error)
}

// NewDeletions builds the receiver, keeping its cursor in Prisme's database.
func NewDeletions(db *gorm.DB, client directory.Client, eraser Eraser, logger *slog.Logger) *Deletions {
	return &Deletions{receiver: &directory.Deletions{
		Client:          client,
		OnPersonDeleted: eraser.OnPersonDeleted,
		Erase:           eraser.Erase,
		Holds:           eraser.Holds,
		State:           NewDeletionState(db),
		Logger:          logger,
	}}
}

// Run polls until ctx ends. Start it in exactly one replica.
func (d *Deletions) Run(ctx context.Context) {
	d.receiver.Run(ctx)
}

// Erased reports whether the receiver handled this subject of the issuer:
// porte/oidc's Deps.Erased, which then refuses it an account. It reads the
// saved state on every call.
func (d *Deletions) Erased(ctx context.Context, provider, subject string) (bool, error) {
	return d.receiver.Erased(ctx, provider, subject)
}

// Handler verifies and applies Registre's signed user.deleted push.
func (d *Deletions) Handler() http.Handler {
	return d.receiver.Handler()
}

// deletionState keeps the receiver's cursor and handled subjects in one row.
type deletionState struct {
	db *gorm.DB
}

// NewDeletionState is the receiver's state in Prisme's database.
func NewDeletionState(db *gorm.DB) directory.StateStore {
	return deletionState{db: db}
}

// Load returns the saved state, or nil before the first save.
func (s deletionState) Load(ctx context.Context) ([]byte, error) {
	var rows []struct{ State []byte }
	err := s.db.WithContext(ctx).Raw(`SELECT state FROM directory_deletion_state WHERE id = 1`).Scan(&rows).Error
	if err != nil || len(rows) == 0 {
		return nil, err
	}
	return rows[0].State, nil
}

// Save replaces the saved state.
func (s deletionState) Save(ctx context.Context, state []byte) error {
	return s.db.WithContext(ctx).Exec(
		`INSERT INTO directory_deletion_state (id, state, updated_at) VALUES (1, ?, now())
		 ON CONFLICT (id) DO UPDATE SET state = EXCLUDED.state, updated_at = now()`, state).Error
}
