package spaces

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/FacileStudio/Prisme/apps/api/internal/matrix"
	"github.com/FacileStudio/porte/directory"

	"gorm.io/gorm"
)

// Deletions erases the people Registre erased. It polls Registre's deletion
// list with the same client credentials as the snapshot sync and accepts
// Registre's signed push on directory.DeletionsPath.
type Deletions struct {
	db       *gorm.DB
	matrix   *matrix.Client
	receiver *directory.Deletions
}

// NewDeletions builds the receiver over Prisme's tables. erase is porte's own
// rows eraser, pg.Store.EraseIdentity.
func NewDeletions(
	service *Service, client directory.Client, erase func(context.Context, string, string) error, logger *slog.Logger,
) *Deletions {
	d := &Deletions{db: service.db, matrix: service.matrix}
	d.receiver = &directory.Deletions{
		Client:          client,
		OnPersonDeleted: d.ErasePerson,
		Erase:           erase,
		State:           deletionState{db: service.db},
		Logger:          logger,
	}
	return d
}

// Run polls until ctx ends. Start it in exactly one replica.
func (d *Deletions) Run(ctx context.Context) {
	d.receiver.Run(ctx)
}

// Handler verifies and applies Registre's signed user.deleted push.
func (d *Deletions) Handler() http.Handler {
	return d.receiver.Handler()
}

var erasureStatements = []string{
	`UPDATE calls SET started_by = '' WHERE started_by IN @matrix_ids`,
	`UPDATE spaces SET created_by = 0 WHERE created_by = @id`,
	`DELETE FROM matrix_sessions WHERE user_id = @id`,
	`UPDATE agents SET owner_id = 0, name = '', sealed_token = '', device_id = '', avatar_url = '',
		avatar_source = '', deactivated_at = coalesce(deactivated_at, now()) WHERE owner_id = @id`,
	`UPDATE accounts SET user_id = -user_id, avatar_url = '' WHERE user_id = @id`,
	`DELETE FROM porte_users WHERE id = @id`,
}

// ErasePerson erases the Matrix accounts of a Registre identity and what Prisme
// stores for them: the account and its agents are deactivated and erased in
// Synapse, the sessions and the porte user are deleted, and the agents lose
// their owner, name and tokens. Matrix never reissues an id, so the account row
// and the agent rows stay as reservations of the deactivated ids: the account
// under a negated user id, the agents deactivated. Spaces and channels the
// person created stay, with no creator, and their messages stay in the rooms
// under a deactivated account. A subject that is already gone is a no-op, so a
// replayed notice changes nothing.
func (d *Deletions) ErasePerson(ctx context.Context, provider, subject string) error {
	var owners []int64
	err := d.db.WithContext(ctx).Table("porte_identities").
		Where("provider = ? AND subject = ?", provider, subject).Pluck("user_id", &owners).Error
	if err != nil || len(owners) == 0 {
		return err
	}
	userID := owners[0]
	matrixIDs, err := d.matrixIDs(ctx, userID)
	if err != nil {
		return err
	}
	for _, id := range matrixIDs {
		if err := d.matrix.Erase(ctx, id); err != nil && !matrix.IsCode(err, "M_USER_DEACTIVATED") {
			return err
		}
	}
	return d.deleteRows(ctx, userID, matrixIDs)
}

func (d *Deletions) matrixIDs(ctx context.Context, userID int64) ([]string, error) {
	var ids []string
	db := d.db.WithContext(ctx)
	if err := db.Table("accounts").Where("user_id = ?", userID).Pluck("matrix_user_id", &ids).Error; err != nil {
		return nil, err
	}
	var agents []string
	err := db.Table("agents").Where("owner_id = ? AND deactivated_at IS NULL", userID).Pluck("matrix_user_id", &agents).Error
	return append(ids, agents...), err
}

func (d *Deletions) deleteRows(ctx context.Context, userID int64, matrixIDs []string) error {
	args := map[string]any{"id": userID, "matrix_ids": append(matrixIDs, "")}
	return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, statement := range erasureStatements {
			if err := tx.Exec(statement, args).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// deletionState keeps the receiver's cursor and handled subjects in one row.
type deletionState struct {
	db *gorm.DB
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
