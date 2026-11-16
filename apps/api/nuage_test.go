package main

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/FacileStudio/Prisme/apps/api/modules/auth"
	"github.com/FacileStudio/Prisme/apps/api/modules/nuage"
	"github.com/FacileStudio/porte"
)

// The Nuage picker acts as the signed-in person in their drive, so it answers
// only a session: without one every route stops before the handler, which
// would otherwise answer 503 with no Nuage wired. Its Require("") scope never
// reads a service key, so Bulle's key is only ever tried as a session token.
func TestTheNuagePickerNeedsABrowserSession(t *testing.T) {
	quiet := slog.New(slog.NewTextHandler(io.Discard, nil))
	authService, err := auth.NewService(context.Background(), nil, porte.Config{}, auth.Guards{}, quiet)
	if err != nil {
		t.Fatal(err)
	}
	router := buildRouter(app{logger: quiet, auth: authService, suite: suiteApps{nuage: nuage.NewService(nil, quiet)}})
	for _, call := range []struct{ method, path string }{
		{http.MethodGet, "/api/nuage/spaces"},
		{http.MethodGet, "/api/nuage/list"},
		{http.MethodGet, "/api/nuage/search?q=plan"},
		{http.MethodPost, "/api/nuage/share"},
	} {
		request := httptest.NewRequest(call.method, call.path, strings.NewReader(`{"file_id":1}`))
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("X-Facile-CSRF", "1")
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, request)
		if recorder.Code != http.StatusUnauthorized {
			t.Errorf("%s %s without a session = %d, want 401", call.method, call.path, recorder.Code)
		}
	}
}
