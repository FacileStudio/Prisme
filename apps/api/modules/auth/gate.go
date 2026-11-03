package auth

import (
	"net/http"
	"slices"

	"github.com/FacileStudio/porte"
	"github.com/FacileStudio/tronc/errors"
	"github.com/FacileStudio/tronc/httpjson"
)

// Scopes a delegated Registre token may carry for Prisme.
const (
	ScopeRead   = "prisme.read"
	ScopeAgents = "prisme.agents"
)

// gate refuses a delegated bearer (a token with no browser session behind
// it) unless the route accepts its scope. porte's own ScopePrefix gate knows
// only read and write; Prisme grants agents management as its own scope.
func gate(scope string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		identity, ok := porte.From(r.Context())
		if !ok || identity.UserID == 0 {
			httpjson.WriteError(w, errors.Unauthorized("missing auth"))
			return
		}
		if Delegated(identity) && (scope == "" || !slices.Contains(identity.Scopes, scope)) {
			httpjson.WriteError(w, errors.Forbidden("insufficient scope"))
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Delegated reports whether the identity came from a bearer JWT rather than a
// Prisme session.
func Delegated(identity porte.Identity) bool {
	return identity.SessionID == 0
}
