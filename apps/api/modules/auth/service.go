// Package auth is Prisme's half of porte: porte owns the credential and the
// OIDC flow, this package wires it and decides which callers reach which route.
package auth

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/FacileStudio/porte"
	"github.com/FacileStudio/porte/oidc"
	"github.com/FacileStudio/porte/pg"
	"github.com/FacileStudio/porte/session"
)

// Service holds one session manager and the OIDC kit built over it.
type Service struct {
	sessions *session.Manager
	kit      *oidc.Kit
	services serviceKeys
}

// NewService builds porte over the identity tables. Prisme signs people in
// through Registre only, so there is no password kit.
func NewService(ctx context.Context, db *sql.DB, cfg porte.Config, keys [][]byte, logger *slog.Logger) (*Service, error) {
	store := pg.New(db)
	sessions, err := session.New(cfg, session.Deps{Sessions: store.Sessions(), Logger: logger})
	if err != nil {
		return nil, err
	}
	kit, err := oidc.New(ctx, cfg, oidc.Deps{
		Users:      pictureUsers{UserStore: store.Users()},
		Identities: store.Identities(),
		Codes:      store.LoginCodes(),
		Sessions:   sessions,
		Logger:     logger,
	})
	if err != nil {
		return nil, err
	}
	return &Service{sessions: sessions, kit: kit, services: serviceKeys{keys: keys, db: db}}, nil
}

// Optional authenticates when a credential is present and passes through
// otherwise.
func (s *Service) Optional(next http.Handler) http.Handler {
	return s.sessions.Optional(next)
}

// Require authenticates the caller. A route with a scope also accepts a
// service-key call acting for a user, and a delegated token carrying that
// scope. An empty scope keeps the route to browser sessions.
func (s *Service) Require(scope string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		sessions := s.sessions.RequireAuth(gate(scope, next))
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if scope != "" && s.services.matches(r) {
				s.services.serve(w, r, next)
				return
			}
			sessions.ServeHTTP(w, r)
		})
	}
}
