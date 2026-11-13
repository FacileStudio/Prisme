package schemas

import (
	"context"
	"log/slog"

	"github.com/FacileStudio/Prisme/apps/api/migrations"
	"gorm.io/gorm"
)

// Migrate applies the ordered migrations under migrations/ under a lock of
// the current schema's own. The tests build their schemas through it, so they
// run exactly the SQL production runs.
func Migrate(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	ctx := context.Background()
	cfg, err := migrations.SchemaConfig(ctx, sqlDB, nil)
	if err != nil {
		return err
	}
	return migrations.Up(ctx, cfg)
}

// MigrateWithLogger is what the boot runs: the migrations under tronc's
// default lock, logged through the app's logger. Prisme runs no data step
// after them: the two the old boot ran had nothing left to change, see the
// 0.18.0 entry of CHANGELOG.md.
func MigrateWithLogger(ctx context.Context, db *gorm.DB, logger *slog.Logger) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return migrations.Up(ctx, migrations.Config(sqlDB, logger))
}
