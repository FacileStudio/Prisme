package agents

import (
	"net/http"
	"net/url"

	"github.com/FacileStudio/Prisme/apps/api/modules/auth"
	"github.com/FacileStudio/porte"
	"github.com/FacileStudio/tronc/httpjson"

	"github.com/go-chi/chi/v5"
)

// RegisterRoutes mounts the agent routes. Each accepts a browser session or a
// delegated token carrying prisme.agents, which is how Bulle provisions.
func RegisterRoutes(router chi.Router, service *Service, authService *auth.Service) {
	h := handlers{service: service}
	router.Route("/agents", func(r chi.Router) {
		r.Use(authService.Require(auth.ScopeAgents))
		r.Get("/", h.list)
		r.Post("/", h.create)
		r.Get("/check", h.check)
		r.Put("/{user_id}/rooms", h.rooms)
		r.Delete("/{user_id}", h.deactivate)
	})
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
	var req CreateRequest
	if err := httpjson.DecodeJSON(w, r, &req); err != nil {
		httpjson.WriteError(w, err)
		return
	}
	identity, _ := porte.From(r.Context())
	resp, err := h.service.Create(r.Context(), identity.UserID, req)
	respond(w, http.StatusOK, resp, err)
}

func (h handlers) check(w http.ResponseWriter, r *http.Request) {
	identity, _ := porte.From(r.Context())
	name := r.URL.Query().Get("name")
	resp, err := h.service.CheckAvailable(r.Context(), identity.UserID, name)
	respond(w, http.StatusOK, resp, err)
}

func (h handlers) rooms(w http.ResponseWriter, r *http.Request) {
	var req RoomsRequest
	if err := httpjson.DecodeJSON(w, r, &req); err != nil {
		httpjson.WriteError(w, err)
		return
	}
	identity, _ := porte.From(r.Context())
	resp, err := h.service.SetRooms(r.Context(), identity.UserID, agentID(r), req.RoomIDs)
	respond(w, http.StatusOK, resp, err)
}

func (h handlers) deactivate(w http.ResponseWriter, r *http.Request) {
	identity, _ := porte.From(r.Context())
	if err := h.service.Deactivate(r.Context(), identity.UserID, agentID(r)); err != nil {
		httpjson.WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func agentID(r *http.Request) string {
	raw := chi.URLParam(r, "user_id")
	if decoded, err := url.PathUnescape(raw); err == nil {
		return decoded
	}
	return raw
}

func respond(w http.ResponseWriter, status int, body any, err error) {
	if err != nil {
		httpjson.WriteError(w, err)
		return
	}
	httpjson.WriteJSON(w, status, body)
}
