package main

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/FacileStudio/Prisme/apps/api/modules/auth"
	"github.com/FacileStudio/Prisme/apps/api/modules/spaces"
	"github.com/FacileStudio/porte"
	"github.com/FacileStudio/porte/directory"
	"github.com/FacileStudio/tronc/apiref"
	"github.com/go-chi/chi/v5"
)

func testRouter(t *testing.T) chi.Router {
	quiet := slog.New(slog.NewTextHandler(io.Discard, nil))
	authService, err := auth.NewService(context.Background(), nil, porte.Config{}, nil, quiet)
	if err != nil {
		t.Fatal(err)
	}
	return buildRouter(app{logger: quiet, auth: authService})
}

// A hand-written registry rots the moment a route is added without it, which is
// how an API ends up with documentation that lies. This turns that drift into a
// failing build. If it fails, fix the registry in the module's
// documentation.go - never the test.
func TestEveryRouteIsDocumented(t *testing.T) {
	if missing := apiref.Undocumented(testRouter(t), refConfig(), "/api/auth", "/.well-known", "/version", "/_matrix"); len(missing) > 0 {
		t.Errorf("routes missing from the API registry:\n  %s", strings.Join(missing, "\n  "))
	}
}

func registreRouter(t *testing.T) chi.Router {
	quiet := slog.New(slog.NewTextHandler(io.Discard, nil))
	authService, err := auth.NewService(context.Background(), nil, porte.Config{}, nil, quiet)
	if err != nil {
		t.Fatal(err)
	}
	deletions := spaces.NewDeletions(&spaces.Service{}, directory.Client{}, nil, quiet)
	return buildRouter(app{logger: quiet, auth: authService, directory: directorySync{syncer: &directory.Syncer{}, deletions: deletions}})
}

// Registre nudges anonymously, so the route must answer without a session.
func TestTheRegistreNudgeNeedsNoCredentials(t *testing.T) {
	recorder := httptest.NewRecorder()
	registreRouter(t).ServeHTTP(recorder, httptest.NewRequest(http.MethodPost, directory.SyncPath, nil))
	if recorder.Code != http.StatusAccepted {
		t.Fatalf("POST %s = %d, want 202", directory.SyncPath, recorder.Code)
	}
}

// The deletion push is verified by its signature, not by a session: an unsigned
// body reaches the receiver and is refused as a bad notice, not as a missing login.
func TestTheDeletionPushReachesTheReceiverWithoutASession(t *testing.T) {
	recorder := httptest.NewRecorder()
	registreRouter(t).ServeHTTP(recorder, httptest.NewRequest(http.MethodPost, directory.DeletionsPath, strings.NewReader("x")))
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("POST %s = %d, want 400", directory.DeletionsPath, recorder.Code)
	}
}

func TestTheSecurityHeadersAllowCallsAndNameTheHomeserver(t *testing.T) {
	router := buildRouter(app{logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
		edge: edge{homeserver: "https://chat.example.test", liveKit: "wss://live.example.test:7881"}})
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/version", nil))
	if recorder.Code != http.StatusOK || !strings.Contains(recorder.Body.String(), `"version":"dev"`) {
		t.Fatalf("GET /version = %d %s", recorder.Code, recorder.Body.String())
	}
	policy := recorder.Header().Get("Permissions-Policy")
	for _, want := range []string{"camera=(self)", "microphone=(self)"} {
		if !strings.Contains(policy, want) {
			t.Errorf("Permissions-Policy = %q, missing %q", policy, want)
		}
	}
	csp := recorder.Header().Get("Content-Security-Policy-Report-Only")
	for _, want := range []string{"https://chat.example.test", "wss://chat.example.test", "wss://live.example.test:7881"} {
		if !strings.Contains(csp, want) {
			t.Errorf("CSP = %q, missing %q", csp, want)
		}
	}
}

func TestReferenceIsServed(t *testing.T) {
	router := testRouter(t)

	page := httptest.NewRecorder()
	router.ServeHTTP(page, httptest.NewRequest(http.MethodGet, "/docs", nil))
	if page.Code != http.StatusOK {
		t.Fatalf("GET /docs = %d, want 200", page.Code)
	}
	if !strings.Contains(page.Body.String(), `data-url="/docs/openapi.json"`) {
		t.Errorf("the reference page does not point at the spec: %s", page.Body.String())
	}

}

func TestSpecIsServed(t *testing.T) {
	spec := httptest.NewRecorder()
	testRouter(t).ServeHTTP(spec, httptest.NewRequest(http.MethodGet, "/docs/openapi.json", nil))
	if spec.Code != http.StatusOK {
		t.Fatalf("GET /docs/openapi.json = %d, want 200", spec.Code)
	}
	var document struct {
		OpenAPI string `json:"openapi"`
		Info    struct {
			Title string `json:"title"`
		} `json:"info"`
		Paths map[string]any `json:"paths"`
	}
	if err := json.Unmarshal(spec.Body.Bytes(), &document); err != nil {
		t.Fatalf("the spec is not JSON: %v", err)
	}
	if !strings.HasPrefix(document.OpenAPI, "3.1") {
		t.Errorf("openapi = %q, want 3.1.x", document.OpenAPI)
	}
	if document.Info.Title != apiTitle {
		t.Errorf("info.title = %q, want %q", document.Info.Title, apiTitle)
	}
	if len(document.Paths) == 0 {
		t.Error("the spec documents no paths")
	}
}
