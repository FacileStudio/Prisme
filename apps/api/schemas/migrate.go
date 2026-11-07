package schemas

import (
	"github.com/FacileStudio/porte/pg"
	"gorm.io/gorm"
)

// Migrate creates porte's identity tables, then Prisme's own. Both are
// idempotent, and Prisme runs a single replica, so boot is the right place.
func Migrate(db *gorm.DB) error {
	if err := db.Exec(pg.Schema).Error; err != nil {
		return err
	}
	if err := db.AutoMigrate(&Account{}, &MatrixSession{}, &Space{}, &Channel{}, &Agent{}, &Call{}); err != nil {
		return err
	}
	return db.Exec(deletionState).Error
}

const deletionState = `CREATE TABLE IF NOT EXISTS directory_deletion_state (
	id         integer PRIMARY KEY,
	state      bytea NOT NULL,
	updated_at timestamptz NOT NULL DEFAULT now()
)`
