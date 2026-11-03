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
	"github.com/FacileStudio/porte"
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
	if missing := apiref.Undocumented(testRouter(t), refConfig(), "/api/auth", "/.well-known"); len(missing) > 0 {
		t.Errorf("routes missing from the API registry:\n  %s", strings.Join(missing, "\n  "))
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
