package main

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/FacileStudio/Prisme/apps/api/internal/accounts"
	"github.com/FacileStudio/Prisme/apps/api/internal/database"
	"github.com/FacileStudio/Prisme/apps/api/internal/echo"
	"github.com/FacileStudio/Prisme/apps/api/internal/env"
	"github.com/FacileStudio/Prisme/apps/api/internal/matrix"
	"github.com/FacileStudio/Prisme/apps/api/internal/seal"
	"github.com/FacileStudio/Prisme/apps/api/modules/agents"
	"github.com/FacileStudio/Prisme/apps/api/modules/auth"
	"github.com/FacileStudio/Prisme/apps/api/modules/calls"
	"github.com/FacileStudio/Prisme/apps/api/modules/session"
	"github.com/FacileStudio/Prisme/apps/api/modules/spaces"
	"github.com/FacileStudio/Prisme/apps/api/modules/users"
	"github.com/FacileStudio/Prisme/apps/api/schemas"

	"github.com/FacileStudio/Journal/sdk/journal"
	"github.com/FacileStudio/tronc/health"
	"github.com/FacileStudio/tronc/logger"
	"gorm.io/gorm"
)

func run() int {
	cfg, err := env.Load()
	bootLog := logger.New(logger.Config{})
	if err != nil {
		bootLog.Error("failed to load config", slog.Any("error", err))
		return 1
	}
	appLogger, flush := newLogger(cfg)
	defer flush()
	db, err := database.Open(cfg.DatabaseURL)
	if err == nil {
		err = schemas.Migrate(db)
	}
	if err != nil {
		appLogger.Error("failed to prepare the database", slog.Any("error", err))
		return 1
	}
	sqlDB, err := db.DB()
	if err != nil {
		appLogger.Error("failed to access the database handle", slog.Any("error", err))
		return 1
	}
	defer func() { _ = sqlDB.Close() }()
	services, err := buildApp(cfg, db, sqlDB, appLogger)
	if err != nil {
		appLogger.Error("failed to build the services", slog.Any("error", err))
		return 1
	}
	return serve(":"+strconv.Itoa(cfg.Port), buildRouter(services), appLogger)
}

func newLogger(cfg env.Config) (*slog.Logger, func()) {
	var client *journal.Client
	appLogger := logger.New(logger.Config{
		Level: cfg.LogLevel,
		Wrap: func(handler slog.Handler) slog.Handler {
			if cfg.JournalURL == "" || cfg.JournalToken == "" {
				return handler
			}
			client = journal.New(journal.Config{URL: cfg.JournalURL, Token: cfg.JournalToken})
			return journal.NewHandler(client, handler)
		},
	})
	return appLogger, func() {
		if client != nil {
			client.Close()
		}
	}
}

func buildApp(cfg env.Config, db *gorm.DB, sqlDB *sql.DB, appLogger *slog.Logger) (app, error) {
	authService, err := auth.NewService(context.Background(), sqlDB, cfg.Porte, auth.ParseServiceKeys(cfg.ServiceKeys), appLogger)
	if err != nil {
		return app{}, err
	}
	box, err := seal.New(cfg.SecretKey)
	if err != nil {
		return app{}, err
	}
	client := matrix.New(matrix.Config{BaseURL: cfg.Matrix.InternalURL, ServerName: cfg.Matrix.ServerName, ASToken: cfg.Matrix.ASToken})
	accountService := accounts.NewService(db, client, appLogger)
	var echoClient *echo.Client
	if cfg.CallsEnabled() {
		echoClient = echo.New(cfg.Echo.URL, cfg.Echo.Key)
	}
	return app{
		logger: appLogger, cors: cfg.CORSAllowedOrigins, dbCheck: health.DB(sqlDB), homeserver: cfg.Matrix.PublicURL,
		auth: authService,
		session: session.NewService(session.Deps{
			DB: db, Matrix: client, Accounts: accountService, Box: box,
			Homeserver: cfg.Matrix.PublicURL, CallsEnabled: cfg.CallsEnabled(),
		}),
		spaces: spaces.NewService(db, client, accountService),
		agents: agents.NewService(agents.Deps{DB: db, Matrix: client, Accounts: accountService, Box: box, Homeserver: cfg.Matrix.PublicURL, Logger: appLogger}),
		calls:  calls.NewService(echoClient, client, accountService),
		users:  users.NewService(db),
	}, nil
}

func serve(addr string, router http.Handler, appLogger *slog.Logger) int {
	server := &http.Server{
		Addr: addr, Handler: router,
		ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 30 * time.Second,
		WriteTimeout: 60 * time.Second, IdleTimeout: 60 * time.Second,
	}
	serverErr := make(chan error, 1)
	go func() { serverErr <- server.ListenAndServe() }()
	shutdownSignal, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	appLogger.Info("server starting", slog.String("addr", addr))
	select {
	case err := <-serverErr:
		if !errors.Is(err, http.ErrServerClosed) {
			appLogger.Error("server stopped", slog.Any("error", err))
			return 1
		}
		return 0
	case <-shutdownSignal.Done():
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			return 1
		}
		return 0
	}
}
