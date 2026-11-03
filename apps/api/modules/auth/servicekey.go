package auth

import (
	"context"
	"crypto/subtle"
	"database/sql"
	"net/http"
	"net/mail"
	"strings"

	"github.com/FacileStudio/porte"
	"github.com/FacileStudio/tronc/errors"
	"github.com/FacileStudio/tronc/httpjson"
)

// ActingForHeader names the Prisme user a service call acts for, by email.
const ActingForHeader = "X-Prisme-Acting-For"

// serviceActor is the porte actor recorded for a service-key call.
const serviceActor = "service-key"

// serviceKeys authenticates suite services (Bulle) that call Prisme with a
// shared key and name the user they act for.
type serviceKeys struct {
	keys [][]byte
	db   *sql.DB
}

func (s serviceKeys) matches(r *http.Request) bool {
	token, ok := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
	if !ok || token == "" {
		return false
	}
	found := 0
	for _, key := range s.keys {
		found |= subtle.ConstantTimeCompare([]byte(token), key)
	}
	return found == 1
}

func (s serviceKeys) serve(w http.ResponseWriter, r *http.Request, next http.Handler) {
	address, err := mail.ParseAddress(strings.TrimSpace(r.Header.Get(ActingForHeader)))
	if err != nil {
		httpjson.WriteError(w, errors.Invalid(ActingForHeader+" must carry the user's email"))
		return
	}
	email := strings.ToLower(address.Address)
	userID, err := s.userFor(r.Context(), email)
	if err != nil {
		httpjson.WriteError(w, errors.Internal("failed to resolve the user", err))
		return
	}
	identity := porte.Identity{UserID: userID, Email: email, EmailVerified: true, Actor: serviceActor}
	next.ServeHTTP(w, r.WithContext(porte.WithIdentity(r.Context(), identity)))
}

// userFor returns the porte user with this email, creating it as a first
// sign-in would. The address counts as verified: the calling service vouches
// for it, and porte links the Registre identity to it on first SSO sign-in.
func (s serviceKeys) userFor(ctx context.Context, email string) (int64, error) {
	var userID int64
	err := s.db.QueryRowContext(ctx, `
		INSERT INTO porte_users (email, email_verified) VALUES ($1, true)
		ON CONFLICT (email) DO UPDATE SET email = EXCLUDED.email
		RETURNING id`, email).Scan(&userID)
	return userID, err
}

// ParseServiceKeys splits PRISME_SERVICE_KEY on commas, dropping blanks.
func ParseServiceKeys(raw string) [][]byte {
	var keys [][]byte
	for _, key := range strings.Split(raw, ",") {
		if key = strings.TrimSpace(key); key != "" {
			keys = append(keys, []byte(key))
		}
	}
	return keys
}
