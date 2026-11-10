package main

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/FacileStudio/Prisme/apps/api/modules/auth"
	"github.com/FacileStudio/porte"
	"github.com/FacileStudio/porte/directory"
	"github.com/FacileStudio/porte/kit"
)

type noIdentities struct{}

func (noIdentities) ListByUser(context.Context, int64) ([]porte.StoredIdentity, error) {
	return nil, nil
}

// The invitations badge reads the signed-in person's Registre subject, so it
// answers only behind a session, and only where the Registre sync is wired.
func TestTheInvitationsBadgeNeedsASession(t *testing.T) {
	quiet := slog.New(slog.NewTextHandler(io.Discard, nil))
	authService, err := auth.NewService(context.Background(), nil, porte.Config{}, nil, quiet)
	if err != nil {
		t.Fatal(err)
	}
	client := directory.Client{Issuer: "https://sso.test", ClientID: "prisme", ClientSecret: "s"}
	invitations, err := kit.NewInvitations(kit.InvitationsConfig{Directory: client, Identities: noIdentities{}, Logger: quiet})
	if err != nil {
		t.Fatal(err)
	}
	router := buildRouter(app{logger: quiet, auth: authService, directory: directorySync{invitations: invitations}})
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/api/me/invitations", nil))
	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("GET /api/me/invitations without a session = %d, want 401", recorder.Code)
	}
}
