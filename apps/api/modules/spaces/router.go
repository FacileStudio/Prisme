package spaces

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/FacileStudio/Prisme/apps/api/modules/auth"
	"github.com/FacileStudio/porte"
	"github.com/FacileStudio/porte/directory"
	"github.com/FacileStudio/tronc/errors"
	"github.com/FacileStudio/tronc/httpjson"

	"github.com/go-chi/chi/v5"
)

// RegisterRoutes mounts the space routes. Listing accepts a delegated token
// with prisme.read so Bulle can offer a space picker.
func RegisterRoutes(router chi.Router, service *Service, authService *auth.Service) {
	h := handlers{service: service}
	router.With(authService.Require(auth.ScopeRead)).Get("/spaces", h.list)
	router.With(authService.Require("")).Post("/spaces", h.create)
	router.With(authService.Require("")).Post("/spaces/{id}/channels", h.createChannel)
}

// RegisterSync mounts Registre's nudge on the /api group, so it answers at
// directory.SyncPath. It sits outside every auth middleware on purpose: the
// nudge is an anonymous POST that only makes the next authenticated pull
// happen sooner.
func RegisterSync(router chi.Router, syncer *directory.Syncer) {
	router.Method(http.MethodPost, strings.TrimPrefix(directory.SyncPath, "/api"), syncer.Handler())
}

type handlers struct {
	service *Service
}

func (h handlers) list(w http.ResponseWriter, r *http.Request) {
	identity, _ := porte.From(r.Context())
	resp, err := h.service.List(r.Context(), identity.UserID)
	respond(w, http.StatusOK, resp, err)
}

func (h handlers) create(w http.ResponseWriter, r *http.Request) {
	var req CreateSpaceRequest
	if err := httpjson.DecodeJSON(w, r, &req); err != nil {
		httpjson.WriteError(w, err)
		return
	}
	identity, _ := porte.From(r.Context())
	resp, err := h.service.Create(r.Context(), identity.UserID, req)
	respond(w, http.StatusCreated, resp, err)
}

func (h handlers) createChannel(w http.ResponseWriter, r *http.Request) {
	spaceID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		httpjson.WriteError(w, errors.Invalid("space id must be a number"))
		return
	}
	var req CreateChannelRequest
	if err := httpjson.DecodeJSON(w, r, &req); err != nil {
		httpjson.WriteError(w, err)
		return
	}
	identity, _ := porte.From(r.Context())
	resp, err := h.service.CreateChannel(r.Context(), identity.UserID, spaceID, req)
	respond(w, http.StatusCreated, resp, err)
}

func respond(w http.ResponseWriter, status int, body any, err error) {
	if err != nil {
		httpjson.WriteError(w, err)
		return
	}
	httpjson.WriteJSON(w, status, body)
}
