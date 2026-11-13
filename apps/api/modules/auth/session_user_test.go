package auth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/FacileStudio/porte"
)

// TestSessionUserResolvesOnlyALiveSession checks the rate limiter's view of a
// caller: the user behind a session cookie or bearer token, and nobody behind
// a made-up token, a revoked one or no credential at all.
func TestSessionUserResolvesOnlyALiveSession(t *testing.T) {
	db := openAuthDatabase(t)
	service := newTestService(t, db, porte.Config{})
	var userID int64
	if err := db.QueryRow(`INSERT INTO porte_users (email) VALUES ('ada@facile.studio') RETURNING id`).Scan(&userID); err != nil {
		t.Fatalf("seed: %v", err)
	}
	ctx := context.Background()
	token, session, err := service.sessions.Issue(ctx, userID, "test")
	if err != nil {
		t.Fatalf("issue: %v", err)
	}

	bearer := bearerRequest(token)
	wantSessionUser(t, service, bearer, userID)
	wantSessionUser(t, service, cookieRequest(service, token), userID)
	wantSessionUser(t, service, bearerRequest("not-a-session"), 0)
	wantSessionUser(t, service, httptest.NewRequest(http.MethodGet, "/api/spaces", nil), 0)

	if err := service.sessions.Revoke(ctx, userID, session.ID); err != nil {
		t.Fatalf("revoke: %v", err)
	}
	wantSessionUser(t, service, bearer, 0)
}

func wantSessionUser(t *testing.T, service *Service, r *http.Request, want int64) {
	t.Helper()
	got, ok := service.SessionUser(r)
	if ok != (want != 0) || got != want {
		t.Fatalf("SessionUser = %d, %v; want %d", got, ok, want)
	}
}

func bearerRequest(token string) *http.Request {
	request := httptest.NewRequest(http.MethodGet, "/api/spaces", nil)
	request.Header.Set("Authorization", "Bearer "+token)
	return request
}

func cookieRequest(service *Service, token string) *http.Request {
	issued := httptest.NewRecorder()
	service.sessions.SetSessionCookie(issued, httptest.NewRequest(http.MethodGet, "/api/spaces", nil), token)
	request := httptest.NewRequest(http.MethodGet, "/api/spaces", nil)
	for _, cookie := range issued.Result().Cookies() {
		request.AddCookie(cookie)
	}
	return request
}
