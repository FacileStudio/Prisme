package nuage

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

type fakeNuage struct {
	tokens  []string
	created int
	reject  string
	lasting bool
}

func (f *fakeNuage) handler(t *testing.T) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/files/7", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, File{ID: 7, Name: "plan.pdf", MimeType: "application/pdf", Size: 2048})
	})
	mux.HandleFunc("GET /api/shares/by-me", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, map[string]any{"shares": f.shares()})
	})
	mux.HandleFunc("POST /api/shares", func(w http.ResponseWriter, r *http.Request) {
		var body shareRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.FileID != 7 || body.Permission != "view" {
			t.Errorf("share request %+v, %v", body, err)
		}
		f.created++
		writeJSON(w, share{Token: "fresh", FileID: &body.FileID, Permission: "view"})
	})
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		f.tokens = append(f.tokens, token)
		if token == f.reject {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		mux.ServeHTTP(w, r)
	})
}

func (f *fakeNuage) shares() []share {
	expires := "2026-12-01T00:00:00Z"
	other, seven := int64(9), int64(7)
	shares := []share{
		{Token: "other", FileID: &other, Permission: "view"},
		{Token: "brief", FileID: &seven, Permission: "view", ExpiresAt: &expires},
	}
	if f.lasting {
		shares = append(shares, share{Token: "kept", FileID: &seven, Permission: "view"})
	}
	return shares
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}

func clientFor(t *testing.T, fake *fakeNuage) (*Client, string) {
	server := httptest.NewServer(fake.handler(t))
	t.Cleanup(server.Close)
	next := []string{"stale", "fresh"}
	forgot := 0
	return New(server.URL+"/", Tokens{
		Token:  func(context.Context, int64) (string, error) { return next[min(forgot, len(next)-1)], nil },
		Forget: func(int64) { forgot++ },
	}), server.URL
}

func TestShareReusesALastingViewLinkOrMakesOne(t *testing.T) {
	for _, c := range []struct {
		lasting bool
		token   string
		created int
	}{{false, "fresh", 1}, {true, "kept", 0}} {
		fake := &fakeNuage{reject: "Bearer stale", lasting: c.lasting}
		client, base := clientFor(t, fake)
		shared, err := client.Share(context.Background(), 1, 7)
		if err != nil {
			t.Fatal(err)
		}
		if fake.created != c.created || shared.Name != "plan.pdf" || shared.URL != base+"/s/"+c.token {
			t.Errorf("lasting=%v: shared %+v after %d creations", c.lasting, shared, fake.created)
		}
	}
}

func TestARefusedTokenIsForgottenAndRetriedOnce(t *testing.T) {
	fake := &fakeNuage{reject: "Bearer stale"}
	client, _ := clientFor(t, fake)
	if _, err := client.Spaces(context.Background(), 1); err != nil && !errors.Is(err, ErrNotFound) {
		t.Fatalf("spaces after a stale token: %v", err)
	}
	if len(fake.tokens) != 2 || fake.tokens[1] != "Bearer fresh" {
		t.Errorf("tokens sent %v, want the stale one then a fresh one", fake.tokens)
	}
}

func TestATokenNuageKeepsRefusingIsReportedAsRefused(t *testing.T) {
	fake := &fakeNuage{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fake.tokens = append(fake.tokens, r.Header.Get("Authorization"))
		w.WriteHeader(http.StatusForbidden)
	}))
	t.Cleanup(server.Close)
	client := New(server.URL, Tokens{
		Token:  func(context.Context, int64) (string, error) { return "t", nil },
		Forget: func(int64) {},
	})
	if _, _, err := client.List(context.Background(), 1, Place{}); !errors.Is(err, ErrRefused) {
		t.Errorf("list with a forbidden actor: %v, want ErrRefused", err)
	}
}
