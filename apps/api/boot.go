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
	"github.com/FacileStudio/Prisme/apps/api/modules/spaces"
	"github.com/FacileStudio/Prisme/apps/api/modules/users"
	"github.com/FacileStudio/Prisme/apps/api/schemas"

	"github.com/FacileStudio/Journal/sdk/journal"
	"github.com/FacileStudio/porte"
	"github.com/FacileStudio/porte/directory"
	"github.com/FacileStudio/porte/kit"
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
	go services.eraser.Report(shutdownSignal)
	if sync := services.directory; sync.registre != nil {
		go sync.registre.Run(shutdownSignal, sync.syncer)
		go sync.deletions.Run(shutdownSignal)
	}
	return serve(shutdownSignal, addr, buildRouter(services), services.logger)
}

// newDirectory builds the Registre space sync beside the deletion receiver,
// and nils without the receiver: without Registre there is nothing to sync
// and no route to mount.
func newDirectory(cfg env.Config, built app, deletions *spaces.Deletions) directorySync {
	if deletions == nil {
		return directorySync{}
	}
	registre := spaces.NewRegistre(built.spaces, cfg.Porte.Issuer, built.logger)
	client := directoryClient(cfg.Porte)
	syncer := &directory.Syncer{Client: client, Apply: registre.Apply, Logger: built.logger}
	invitations, err := kit.NewInvitations(kit.InvitationsConfig{
		Directory: client, Identities: built.auth.Identities(), Logger: built.logger,
	})
	if err != nil {
		built.logger.Warn("the invitations badge is off", slog.Any("error", err))
	}
	return directorySync{registre: registre, syncer: syncer, deletions: deletions, invitations: invitations}
}

func directoryClient(cfg porte.Config) directory.Client {
	return directory.Client{Issuer: cfg.Issuer, ClientID: cfg.ClientID, ClientSecret: cfg.ClientSecret}
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
	client := matrix.New(matrix.Config{BaseURL: cfg.Matrix.InternalURL, ServerName: cfg.Matrix.ServerName, ASToken: cfg.Matrix.ASToken})
	eraser, err := newEraser(cfg, db, sqlDB, client, appLogger)
	if err != nil {
		return app{}, err
	}
	authService, deletions, err := newAuth(cfg, db, sqlDB, eraser, appLogger)
	if err != nil {
		return app{}, err
	}
	box, err := seal.New(cfg.SecretKey)
	if err != nil {
		return app{}, err
	}
	accountService := accounts.NewService(db, client, appLogger)
	agentService := agents.NewService(agents.Deps{
		DB: db, Matrix: client, Accounts: accountService, Box: box, Homeserver: cfg.Matrix.PublicURL, Logger: appLogger,
	})
	built := app{
		logger: appLogger, dbCheck: health.DB(sqlDB),
		edge: edge{core: cfg.Core, homeserver: cfg.Matrix.PublicURL, liveKit: cfg.LiveKitURL},
		auth: authService, eraser: eraser,
		session: newSession(cfg, db, client, accountService, box),
		spaces: spaces.NewService(spaces.Deps{
			DB: db, Matrix: client, Accounts: accountService, Followers: agentService, Logger: appLogger,
		}),
		agents: agentService, calls: newCalls(cfg, db, client, accountService, appLogger),
		users: users.NewService(db), push: newPush(cfg, db, client, appLogger),
	}
	built.directory = newDirectory(cfg, built, deletions)
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
