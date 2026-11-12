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
	"github.com/FacileStudio/Prisme/apps/api/modules/push"
	"github.com/FacileStudio/porte"
	"github.com/go-chi/chi/v5"
)

type matrixCall struct {
	method, path, token, body string
	status                    int
	want                      string
}

func appserviceRouter(t *testing.T) chi.Router {
	quiet := slog.New(slog.NewTextHandler(io.Discard, nil))
	authService, err := auth.NewService(context.Background(), nil, porte.Config{}, auth.Guards{}, quiet)
	if err != nil {
		t.Fatal(err)
	}
	return buildRouter(app{
		logger: quiet, auth: authService, push: push.NewService(push.Deps{HSToken: "hs", Logger: quiet}),
		edge: edge{homeserver: "https://chat.example.test"},
	})
}

// Synapse calls these routes, and Matrix clients read the well-known. Their
// paths, status codes and bodies are a contract with software Prisme does not
// ship: a change here breaks push or room delivery without failing anything
// else, so the bodies are pinned byte for byte.
func TestTheMatrixFacingRoutesKeepTheirContract(t *testing.T) {
	router := appserviceRouter(t)
	forbidden := `{"error":{"code":"permission_denied","message":"bad hs_token"}}`
	unhandled := `{"error":{"code":"not_found","message":"not handled by this appservice"}}`
	calls := []matrixCall{
		{http.MethodPut, "/_matrix/app/v1/transactions/1", "Bearer hs", `{"events":[]}`, http.StatusOK, `{}`},
		{http.MethodPut, "/_matrix/app/v1/transactions/1", "", `{"events":[]}`, http.StatusForbidden, forbidden},
		{http.MethodPut, "/_matrix/app/v1/transactions/1", "Bearer as", `{"events":[]}`, http.StatusForbidden, forbidden},
		{http.MethodPut, "/_matrix/app/v1/transactions/1", "Bearer hs", `{`, http.StatusBadRequest,
			`{"error":{"code":"invalid_argument","message":"unreadable transaction"}}`},
		{http.MethodGet, "/_matrix/app/v1/users/@ada:hs", "Bearer hs", "", http.StatusNotFound, unhandled},
		{http.MethodGet, "/_matrix/app/v1/rooms/%23prisme-a:hs", "Bearer hs", "", http.StatusNotFound, unhandled},
		{http.MethodPost, "/_matrix/app/v1/ping", "Bearer hs", `{}`, http.StatusNotFound, unhandled},
		{http.MethodGet, "/.well-known/matrix/client", "", "", http.StatusOK,
			`{"m.homeserver":{"base_url":"https://chat.example.test"}}`},
	}
	for _, call := range calls {
		assertMatrixCall(t, router, call)
	}
}

func assertMatrixCall(t *testing.T, router http.Handler, call matrixCall) {
	t.Helper()
	request := httptest.NewRequest(call.method, call.path, strings.NewReader(call.body))
	if call.token != "" {
		request.Header.Set("Authorization", call.token)
	}
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	body := strings.TrimSpace(recorder.Body.String())
	if recorder.Code != call.status || body != call.want {
		t.Errorf("%s %s = %d %s, want %d %s", call.method, call.path, recorder.Code, body, call.status, call.want)
	}
	if ctype := recorder.Header().Get("Content-Type"); !strings.HasPrefix(ctype, "application/json") {
		t.Errorf("%s %s content-type = %q, want JSON", call.method, call.path, ctype)
	}
}
