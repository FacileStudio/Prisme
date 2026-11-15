package capsule

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func fakeCapsule(t *testing.T) *httptest.Server {
	t.Helper()
	burned := false
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/pastes", func(w http.ResponseWriter, r *http.Request) {
		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()
		var req createRequest
		if err := dec.Decode(&req); err != nil || !req.BurnAfterRead || req.Content != "cipher" || req.ExpiresIn != "24h" {
			http.Error(w, "bad", http.StatusBadRequest)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"id": "cap_0123456789abcdef", "delete_token": "secret"})
	})
	mux.HandleFunc("GET /api/pastes/{id}", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{"id": r.PathValue("id"), "exists": !burned})
	})
	mux.HandleFunc("POST /api/pastes/{id}/content", func(w http.ResponseWriter, _ *http.Request) {
		if burned {
			http.Error(w, `{"error":"not_found"}`, http.StatusNotFound)
			return
		}
		burned = true
		_ = json.NewEncoder(w).Encode(map[string]string{"content": "cipher"})
	})
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)
	return server
}

func TestSealOpenBurns(t *testing.T) {
	client := New(fakeCapsule(t).URL)
	ctx := context.Background()
	sealed, err := client.Seal(ctx, "cipher", "24h")
	if err != nil || sealed.ID != "cap_0123456789abcdef" {
		t.Fatalf("seal: %v %+v", err, sealed)
	}
	if meta, err := client.Meta(ctx, sealed.ID); err != nil || !meta.Exists {
		t.Fatalf("meta before open: %v %+v", err, meta)
	}
	if content, err := client.Open(ctx, sealed.ID); err != nil || content != "cipher" {
		t.Fatalf("open: %v %q", err, content)
	}
	if _, err := client.Open(ctx, sealed.ID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("second open: want ErrNotFound, got %v", err)
	}
	if meta, _ := client.Meta(ctx, sealed.ID); meta.Exists {
		t.Fatal("meta after open still exists")
	}
}

func TestValidation(t *testing.T) {
	for id, want := range map[string]bool{"cap_0123456789abcdef": true, "cap_0123": false, "../x": false, "cap_0123456789ABCDEF": false} {
		if ValidID(id) != want {
			t.Errorf("ValidID(%q) = %v", id, !want)
		}
	}
	if !ValidExpiry("7d") || ValidExpiry("2d") {
		t.Error("ValidExpiry")
	}
}
