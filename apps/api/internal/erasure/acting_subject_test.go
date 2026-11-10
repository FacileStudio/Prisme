package erasure

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"
	"time"

	"github.com/FacileStudio/Prisme/apps/api/internal/accounts"
	"github.com/FacileStudio/Prisme/apps/api/modules/auth"
	"github.com/FacileStudio/porte"
	"github.com/FacileStudio/porte/directory"
)

// fakeRegistreIdP is fakeRegistre with an OIDC discovery document, so a real
// auth service can be built over it: sign-in identities are filed under its
// URL, and its deletion list names the given subject.
func fakeRegistreIdP(t *testing.T, subject string) string {
	t.Helper()
	at := time.Now().Add(-time.Minute).UTC()
	var server *httptest.Server
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/.well-known/openid-configuration" {
			_ = json.NewEncoder(w).Encode(map[string]any{
				"issuer": server.URL, "authorization_endpoint": server.URL + "/authorize",
				"token_endpoint": server.URL + "/token", "jwks_uri": server.URL + "/jwks",
				"userinfo_endpoint": server.URL + "/userinfo", "id_token_signing_alg_values_supported": []string{"RS256"},
			})
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"deleted": []directory.Deleted{{Subject: subject, DeletedAt: at}}})
	}))
	t.Cleanup(server.Close)
	return server.URL
}

// bulleCall is one service-key call as Bulle makes it, to a route that
// gives the person a Matrix account, as POST /api/agents does.
func bulleCall(t *testing.T, w world, service *auth.Service, headers map[string]string) {
	t.Helper()
	quiet := slog.New(slog.NewTextHandler(io.Discard, nil))
	ensure := http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		identity, _ := porte.From(r.Context())
		if _, err := accounts.NewService(w.db, w.eraser.deps.Matrix, quiet).Ensure(r.Context(), identity.UserID); err != nil {
			t.Errorf("ensure: %v", err)
		}
		rw.WriteHeader(http.StatusNoContent)
	})
	request := httptest.NewRequest(http.MethodPost, "/api/agents", nil)
	request.Header.Set("Authorization", "Bearer bulle-key")
	for name, value := range headers {
		request.Header.Set(name, value)
	}
	recorder := httptest.NewRecorder()
	service.Require(auth.ScopeAgents)(ensure).ServeHTTP(recorder, request)
	if recorder.Code != http.StatusNoContent {
		t.Fatalf("service call %v = %d %s", headers, recorder.Code, recorder.Body.String())
	}
}

// actingWorld is seedWorld behind a real auth service whose issuer is a
// fake Registre that lists sub-bulle as erased, with its deletion receiver.
type actingWorld struct {
	world
	service   *auth.Service
	deletions *directory.Deletions
	issuer    string
}

func newActingWorld(t *testing.T) actingWorld {
	t.Helper()
	w := seedWorld(t)
	issuer := fakeRegistreIdP(t, "sub-bulle")
	mustExec(t, w.db, `UPDATE porte_identities SET provider = ?`, issuer)
	sqlDB, _ := w.db.DB()
	quiet := slog.New(slog.NewTextHandler(io.Discard, nil))
	service, err := auth.NewService(context.Background(), sqlDB, porte.Config{
		Issuer: issuer, ClientID: "prisme", ClientSecret: "s",
		RedirectURL: "http://localhost/api/auth/oidc/callback", SuccessURL: "http://localhost/",
	}, auth.ParseServiceKeys("bulle-key"), quiet)
	if err != nil {
		t.Fatalf("auth service: %v", err)
	}
	w.eraser.deps.Issuer = issuer
	return actingWorld{world: w, service: service, deletions: receiver(t, w, issuer), issuer: issuer}
}

// A person Bulle provisions by email alone cannot be mapped by a deletion
// notice, which stays pending. Once Bulle calls again with the subject, the
// next delivery erases them from Prisme and Synapse and is answered as done.
func TestAnAccountLinkedByTheActingSubjectIsErasedByANotice(t *testing.T) {
	a := newActingWorld(t)
	w, service, deletions, issuer := a.world, a.service, a.deletions, a.issuer
	ctx := context.Background()
	bulleCall(t, w, service, map[string]string{auth.ActingForHeader: "bulle@facile.studio"})
	var mxid string
	w.db.Raw(`SELECT a.matrix_user_id FROM accounts a JOIN porte_users u ON u.id = a.user_id WHERE u.email = 'bulle@facile.studio'`).Scan(&mxid)
	if mxid == "" {
		t.Fatal("the service call made no Matrix account")
	}
	if err := deletions.Poll(ctx); err == nil {
		t.Fatal("a notice for an unlinked account was answered as done")
	}
	if err := w.eraser.ErasePerson(ctx, issuer, "sub-bulle"); !errors.Is(err, ErrUnproven) {
		t.Fatalf("notice before the subject = %v, want ErrUnproven", err)
	}

	bulleCall(t, w, service, map[string]string{auth.ActingForHeader: "bulle@facile.studio", auth.ActingSubjectHeader: "sub-bulle"})
	if err := deletions.Poll(ctx); err != nil {
		t.Fatalf("notice after the subject = %v, want done", err)
	}
	assertErasedBy(t, w, mxid)
	if err := w.eraser.ErasePerson(ctx, issuer, "sub-bulle"); err != nil {
		t.Fatalf("a repeated notice = %v, want done", err)
	}
	if err := deletions.Poll(ctx); err != nil {
		t.Fatalf("a repeated poll = %v, want done", err)
	}
}

// assertErasedBy checks the Bulle-provisioned person left Prisme and Synapse
// and the seeded people kept everything.
func assertErasedBy(t *testing.T, w world, mxid string) {
	t.Helper()
	if n := countRows(t, w.db, `SELECT count(*) FROM porte_users WHERE email = 'bulle@facile.studio'`); n != 0 {
		t.Error("the account outlived the notice")
	}
	if n := countRows(t, w.db, `SELECT count(*) FROM accounts WHERE matrix_user_id = ?`, mxid); n != 0 {
		t.Error("the Matrix account row outlived the notice")
	}
	w.synapse.read(func(f *fakeSynapse) {
		if !slices.Contains(f.calls, "POST /v1/deactivate/"+mxid) || !f.users[mxid] {
			t.Errorf("Synapse did not deactivate %s: %v", mxid, f.calls)
		}
	})
	assertUntouched(t, w.db, w.gone)
	assertUntouched(t, w.db, w.other)
}
