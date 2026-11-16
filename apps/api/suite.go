package main

import (
	"log/slog"

	"github.com/FacileStudio/Prisme/apps/api/internal/accounts"
	"github.com/FacileStudio/Prisme/apps/api/internal/capsule"
	"github.com/FacileStudio/Prisme/apps/api/internal/echo"
	"github.com/FacileStudio/Prisme/apps/api/internal/env"
	giphy "github.com/FacileStudio/Prisme/apps/api/internal/gifs"
	"github.com/FacileStudio/Prisme/apps/api/internal/matrix"
	linkpreview "github.com/FacileStudio/Prisme/apps/api/internal/unfurl"
	"github.com/FacileStudio/Prisme/apps/api/modules/calls"
	"github.com/FacileStudio/Prisme/apps/api/modules/capsules"
	"github.com/FacileStudio/Prisme/apps/api/modules/gifs"
	"github.com/FacileStudio/Prisme/apps/api/modules/unfurl"

	"gorm.io/gorm"
)

func newCalls(cfg env.Config, db *gorm.DB, client *matrix.Client, accountService *accounts.Service, appLogger *slog.Logger) *calls.Service {
	return calls.NewService(calls.Deps{
		DB: db, Echo: newEcho(cfg), Matrix: client, Accounts: accountService,
		Secret: cfg.SecretKey, PublicURL: cfg.PublicURL, Logger: appLogger,
	})
}

func newSuiteApps(cfg env.Config, db *gorm.DB, client *matrix.Client, accountService *accounts.Service, appLogger *slog.Logger) suiteApps {
	return suiteApps{
		calls:    newCalls(cfg, db, client, accountService, appLogger),
		capsules: capsules.NewService(capsule.New(cfg.CapsuleURL), appLogger),
		unfurl:   unfurl.NewService(linkpreview.New(linkpreview.SuiteHosts())),
		gifs:     gifs.NewService(newGiphy(cfg), appLogger),
	}
}

func newEcho(cfg env.Config) *echo.Client {
	if !cfg.CallsEnabled() {
		return nil
	}
	return echo.New(cfg.Echo.URL, cfg.Echo.Key)
}

func newGiphy(cfg env.Config) *giphy.Client {
	if cfg.GiphyKey == "" {
		return nil
	}
	return giphy.New(cfg.GiphyKey)
}
