package auth

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/FacileStudio/porte"
	"github.com/FacileStudio/porte/directory"
	"github.com/FacileStudio/porte/oidc/jwt"
)

// handledState is a deletion receiver's saved state naming one handled subject.
type handledState string

func (s handledState) Load(context.Context) ([]byte, error) {
	return []byte(`{"done":{"` + string(s) + `":"` + time.Now().UTC().Format(time.RFC3339) + `"}}`), nil
}

func (handledState) Save(context.Context, []byte) error { return nil }

// discoveryIssuer serves the OIDC discovery document porte's kit needs to start.
func discoveryIssuer(t *testing.T) string {
	t.Helper()
	var server *httptest.Server
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		err := json.NewEncoder(w).Encode(map[string]any{
			"issuer": server.URL, "authorization_endpoint": server.URL + "/authorize",
			"token_endpoint": server.URL + "/token", "jwks_uri": server.URL + "/jwks",
			"id_token_signing_alg_values_supported": []string{"RS256"},
		})
		if err != nil {
			t.Errorf("discovery: %v", err)
		}
	}))
	t.Cleanup(server.Close)
	return server.URL
}

// A subject the deletion receiver handled gets no account from porte's own
// write path, even one that never held an account here and so left no
// tombstone in erased_subjects: a bearer token minted before the erasure is
// refused and nothing is written.
func TestTheDeletionReceiverRefusesAnAccountToAHandledSubject(t *testing.T) {
	db := openAuthDatabase(t)
	issuer := discoveryIssuer(t)
	receiver := &directory.Deletions{Client: directory.Client{Issuer: issuer}, State: handledState("sub-gone")}
	s, err := NewService(context.Background(), db, porte.Config{
		Issuer: issuer, ClientID: "prisme", ClientSecret: "s",
		RedirectURL: "http://localhost/api/auth/oidc/callback", SuccessURL: "http://localhost/",
	}, Guards{Erased: receiver}, slog.New(slog.NewTextHandler(io.Discard, nil)))
	if err != nil {
		t.Fatalf("auth service: %v", err)
	}
	verified := true
	_, err = s.kit.ProvisionBearer(context.Background(), jwt.Claims{
		Subject: "sub-gone", Email: "gone@facile.studio", Profile: jwt.Profile{EmailVerified: &verified},
	})
	if !errors.Is(err, porte.ErrSubjectErased) {
		t.Fatalf("ProvisionBearer = %v, want porte.ErrSubjectErased", err)
	}
	if n := countIn(t, db, `SELECT count(*) FROM porte_users`) + countIn(t, db, `SELECT count(*) FROM porte_identities`); n != 0 {
		t.Errorf("a refused provisioning wrote %d rows", n)
	}
}
