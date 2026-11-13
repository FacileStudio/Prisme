package users

import (
	"net/http"

	"github.com/FacileStudio/Prisme/apps/api/modules/auth"
	"github.com/FacileStudio/porte"
	"github.com/FacileStudio/tronc/httpjson"

	"github.com/go-chi/chi/v5"
)

// RegisterRoutes mounts the directory routes for browser sessions.
func RegisterRoutes(router chi.Router, service *Service, authService *auth.Service) {
	router.With(authService.Require("")).Get("/users", func(w http.ResponseWriter, r *http.Request) {
		identity, _ := porte.From(r.Context())
		resp, err := service.Search(r.Context(), identity.UserID, r.URL.Query().Get("q"))
		respond(w, resp, err)
	})
	router.With(authService.Require("")).Get("/users/agents", func(w http.ResponseWriter, r *http.Request) {
		identity, _ := porte.From(r.Context())
		resp, err := service.Agents(r.Context(), identity.UserID)
		respond(w, resp, err)
	})
}

func respond(w http.ResponseWriter, body any, err error) {
	if err != nil {
		httpjson.WriteError(w, err)
		return
	}
	httpjson.WriteJSON(w, http.StatusOK, body)
}
