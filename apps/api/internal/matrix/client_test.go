package matrix

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCallsAssertTheUserAndTreatAnExistingUserAsRegistered(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer as" {
			t.Errorf("Authorization = %q", r.Header.Get("Authorization"))
		}
		switch r.URL.Path {
		case "/_matrix/client/v3/register":
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"errcode":"M_USER_IN_USE","error":"taken"}`))
		case "/_matrix/client/v3/joined_rooms":
			if r.URL.Query().Get("user_id") != "@a:hs" {
				t.Errorf("user_id = %q", r.URL.Query().Get("user_id"))
			}
			_, _ = w.Write([]byte(`{"joined_rooms":["!r:hs"]}`))
		}
	}))
	defer server.Close()
	client := New(Config{BaseURL: server.URL, ServerName: "hs", ASToken: "as"})
	if err := client.Register(context.Background(), "a"); err != nil {
		t.Fatalf("Register = %v, want nil for M_USER_IN_USE", err)
	}
	if !client.IsJoined(context.Background(), client.UserID("a"), "!r:hs") {
		t.Fatal("IsJoined = false")
	}
}
