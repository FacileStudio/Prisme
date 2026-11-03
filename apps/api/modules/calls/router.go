package calls

import (
	"io"
	"net/http"

	"github.com/FacileStudio/Prisme/apps/api/modules/auth"
	"github.com/FacileStudio/porte"
	"github.com/FacileStudio/tronc/errors"
	"github.com/FacileStudio/tronc/httpjson"

	"github.com/go-chi/chi/v5"
)

const maxEventBytes = 64 << 10

// RegisterRoutes mounts POST /calls for browser sessions and POST
// /echo/events for Echo, which authenticates by signing the body.
func RegisterRoutes(router chi.Router, service *Service, authService *auth.Service) {
	router.With(authService.Require("")).Post("/calls", func(w http.ResponseWriter, r *http.Request) {
		var req Request
		if err := httpjson.DecodeJSON(w, r, &req); err != nil {
			httpjson.WriteError(w, err)
			return
		}
		identity, _ := porte.From(r.Context())
		resp, err := service.Start(r.Context(), identity.UserID, req)
		if err != nil {
			httpjson.WriteError(w, err)
			return
		}
		httpjson.WriteJSON(w, http.StatusOK, resp)
	})
	router.Post(EventsPath, func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, maxEventBytes))
		if err != nil {
			httpjson.WriteError(w, errors.Invalid("unreadable Echo event"))
			return
		}
		if err := service.Receive(r.Context(), body, r.Header.Get("X-Echo-Signature")); err != nil {
			httpjson.WriteError(w, err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})
}
