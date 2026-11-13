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
	"time"

	"github.com/FacileStudio/Prisme/apps/api/modules/auth"
	"github.com/FacileStudio/porte"
	"github.com/go-chi/chi/v5"
)

var testBudgets = budgets{
	visitor: budget{limit: 50, burst: 2},
	write:   budget{limit: 50, burst: 1},
	signIn:  budget{limit: 50, burst: 1},
	service: budget{limit: 50, burst: 3},
}

// sessionCallers is auth.Service with a session table of bearer tokens, so
// the limiter's tests need no database: auth's own tests cover SessionUser.
type sessionCallers struct {
	*auth.Service
	sessions map[string]int64
	lookups  map[string]int
}

func (c sessionCallers) SessionUser(r *http.Request) (int64, bool) {
	c.lookups[r.URL.Path]++
	userID, ok := c.sessions[strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")]
	return userID, ok
}

func limitedHandler(t *testing.T) http.Handler {
	t.Helper()
	handler, _ := countedHandler(t)
	return handler
}

func countedHandler(t *testing.T) (http.Handler, map[string]int) {
	t.Helper()
	quiet := slog.New(slog.NewTextHandler(io.Discard, nil))
	authService, err := auth.NewService(context.Background(), nil, porte.Config{}, auth.Guards{ServiceKeys: auth.ParseServiceKeys("key-one")}, quiet)
	if err != nil {
		t.Fatal(err)
	}
	who := sessionCallers{Service: authService, sessions: map[string]int64{"ada-session": 1, "bob-session": 2}, lookups: map[string]int{}}
	ok := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusNoContent) })
	limits := newRateLimits(testBudgets, who)
	router := chi.NewRouter()
	router.Route("/api", func(r chi.Router) {
		r.Use(limits.requests)
		r.With(limits.signIn).Handle("/auth/*", ok)
		r.Handle("/*", ok)
	})
	return router, who.lookups
}

func hit(handler http.Handler, method, addr, key string) *httptest.ResponseRecorder {
	return hitPath(handler, method, "/api/spaces", addr, key)
}

func hitPath(handler http.Handler, method, path, addr, key string) *httptest.ResponseRecorder {
	request := httptest.NewRequest(method, path, nil)
	request.RemoteAddr = addr + ":40000"
	if key != "" {
		request.Header.Set("Authorization", "Bearer "+key)
	}
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	return recorder
}

func expectRefused(t *testing.T, recorder *httptest.ResponseRecorder) {
	t.Helper()
	if recorder.Code != http.StatusTooManyRequests {
		t.Fatalf("status %d, want 429", recorder.Code)
	}
	if recorder.Header().Get("Retry-After") == "" {
		t.Fatal("a refusal carries no Retry-After")
	}
	var body struct {
		Error struct {
			Code string `json:"code"`
		} `json:"error"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil || body.Error.Code != "rate_limited" {
		t.Fatalf("body %q, want the rate_limited envelope", recorder.Body.String())
	}
}

func expectPassed(t *testing.T, recorder *httptest.ResponseRecorder) {
	t.Helper()
	if recorder.Code != http.StatusNoContent {
		t.Fatalf("status %d, want the request through", recorder.Code)
	}
}

func TestVisitorLimitTripsAndResets(t *testing.T) {
	handler := limitedHandler(t)
	expectPassed(t, hit(handler, http.MethodGet, "198.51.100.7", ""))
	expectPassed(t, hit(handler, http.MethodGet, "198.51.100.7", ""))
	expectRefused(t, hit(handler, http.MethodGet, "198.51.100.7", ""))
	expectPassed(t, hit(handler, http.MethodGet, "198.51.100.8", ""))
	time.Sleep(60 * time.Millisecond)
	expectPassed(t, hit(handler, http.MethodGet, "198.51.100.7", ""))
}

func TestWritesSpendTheirOwnBucket(t *testing.T) {
	handler := limitedHandler(t)
	expectPassed(t, hit(handler, http.MethodPost, "198.51.100.7", ""))
	expectRefused(t, hit(handler, http.MethodPost, "198.51.100.7", ""))
	expectPassed(t, hit(handler, http.MethodPost, "198.51.100.8", ""))
}

func TestServiceKeyHasItsOwnBucket(t *testing.T) {
	handler := limitedHandler(t)
	for range 2 {
		expectPassed(t, hit(handler, http.MethodGet, "198.51.100.7", ""))
	}
	expectRefused(t, hit(handler, http.MethodGet, "198.51.100.7", ""))
	for range 3 {
		expectPassed(t, hit(handler, http.MethodGet, "198.51.100.7", "key-one"))
	}
	expectRefused(t, hit(handler, http.MethodGet, "198.51.100.7", "key-one"))
	expectRefused(t, hit(handler, http.MethodGet, "198.51.100.7", "not-a-key"))
	expectPassed(t, hit(handler, http.MethodGet, "198.51.100.9", ""))
}

func TestSignInRoutesAreLimited(t *testing.T) {
	router := testRouter(t)
	last := 0
	for range productionBudgets.signIn.burst + 1 {
		request := httptest.NewRequest(http.MethodGet, "/api/auth/config", nil)
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, request)
		last = recorder.Code
	}
	if last != http.StatusTooManyRequests {
		t.Fatalf("request %d to /api/auth/config answered %d, want 429", productionBudgets.signIn.burst+1, last)
	}
}

// TestSignedInUsersBehindOneAddressHaveTheirOwnBudget is the office case: two
// people and an anonymous visitor share one address, and each spends from a
// bucket of their own.
func TestSignedInUsersBehindOneAddressHaveTheirOwnBudget(t *testing.T) {
	handler := limitedHandler(t)
	for range 2 {
		expectPassed(t, hit(handler, http.MethodGet, "198.51.100.7", "ada-session"))
	}
	expectRefused(t, hit(handler, http.MethodGet, "198.51.100.7", "ada-session"))
	expectRefused(t, hit(handler, http.MethodGet, "203.0.113.4", "ada-session"))
	for range 2 {
		expectPassed(t, hit(handler, http.MethodGet, "198.51.100.7", "bob-session"))
	}
	expectRefused(t, hit(handler, http.MethodGet, "198.51.100.7", "bob-session"))
	for range 2 {
		expectPassed(t, hit(handler, http.MethodGet, "198.51.100.7", ""))
	}
	expectRefused(t, hit(handler, http.MethodGet, "198.51.100.7", "not-a-session"))
}

// TestAUserSpendsTheirOwnWriteBudget checks writes follow the user as well.
func TestAUserSpendsTheirOwnWriteBudget(t *testing.T) {
	handler := limitedHandler(t)
	expectPassed(t, hit(handler, http.MethodPost, "198.51.100.7", "ada-session"))
	expectRefused(t, hit(handler, http.MethodPost, "198.51.100.7", "ada-session"))
	expectPassed(t, hit(handler, http.MethodPost, "198.51.100.7", "bob-session"))
	expectPassed(t, hit(handler, http.MethodPost, "198.51.100.7", ""))
}

// TestSignInRoutesStayPerAddress checks porte's routes ignore the session:
// they never look it up, and two people behind one address share that
// address's buckets there.
func TestSignInRoutesStayPerAddress(t *testing.T) {
	handler, lookups := countedHandler(t)
	expectPassed(t, hitPath(handler, http.MethodGet, "/api/auth/config", "198.51.100.7", "ada-session"))
	expectRefused(t, hitPath(handler, http.MethodGet, "/api/auth/config", "198.51.100.7", "bob-session"))
	expectPassed(t, hitPath(handler, http.MethodGet, "/api/auth/config", "198.51.100.8", "bob-session"))
	if lookups["/api/auth/config"] != 0 {
		t.Fatalf("porte's routes looked the session up %d times, want none", lookups["/api/auth/config"])
	}
	expectPassed(t, hit(handler, http.MethodGet, "198.51.100.7", "ada-session"))
	if lookups["/api/spaces"] != 1 {
		t.Fatalf("a request looked its session up %d times, want once", lookups["/api/spaces"])
	}
}
