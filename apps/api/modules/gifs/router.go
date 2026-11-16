package gifs

import (
	"net/http"
	"strconv"

	"github.com/FacileStudio/Prisme/apps/api/modules/auth"
	"github.com/FacileStudio/tronc/httpjson"

	"github.com/go-chi/chi/v5"
)

// RegisterRoutes mounts GET /gifs and GET /gifs/media for browser sessions.
func RegisterRoutes(router chi.Router, service *Service, authService *auth.Service) {
	router.With(authService.Require("")).Get("/gifs", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		offset, _ := strconv.Atoi(query.Get("offset"))
		page, err := service.Search(r.Context(), query.Get("q"), query.Get("lang"), offset)
		if err != nil {
			httpjson.WriteError(w, err)
			return
		}
		httpjson.WriteJSON(w, http.StatusOK, page)
	})
	router.With(authService.Require("")).Get("/gifs/media", func(w http.ResponseWriter, r *http.Request) {
		media, err := service.Media(r.Context(), r.URL.Query().Get("u"))
		if err != nil {
			httpjson.WriteError(w, err)
			return
		}
		w.Header().Set("Content-Type", media.ContentType)
		w.Header().Set("Cache-Control", "private, max-age=86400")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		_, _ = w.Write(media.Body)
	})
}
