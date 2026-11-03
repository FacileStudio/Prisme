// Package calls starts Echo calls from Matrix rooms: one Echo room per Matrix
// room, a LiveKit token per join, and a call message in the room.
package calls

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"strings"

	"github.com/FacileStudio/Prisme/apps/api/internal/accounts"
	"github.com/FacileStudio/Prisme/apps/api/internal/echo"
	"github.com/FacileStudio/Prisme/apps/api/internal/matrix"
	"github.com/FacileStudio/tronc/errors"
)

// ErrNotConfigured is the answer when the deployment has no Echo.
const ErrNotConfigured = "Calls are not set up on this Prisme: ECHO_URL and ECHO_INTEGRATION_KEY are unset."

// Service starts and joins calls. A nil Echo client turns calls off.
type Service struct {
	echo     *echo.Client
	matrix   *matrix.Client
	accounts *accounts.Service
}

// NewService builds the call service.
func NewService(echoClient *echo.Client, client *matrix.Client, accountService *accounts.Service) *Service {
	return &Service{echo: echoClient, matrix: client, accounts: accountService}
}

// Start gets or creates the Echo room of a Matrix room the user is in, mints
// the user a token and, unless joining, announces the call in the room.
func (s *Service) Start(ctx context.Context, userID int64, req Request) (Response, error) {
	if s.echo == nil {
		return Response{}, errors.Unavailable(ErrNotConfigured)
	}
	account, err := s.accounts.Ensure(ctx, userID)
	if err != nil {
		return Response{}, err
	}
	if req.RoomID == "" || !s.matrix.IsJoined(ctx, account.MatrixUserID, req.RoomID) {
		return Response{}, errors.NotFound("no such room")
	}
	profile, err := s.accounts.Profile(ctx, userID)
	if err != nil {
		return Response{}, err
	}
	room, err := s.echo.Room(ctx, req.RoomID, s.title(ctx, account.MatrixUserID, req.RoomID))
	if err != nil {
		return Response{}, errors.New("unavailable", "Echo did not open the call", err)
	}
	name := displayName(profile, account.Localpart)
	token, err := s.echo.Token(ctx, room.RoomID, echo.Participant{Identity: account.MatrixUserID + "#" + suffix(), Name: name})
	if err != nil {
		return Response{}, errors.New("unavailable", "Echo did not issue a join token", err)
	}
	if !req.Join {
		s.announce(ctx, account.MatrixUserID, req.RoomID, name, room)
	}
	return Response{CallID: room.RoomID, LiveKitURL: token.LiveKitURL, Token: token.Token, JoinURL: room.JoinURL}, nil
}

func (s *Service) title(ctx context.Context, userID, roomID string) string {
	if name := s.matrix.RoomName(ctx, userID, roomID); name != "" {
		return name
	}
	return "Prisme call"
}

func (s *Service) announce(ctx context.Context, userID, roomID, name string, room echo.Room) {
	content := map[string]any{
		"msgtype":            "m.text",
		"body":               name + " started a call: " + room.JoinURL,
		"studio.facile.call": CallContent{CallID: room.RoomID, JoinURL: room.JoinURL, StartedBy: userID},
	}
	_, _ = s.matrix.SendMessage(ctx, userID, roomID, content)
}

func displayName(profile accounts.Profile, fallback string) string {
	if name := strings.TrimSpace(profile.Name); name != "" {
		return name
	}
	return fallback
}

func suffix() string {
	random := make([]byte, 4)
	_, _ = rand.Read(random)
	return hex.EncodeToString(random)
}
