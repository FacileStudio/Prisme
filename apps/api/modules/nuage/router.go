package nuage

import (
	"net/http"
	"strconv"

	"github.com/FacileStudio/Prisme/apps/api/internal/nuage"
	"github.com/FacileStudio/Prisme/apps/api/modules/auth"
	"github.com/FacileStudio/porte"
	"github.com/FacileStudio/tronc/errors"
	"github.com/FacileStudio/tronc/httpjson"

	"github.com/go-chi/chi/v5"
)

const maxShareBytes = 4 << 10

// RegisterRoutes mounts the Nuage picker for browser sessions only.
func RegisterRoutes(router chi.Router, service *Service, authService *auth.Service) {
	router.With(authService.Require("")).Route("/nuage", func(r chi.Router) {
		r.Get("/spaces", func(w http.ResponseWriter, req *http.Request) {
			resp, err := service.Spaces(req.Context(), userID(req))
			respond(w, resp, err)
		})
		r.Get("/list", func(w http.ResponseWriter, req *http.Request) {
			place, err := placeOf(req)
			if err != nil {
				httpjson.WriteError(w, err)
				return
			}
			resp, err := service.List(req.Context(), userID(req), place)
			respond(w, resp, err)
		})
		r.Get("/search", func(w http.ResponseWriter, req *http.Request) {
			resp, err := service.Search(req.Context(), userID(req), req.URL.Query().Get("q"))
			respond(w, resp, err)
		})
		r.Post("/share", func(w http.ResponseWriter, req *http.Request) {
			var body ShareRequest
			if err := httpjson.DecodeJSONLimit(w, req, &body, maxShareBytes); err != nil {
				httpjson.WriteError(w, err)
				return
			}
			resp, err := service.Share(req.Context(), userID(req), body)
			respond(w, resp, err)
		})
	})
}

func userID(req *http.Request) int64 {
	identity, _ := porte.From(req.Context())
	return identity.UserID
}

func placeOf(req *http.Request) (nuage.Place, error) {
	folder, err := optionalID(req, "folder_id")
	if err != nil {
		return nuage.Place{}, err
	}
	space, err := optionalID(req, "space_id")
	return nuage.Place{FolderID: folder, SpaceID: space}, err
}

func optionalID(req *http.Request, key string) (*int64, error) {
	raw := req.URL.Query().Get(key)
	if raw == "" {
		return nil, nil
	}
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || id <= 0 {
		return nil, errors.Invalid(key + " must be a positive integer")
	}
	return &id, nil
}

func respond(w http.ResponseWriter, body any, err error) {
	if err != nil {
		httpjson.WriteError(w, err)
		return
	}
	w.Header().Set("Cache-Control", "private, no-store")
	httpjson.WriteJSON(w, http.StatusOK, body)
}
