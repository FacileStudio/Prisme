package main

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/FacileStudio/Prisme/apps/api/modules/auth"
	"github.com/FacileStudio/porte"
)

var testBudgets = budgets{
	visitor: budget{limit: 50, burst: 2},
	write:   budget{limit: 50, burst: 1},
	signIn:  budget{limit: 50, burst: 1},
	service: budget{limit: 50, burst: 3},
}

func limitedHandler(t *testing.T) http.Handler {
	t.Helper()
	quiet := slog.New(slog.NewTextHandler(io.Discard, nil))
	authService, err := auth.NewService(context.Background(), nil, porte.Config{}, auth.Guards{ServiceKeys: auth.ParseServiceKeys("key-one")}, quiet)
	if err != nil {
		t.Fatal(err)
	}
	ok := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusNoContent) })
	return newRateLimits(testBudgets, authService).requests(ok)
}

func hit(handler http.Handler, method, addr, key string) *httptest.ResponseRecorder {
	request := httptest.NewRequest(method, "/api/spaces", nil)
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
