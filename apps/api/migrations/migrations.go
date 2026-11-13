// Package migrations embeds the ordered SQL that owns Prisme's schema and
// applies it with tronc/migrate.
//
// Until v0.18.0 the schema came from GORM AutoMigrate, porte/pg's Schema and
// hand-written CREATE TABLE IF NOT EXISTS at every boot. 00001_baseline.sql is
// the production schema those built, dumped with pg_dump; every later file is
// an ordinary migration.
package migrations

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"hash/crc64"
	"log/slog"

	"github.com/FacileStudio/tronc/migrate"
)

// FS holds the migration files at its root, which is the layout tronc/migrate
// expects.
//
//go:embed *.sql
var FS embed.FS

// BaselineVersion is the migration that reproduces the schema AutoMigrate built.
const BaselineVersion = 1

// Config is the tronc/migrate configuration over Prisme's migrations.
func Config(db *sql.DB, logger *slog.Logger) migrate.Config {
	return migrate.Config{DB: db, FS: FS, Logger: logger}
}

// SchemaConfig is Config under an advisory lock of the current schema's own.
// Advisory locks are scoped to the database, so test packages that each
// migrate a schema of their own in one test database would otherwise queue on
// tronc/migrate's default lock, five seconds per retry. Production keeps
// Config, so the boot and the migrate subcommand share one lock.
func SchemaConfig(ctx context.Context, db *sql.DB, logger *slog.Logger) (migrate.Config, error) {
	var schema string
	if err := db.QueryRowContext(ctx, `SELECT current_schema()`).Scan(&schema); err != nil {
		return migrate.Config{}, fmt.Errorf("migrations: read the schema: %w", err)
	}
	cfg := Config(db, logger)
	cfg.LockID = int64(crc64.Checksum([]byte(schema), crc64.MakeTable(crc64.ISO)))
	return cfg, nil
}

// Up applies every pending migration. A database built by AutoMigrate, with
// Prisme's tables and no applied migration, is first stamped at the baseline
// through tronc/migrate's own baseline command, so the baseline is recorded
// as applied without running and only the later migrations execute.
func Up(ctx context.Context, cfg migrate.Config) error {
	legacy, err := builtByAutoMigrate(ctx, cfg.DB)
	if err != nil {
		return err
	}
	if legacy {
		if cfg.Logger == nil {
			cfg.Logger = slog.Default()
		}
		cfg.Logger.Info("database predates migrations, recording the baseline", slog.Int("version", BaselineVersion))
		args := []string{"prisme", "migrate", "baseline", fmt.Sprint(BaselineVersion)}
		if _, err := migrate.Command(ctx, args, cfg); err != nil {
			return err
		}
	}
	return migrate.Run(ctx, cfg)
}

// builtByAutoMigrate reports whether the current schema holds Prisme's
// accounts table but no applied migration. A ledger can exist without one:
// goose creates it, with only its version 0 sentinel, the first time anything
// reads it, `migrate status` included. Every name resolves through
// search_path, so a test schema is judged on its own tables.
func builtByAutoMigrate(ctx context.Context, db *sql.DB) (bool, error) {
	var hasAccounts, hasLedger bool
	err := db.QueryRowContext(ctx,
		`SELECT to_regclass('accounts') IS NOT NULL, to_regclass($1) IS NOT NULL`, migrate.TableName,
	).Scan(&hasAccounts, &hasLedger)
	if err != nil {
		return false, fmt.Errorf("migrations: inspect schema: %w", err)
	}
	if !hasAccounts || !hasLedger {
		return hasAccounts, nil
	}
	var applied bool
	err = db.QueryRowContext(ctx, `SELECT EXISTS (SELECT 1 FROM `+migrate.TableName+` WHERE version_id > 0)`).Scan(&applied)
	if err != nil {
		return false, fmt.Errorf("migrations: read the ledger: %w", err)
	}
	return !applied, nil
}
