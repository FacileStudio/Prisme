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
	portepg "github.com/FacileStudio/porte/pg"
	"github.com/FacileStudio/porte/session"
)

// Service holds one session manager and the OIDC kit built over it.
type Service struct {
	sessions   *session.Manager
	kit        *oidc.Kit
	services   serviceKeys
	identities *portepg.IdentityStore
}

// NewService builds porte over the identity tables. Prisme signs people in
// through Registre only, so there is no password kit.
func NewService(ctx context.Context, db *sql.DB, cfg porte.Config, keys [][]byte, logger *slog.Logger) (*Service, error) {
	store, err := newPorteStore(db)
	if err != nil {
		return nil, err
	}
	resealInBackground(ctx, store, logger)
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
	services := serviceKeys{keys: keys, db: db, issuer: cfg.Issuer, identities: store.Identities(), logger: logger}
	return &Service{sessions: sessions, kit: kit, services: services, identities: store.Identities()}, nil
}

// Identities is porte's identity store, with the token sealer the sign-in
// uses, for readers such as the invitations badge.
func (s *Service) Identities() *portepg.IdentityStore {
	return s.identities
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
			if key := s.serviceKey(scope, r); key != 0 {
				s.services.serve(w, r, next, key)
				return
			}
			sessions.ServeHTTP(w, r)
		})
	}
}

// serviceKey is the matching service key's position, or 0 on a route that
// takes no service call or a request that carries no key.
func (s *Service) serviceKey(scope string, r *http.Request) int {
	if scope == "" {
		return 0
	}
	return s.services.match(r)
}
