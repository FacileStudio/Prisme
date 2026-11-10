package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/FacileStudio/Prisme/apps/api/internal/database"
	"github.com/FacileStudio/Prisme/apps/api/internal/env"
	"github.com/FacileStudio/Prisme/apps/api/internal/erasure"
	"github.com/FacileStudio/Prisme/apps/api/internal/matrix"
)

// eraseUser runs `api erase-user <email>`: the operator's path for an
// account that holds no Registre subject, which no deletion notice can reach
// because a notice carries only the subject. It erases the account exactly
// as a notice would, Synapse first, and refuses a linked account.
func eraseUser(args []string, out io.Writer) int {
	if len(args) != 1 {
		_, _ = fmt.Fprintln(out, "usage: api erase-user <email>")
		return 2
	}
	eraser, err := openEraser()
	if err != nil {
		_, _ = fmt.Fprintln(out, "erase-user:", err)
		return 1
	}
	found, err := eraser.EraseUnlinked(context.Background(), args[0])
	switch {
	case err != nil:
		_, _ = fmt.Fprintln(out, "erase-user:", err)
		return 1
	case !found:
		_, _ = fmt.Fprintln(out, "erase-user: no account holds this email")
		return 1
	}
	_, _ = fmt.Fprintln(out, "erase-user: erased")
	return 0
}

func openEraser() (*erasure.Eraser, error) {
	cfg, err := env.Load()
	if err != nil {
		return nil, err
	}
	db, err := database.Open(cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	client := matrix.New(matrix.Config{BaseURL: cfg.Matrix.InternalURL, ServerName: cfg.Matrix.ServerName, ASToken: cfg.Matrix.ASToken})
	return newEraser(cfg, db, sqlDB, client, slog.New(slog.NewTextHandler(os.Stderr, nil)))
}
