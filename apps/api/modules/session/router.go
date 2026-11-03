package session

import (
	"log/slog"
	"net/http"

	"github.com/FacileStudio/Prisme/apps/api/modules/auth"
	"github.com/FacileStudio/porte"
	"github.com/FacileStudio/tronc/httpjson"

	"github.com/go-chi/chi/v5"
)

// RegisterRoutes mounts GET /session for browser sessions only: a delegated
// token has no device to hand out.
func RegisterRoutes(router chi.Router, service *Service, authService *auth.Service) {
	router.With(authService.Require("")).Get("/session", func(w http.ResponseWriter, r *http.Request) {
		identity, _ := porte.From(r.Context())
		resp, err := service.Get(r.Context(), identity.UserID, identity.SessionID)
		if err != nil {
			httpjson.WriteError(w, err)
			return
		}
		httpjson.WriteJSON(w, http.StatusOK, resp)
	})
}

// RevokeOnLogout runs before porte's POST /auth/logout and revokes the Matrix
// device of the session being ended, so signing out of Prisme signs the
// browser out of Matrix too.
func RevokeOnLogout(service *Service, authService *auth.Service, logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		revoke := authService.Optional(revoking(service, logger, next))
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost && r.URL.Path == "/api"+porte.RouteLogout {
				revoke.ServeHTTP(w, r)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func revoking(service *Service, logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		identity, ok := porte.From(r.Context())
		if ok && identity.SessionID != 0 {
			if err := service.Revoke(r.Context(), identity.SessionID); err != nil {
				logger.Warn("matrix logout failed", slog.Any("error", err))
			}
		}
		next.ServeHTTP(w, r)
	})
}
