package auth

import "github.com/go-chi/chi/v5"

// Mount registers porte's routes: config, logout and the OIDC flow.
func (s *Service) Mount(router chi.Router) {
	s.sessions.Mount(router)
	s.kit.Mount(router)
}
