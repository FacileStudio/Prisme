package main

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/FacileStudio/Prisme/apps/api/internal/accounts"
	"github.com/FacileStudio/Prisme/apps/api/internal/echo"
	"github.com/FacileStudio/Prisme/apps/api/internal/env"
	"github.com/FacileStudio/Prisme/apps/api/internal/erasure"
	"github.com/FacileStudio/Prisme/apps/api/internal/matrix"
	"github.com/FacileStudio/Prisme/apps/api/internal/seal"
	"github.com/FacileStudio/Prisme/apps/api/modules/auth"
	"github.com/FacileStudio/Prisme/apps/api/modules/calls"
	"github.com/FacileStudio/Prisme/apps/api/modules/push"
	"github.com/FacileStudio/Prisme/apps/api/modules/session"
	"github.com/FacileStudio/Prisme/apps/api/modules/spaces"

	"github.com/FacileStudio/porte"
	"github.com/FacileStudio/porte/pg"

	"gorm.io/gorm"
)

// newPush builds the Android push relay. Unreadable Firebase credentials are
// logged and leave push off rather than keeping the chat from starting.
func newPush(cfg env.Config, db *gorm.DB, client *matrix.Client, appLogger *slog.Logger) *push.Service {
	fcm, err := push.NewFCM(cfg.FCMAccount)
	if err != nil {
		appLogger.Error("android push is off: FCM_SERVICE_ACCOUNT is unreadable", slog.Any("error", err))
	}
	return push.NewService(push.Deps{
		DB: db, Matrix: client, FCM: fcm, HSToken: cfg.Matrix.HSToken, Logger: appLogger,
	})
}

func newCalls(cfg env.Config, db *gorm.DB, client *matrix.Client, accountService *accounts.Service, appLogger *slog.Logger) *calls.Service {
	return calls.NewService(calls.Deps{
		DB: db, Echo: newEcho(cfg), Matrix: client, Accounts: accountService,
		Secret: cfg.SecretKey, PublicURL: cfg.PublicURL, Logger: appLogger,
	})
}

func newEcho(cfg env.Config) *echo.Client {
	if !cfg.CallsEnabled() {
		return nil
	}
	return echo.New(cfg.Echo.URL, cfg.Echo.Key)
}

// newEraser builds the GDPR eraser the deletion receiver and the erase-user
// command share.
func newEraser(cfg env.Config, db *gorm.DB, sqlDB *sql.DB, client *matrix.Client, logger *slog.Logger) (*erasure.Eraser, error) {
	return erasure.New(context.Background(), erasure.Deps{
		DB: db, SQL: sqlDB, Matrix: client, Issuer: cfg.Porte.Issuer,
		Bot: client.UserID(spaces.BotLocalpart), Logger: logger,
	})
}

// newDeletions builds the GDPR deletion receiver when OIDC is configured, with
// Prisme's own client credentials, and nil otherwise. It exists before porte's
// OIDC kit, which refuses an account to every subject it handled.
func newDeletions(cfg env.Config, db *gorm.DB, sqlDB *sql.DB, eraser *erasure.Eraser, logger *slog.Logger) *spaces.Deletions {
	porteCfg := cfg.Porte
	if porteCfg.Issuer == "" || porteCfg.ClientID == "" || porteCfg.ClientSecret == "" {
		return nil
	}
	store := pg.New(sqlDB)
	return spaces.NewDeletions(db, directoryClient(porteCfg), spaces.Eraser{
		OnPersonDeleted: eraser.ErasePerson, Erase: store.EraseIdentity, Holds: store.HoldsIdentity,
	}, logger)
}

// erasedSubjects is the receiver as porte/oidc's guard, or no guard without one.
func erasedSubjects(deletions *spaces.Deletions) porte.ErasedSubjects {
	if deletions == nil {
		return nil
	}
	return deletions
}

// newAuth builds the deletion receiver, then porte with Bulle's service keys,
// refusing an account to every subject the receiver handled.
func newAuth(cfg env.Config, db *gorm.DB, sqlDB *sql.DB, eraser *erasure.Eraser, logger *slog.Logger) (*auth.Service, *spaces.Deletions, error) {
	deletions := newDeletions(cfg, db, sqlDB, eraser, logger)
	service, err := auth.NewService(context.Background(), sqlDB, cfg.Porte, auth.Guards{
		ServiceKeys: auth.ParseServiceKeys(cfg.ServiceKeys), Erased: erasedSubjects(deletions),
	}, logger)
	return service, deletions, err
}

func newSession(cfg env.Config, db *gorm.DB, client *matrix.Client, accountService *accounts.Service, box *seal.Box) *session.Service {
	return session.NewService(session.Deps{
		DB: db, Matrix: client, Accounts: accountService, Box: box,
		Homeserver: cfg.Matrix.PublicURL, CallsEnabled: cfg.CallsEnabled(),
	})
}
