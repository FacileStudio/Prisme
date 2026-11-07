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
	if err := db.Exec(deletionState).Error; err != nil {
		return err
	}
	return db.Exec(pushDevices).Error
}

const deletionState = `CREATE TABLE IF NOT EXISTS directory_deletion_state (
	id         integer PRIMARY KEY,
	state      bytea NOT NULL,
	updated_at timestamptz NOT NULL DEFAULT now()
)`

// pushDevices holds the Firebase token of each signed-in Android app. A row
// belongs to a porte session, so signing out deletes it.
const pushDevices = `CREATE TABLE IF NOT EXISTS push_devices (
	token            text PRIMARY KEY,
	user_id          bigint NOT NULL REFERENCES porte_users(id) ON DELETE CASCADE,
	porte_session_id bigint NOT NULL REFERENCES porte_sessions(id) ON DELETE CASCADE,
	created_at       timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS push_devices_user_idx ON push_devices (user_id)`
