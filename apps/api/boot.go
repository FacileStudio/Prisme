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
	"github.com/FacileStudio/Prisme/apps/api/internal/env"
	"github.com/FacileStudio/Prisme/apps/api/internal/matrix"
	"github.com/FacileStudio/Prisme/apps/api/internal/seal"
	"github.com/FacileStudio/Prisme/apps/api/modules/agents"
	"github.com/FacileStudio/Prisme/apps/api/modules/auth"
	"github.com/FacileStudio/Prisme/apps/api/modules/session"
	"github.com/FacileStudio/Prisme/apps/api/modules/spaces"
	"github.com/FacileStudio/Prisme/apps/api/modules/users"
	"github.com/FacileStudio/Prisme/apps/api/schemas"

	"github.com/FacileStudio/Journal/sdk/journal"
	"github.com/FacileStudio/porte/directory"
	"github.com/FacileStudio/porte/pg"
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
	return start(services, ":"+strconv.Itoa(cfg.Port))
}

// start runs the Registre sync and the deletion poll, when configured, beside
// the HTTP server until a signal ends all of them. Prisme runs one replica, so
// this is the one deletion poller.
func start(services app, addr string) int {
	shutdownSignal, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	if sync := services.directory; sync.registre != nil {
		go sync.registre.Run(shutdownSignal, sync.syncer)
		go sync.deletions.Run(shutdownSignal)
	}
	return serve(shutdownSignal, addr, buildRouter(services), services.logger)
}

// newDirectory builds the Registre space sync and the GDPR deletion receiver
// when OIDC is configured, with Prisme's own client credentials, and nils
// otherwise: without Registre there is nothing to sync and no route to mount.
func newDirectory(cfg env.Config, spaceService *spaces.Service, sqlDB *sql.DB, appLogger *slog.Logger) directorySync {
	porteCfg := cfg.Porte
	if porteCfg.Issuer == "" || porteCfg.ClientID == "" || porteCfg.ClientSecret == "" {
		return directorySync{}
	}
	registre := spaces.NewRegistre(spaceService, porteCfg.Issuer, appLogger)
	client := directory.Client{Issuer: porteCfg.Issuer, ClientID: porteCfg.ClientID, ClientSecret: porteCfg.ClientSecret}
	syncer := &directory.Syncer{Client: client, Apply: registre.Apply, Logger: appLogger}
	deletions := spaces.NewDeletions(spaceService, client, pg.New(sqlDB).EraseIdentity, appLogger)
	return directorySync{registre: registre, syncer: syncer, deletions: deletions}
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
	agentService := agents.NewService(agents.Deps{
		DB: db, Matrix: client, Accounts: accountService, Box: box, Homeserver: cfg.Matrix.PublicURL, Logger: appLogger,
	})
	built := app{
		logger: appLogger, dbCheck: health.DB(sqlDB),
		edge: edge{core: cfg.Core, homeserver: cfg.Matrix.PublicURL, liveKit: cfg.LiveKitURL},
		auth: authService,
		session: session.NewService(session.Deps{
			DB: db, Matrix: client, Accounts: accountService, Box: box,
			Homeserver: cfg.Matrix.PublicURL, CallsEnabled: cfg.CallsEnabled(),
		}),
		spaces: spaces.NewService(spaces.Deps{
			DB: db, Matrix: client, Accounts: accountService, Followers: agentService, Logger: appLogger,
		}),
		agents: agentService, calls: newCalls(cfg, db, client, accountService, appLogger),
		users: users.NewService(db),
	}
	built.push = newPush(cfg, db, client, appLogger)
	built.directory = newDirectory(cfg, built.spaces, sqlDB, appLogger)
	return built, nil
}

func serve(shutdownSignal context.Context, addr string, router http.Handler, appLogger *slog.Logger) int {
	server := &http.Server{
		Addr: addr, Handler: router,
		ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 30 * time.Second,
		WriteTimeout: 60 * time.Second, IdleTimeout: 60 * time.Second,
	}
	serverErr := make(chan error, 1)
	go func() { serverErr <- server.ListenAndServe() }()
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
