package main

import (
	"log/slog"
	"net/http"
	"os"

	documentation "github.com/FacileStudio/Prisme/apps/api/internal/documentation"
	"github.com/FacileStudio/Prisme/apps/api/modules/agents"
	"github.com/FacileStudio/Prisme/apps/api/modules/auth"
	"github.com/FacileStudio/Prisme/apps/api/modules/calls"
	"github.com/FacileStudio/Prisme/apps/api/modules/session"
	"github.com/FacileStudio/Prisme/apps/api/modules/spaces"
	"github.com/FacileStudio/Prisme/apps/api/modules/users"

	"github.com/FacileStudio/porte/directory"
	"github.com/FacileStudio/tronc/apiref"
	troncenv "github.com/FacileStudio/tronc/env"
	"github.com/FacileStudio/tronc/health"
	"github.com/FacileStudio/tronc/healthcheck"
	"github.com/FacileStudio/tronc/httpjson"
	"github.com/FacileStudio/tronc/httpx"
	troncmiddleware "github.com/FacileStudio/tronc/middleware"
	"github.com/FacileStudio/tronc/spa"
	"github.com/go-chi/chi/v5"
)

// apiTitle names this API in the reference page and the OpenAPI document.
// scripts/rename.sh rewrites this line when a fork is retargeted, so keep it a
// single const on one line.
const apiTitle = "Prisme API"

const apiDescription = "Team chat on Matrix: sessions, spaces, agents and calls, provisioned through a Synapse appservice."

func registry() documentation.Response {
	return documentation.Response{
		Modules: []documentation.Module{
			session.Documentation(),
			spaces.Documentation(),
			agents.Documentation(),
			calls.Documentation(),
			users.Documentation(),
		},
	}
}

func refConfig() apiref.Config {
	return apiref.Config{Title: apiTitle, Description: apiDescription, Servers: []string{"/api"}, Registry: registry()}
}

// app is every service the router mounts. A nil auth leaves porte and the
// modules unmounted, which is what the route-table test walks.
type app struct {
	logger     *slog.Logger
	cors       []string
	dbCheck    health.Check
	homeserver string
	auth       *auth.Service
	session    *session.Service
	spaces     *spaces.Service
	agents     *agents.Service
	calls      *calls.Service
	users      *users.Service
	registre   *spaces.Registre
	syncer     *directory.Syncer
	deletions  *spaces.Deletions
	core       troncenv.Core
	liveKit    string
}

func buildRouter(a app) chi.Router {
	clientDir := spa.DirFromEnv()
	router := httpx.NewRouter(httpx.Config{
		Logger:         a.logger,
		CORS:           troncmiddleware.CORSConfig{AllowedOrigins: a.cors, AllowCredentials: true},
		TrustedProxies: a.core.TrustedProxies,
		CDNProxies:     a.core.CDNProxies,
		CDNHeader:      a.core.CDNHeader,
		Security:       securityConfig(clientDir, a.homeserver, a.liveKit, a.logger),
	})
	health.Mount(router, a.dbCheck)
	router.Get("/version", versionHandler)
	apiref.Mount(router, refConfig())
	router.Get("/.well-known/matrix/client", func(w http.ResponseWriter, _ *http.Request) {
		httpjson.WriteJSON(w, http.StatusOK, map[string]any{"m.homeserver": map[string]string{"base_url": a.homeserver}})
	})
	if a.auth != nil {
		router.Route("/api", func(r chi.Router) {
			r.Use(session.RevokeOnLogout(a.session, a.auth, a.logger))
			a.auth.Mount(r)
			session.RegisterRoutes(r, a.session, a.auth)
			spaces.RegisterRoutes(r, a.spaces, a.auth)
			agents.RegisterRoutes(r, a.agents, a.auth)
			calls.RegisterRoutes(r, a.calls, a.auth)
			users.RegisterRoutes(r, a.users, a.auth)
			if a.syncer != nil && a.deletions != nil {
				spaces.RegisterSync(r, a.syncer)
				spaces.RegisterDeletions(r, a.deletions)
			}
		})
	}
	if spa.Available(clientDir) {
		router.Handle("/*", spa.Handler(spa.Config{Dir: clientDir}))
		a.logger.Info("serving client", slog.String("dir", clientDir))
	}
	return router
}

func main() {
	if healthcheck.Handle(os.Args) {
		return
	}
	os.Exit(run())
}
