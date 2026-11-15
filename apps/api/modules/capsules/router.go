package capsules

import (
	"net/http"

	"github.com/FacileStudio/Prisme/apps/api/modules/auth"
	"github.com/FacileStudio/tronc/httpjson"

	"github.com/go-chi/chi/v5"
)

const maxSealBytes = 1 << 20

// RegisterRoutes mounts the capsule relay for browser sessions.
func RegisterRoutes(router chi.Router, service *Service, authService *auth.Service) {
	router.With(authService.Require("")).Post("/capsules", func(w http.ResponseWriter, r *http.Request) {
		var req SealRequest
		if err := httpjson.DecodeJSONLimit(w, r, &req, maxSealBytes); err != nil {
			httpjson.WriteError(w, err)
			return
		}
		resp, err := service.Seal(r.Context(), req)
		respond(w, resp, err)
	})
	router.With(authService.Require("")).Get("/capsules/{id}", func(w http.ResponseWriter, r *http.Request) {
		resp, err := service.Meta(r.Context(), chi.URLParam(r, "id"))
		respond(w, resp, err)
	})
	router.With(authService.Require("")).Post("/capsules/{id}/open", func(w http.ResponseWriter, r *http.Request) {
		resp, err := service.Open(r.Context(), chi.URLParam(r, "id"))
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
