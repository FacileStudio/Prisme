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
	"github.com/FacileStudio/Prisme/apps/api/modules/calls"
	"github.com/FacileStudio/Prisme/apps/api/modules/push"
	"github.com/FacileStudio/Prisme/apps/api/modules/spaces"

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

// withBackground adds what runs beside the routes: the Android push relay, the
// eraser and the Registre directory sync that hands it deletion notices.
func withBackground(built app, cfg env.Config, db *gorm.DB, sqlDB *sql.DB, client *matrix.Client) (app, error) {
	built.push = newPush(cfg, db, client, built.logger)
	eraser, err := newEraser(cfg, db, sqlDB, client, built.logger)
	if err != nil {
		return app{}, err
	}
	built.eraser = eraser
	built.directory = newDirectory(cfg, built, db, sqlDB)
	return built, nil
}
