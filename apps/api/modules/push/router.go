package push

import (
	"encoding/json"
	"net/http"

	"github.com/FacileStudio/Prisme/apps/api/modules/auth"
	"github.com/FacileStudio/porte"
	"github.com/FacileStudio/tronc/errors"
	"github.com/FacileStudio/tronc/httpjson"

	"github.com/go-chi/chi/v5"
)

const maxTransactionBytes = 8 << 20

// RegisterRoutes mounts GET /push and POST /push/devices for signed-in
// sessions only: a delegated token has no phone behind it.
func RegisterRoutes(router chi.Router, service *Service, authService *auth.Service) {
	router.With(authService.Require("")).Get("/push", func(w http.ResponseWriter, _ *http.Request) {
		httpjson.WriteJSON(w, http.StatusOK, service.Config())
	})
	router.With(authService.Require("")).Post("/push/devices", func(w http.ResponseWriter, r *http.Request) {
		var req DeviceRequest
		if err := httpjson.DecodeJSON(w, r, &req); err != nil {
			httpjson.WriteError(w, err)
			return
		}
		identity, _ := porte.From(r.Context())
		if err := service.Register(r.Context(), identity.UserID, identity.SessionID, req.Token); err != nil {
			httpjson.WriteError(w, err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})
}

// RegisterAppService mounts the appservice side Synapse calls, at the router's
// root: the event transactions, and a 404 for every user and room query so
// Synapse never mistakes the SPA fallback for an answer.
func RegisterAppService(router chi.Router, service *Service) {
	router.Put(TransactionsPath, func(w http.ResponseWriter, r *http.Request) {
		if !service.Authorized(r.Header.Get("Authorization")) {
			httpjson.WriteError(w, errors.Forbidden("bad hs_token"))
			return
		}
		var txn transaction
		if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, maxTransactionBytes)).Decode(&txn); err != nil {
			httpjson.WriteError(w, errors.Invalid("unreadable transaction"))
			return
		}
		service.Receive(txn)
		httpjson.WriteJSON(w, http.StatusOK, struct{}{})
	})
	router.Handle("/_matrix/*", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		httpjson.WriteError(w, errors.NotFound("not handled by this appservice"))
	}))
}
