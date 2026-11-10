package auth

import (
	"context"
	"crypto/subtle"
	"database/sql"
	stderrors "errors"
	"log/slog"
	"net/http"
	"net/mail"
	"strings"

	"github.com/FacileStudio/porte"
	portepg "github.com/FacileStudio/porte/pg"
	"github.com/FacileStudio/tronc/errors"
	"github.com/FacileStudio/tronc/httpjson"
)

// ActingForHeader names the Prisme user a service call acts for, by email.
const ActingForHeader = "X-Prisme-Acting-For"

// serviceActor is the porte actor recorded for a service-key call.
const serviceActor = "service-key"

// serviceKeys authenticates suite services (Bulle) that call Prisme with a
// shared key and name the user they act for. The issuer is the Registre
// issuer a sign-in files identities under, empty when SSO is off.
type serviceKeys struct {
	keys       [][]byte
	db         *sql.DB
	issuer     string
	identities *portepg.IdentityStore
	logger     *slog.Logger
}

// match returns the position of the presented key in PRISME_SERVICE_KEY,
// counting from 1, or 0 when the request carries none of them.
func (s serviceKeys) match(r *http.Request) int {
	token, ok := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
	if !ok || token == "" {
		return 0
	}
	found := 0
	for i, key := range s.keys {
		found = subtle.ConstantTimeSelect(subtle.ConstantTimeCompare([]byte(token), key), i+1, found)
	}
	return found
}

func (s serviceKeys) serve(w http.ResponseWriter, r *http.Request, next http.Handler, key int) {
	address, err := mail.ParseAddress(strings.TrimSpace(r.Header.Get(ActingForHeader)))
	if err != nil {
		httpjson.WriteError(w, errors.Invalid(ActingForHeader+" must carry the user's email"))
		return
	}
	email := strings.ToLower(address.Address)
	userID, err := s.resolve(r, email, key)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	identity := porte.Identity{UserID: userID, Email: email, EmailVerified: true, Actor: serviceActor}
	next.ServeHTTP(w, r.WithContext(porte.WithIdentity(r.Context(), identity)))
}

// resolve is UserFor, plus the acting subject when the call carries one.
func (s serviceKeys) resolve(r *http.Request, email string, key int) (int64, error) {
	subject, present, err := actingSubject(r)
	if err != nil {
		return 0, err
	}
	if present {
		return s.attach(r.Context(), email, subject, key)
	}
	userID, err := UserFor(r.Context(), s.db, email)
	if err != nil {
		return 0, errors.Internal("failed to resolve the user", err)
	}
	return userID, nil
}

// UserFor returns the porte user with this email, creating it as a first
// sign-in would. The address counts as verified: the caller (a service key or
// the Registre directory) vouches for it, and porte links the Registre
// identity to it on first SSO sign-in.
func UserFor(ctx context.Context, db *sql.DB, email string) (int64, error) {
	var userID int64
	err := db.QueryRowContext(ctx, `
		INSERT INTO porte_users (email, email_verified) VALUES ($1, true)
		ON CONFLICT (email) DO UPDATE SET email = EXCLUDED.email
		RETURNING id`, email).Scan(&userID)
	return userID, err
}

// refusal is an answer in tronc's error envelope with a code of Prisme's
// own, which tronc's status table does not know.
type refusal struct {
	status  int
	code    string
	message string
}

func (r refusal) Error() string { return r.message }

// writeServiceError answers err in tronc's error envelope, keeping a
// refusal's own status and code.
func writeServiceError(w http.ResponseWriter, err error) {
	var refused refusal
	if !stderrors.As(err, &refused) {
		httpjson.WriteError(w, err)
		return
	}
	httpjson.WriteJSON(w, refused.status, map[string]any{
		"error": map[string]string{"code": refused.code, "message": refused.message},
	})
}

// refused logs a refusal by account and key, never by email or subject.
func (s serviceKeys) refused(message string, userID int64, key int) {
	s.logger.Warn("auth: "+message+"; refused, nothing written",
		slog.Int64("user_id", userID), slog.String("service", serviceActor), slog.Int("service_key", key))
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
