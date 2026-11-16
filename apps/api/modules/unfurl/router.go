package unfurl

import (
	"net/http"

	"github.com/FacileStudio/Prisme/apps/api/modules/auth"
	"github.com/FacileStudio/tronc/httpjson"

	"github.com/go-chi/chi/v5"
)

// RegisterRoutes mounts GET /unfurl for browser sessions.
func RegisterRoutes(router chi.Router, service *Service, authService *auth.Service) {
	router.With(authService.Require("")).Get("/unfurl", func(w http.ResponseWriter, r *http.Request) {
		preview, err := service.Preview(r.Context(), r.URL.Query().Get("url"))
		if err != nil {
			httpjson.WriteError(w, err)
			return
		}
		w.Header().Set("Cache-Control", "private, max-age=600")
		httpjson.WriteJSON(w, http.StatusOK, preview)
	})
}
