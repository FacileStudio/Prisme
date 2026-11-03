// Package calls starts Echo calls from Matrix rooms: one Echo room per Matrix
// room, a LiveKit token per join, a call message in the room, and an edit of
// that message when Echo reports the call over.
package calls

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"strings"
	"time"

	"github.com/FacileStudio/Prisme/apps/api/internal/accounts"
	"github.com/FacileStudio/Prisme/apps/api/internal/echo"
	"github.com/FacileStudio/Prisme/apps/api/internal/matrix"
	"github.com/FacileStudio/Prisme/apps/api/schemas"
	"github.com/FacileStudio/tronc/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ErrNotConfigured is the answer when the deployment has no Echo.
const ErrNotConfigured = "Calls are not set up on this Prisme: ECHO_URL and ECHO_INTEGRATION_KEY are unset."

// EventsPath is where Echo posts call events, under the API prefix.
const EventsPath = "/echo/events"

// Deps is what the call service needs. A nil Echo turns calls off.
// PublicURL is where Echo reaches Prisme; Secret is PRISME_SECRET_KEY, from
// which every room's call encryption key derives.
type Deps struct {
	DB        *gorm.DB
	Echo      *echo.Client
	Matrix    *matrix.Client
	Accounts  *accounts.Service
	Secret    string
	PublicURL string
	Logger    *slog.Logger
}

// Service starts, joins and ends calls.
type Service struct {
	db       *gorm.DB
	echo     *echo.Client
	matrix   *matrix.Client
	accounts *accounts.Service
	keys     []byte
	callback string
	log      *slog.Logger
}

// NewService builds the call service.
func NewService(deps Deps) *Service {
	callback := ""
	if deps.PublicURL != "" {
		callback = deps.PublicURL + "/api" + EventsPath
	}
	return &Service{
		db: deps.DB, echo: deps.Echo, matrix: deps.Matrix, accounts: deps.Accounts,
		keys: keyMaterial(deps.Secret), callback: callback, log: deps.Logger,
	}
}

// Start gets or creates the Echo room of a Matrix room the user is in, mints
// the user a token and announces the call in the room, unless joining or a
// call is already running there. The room's encryption key goes only to
// members and never to Echo.
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
	room, err := s.echo.Room(ctx, req.RoomID, s.title(ctx, account.MatrixUserID, req.RoomID), s.callback)
	if err != nil {
		return Response{}, errors.New("unavailable", "Echo did not open the call", err)
	}
	name := displayName(profile, account.Localpart)
	token, err := s.echo.Token(ctx, room.RoomID, echo.Participant{Identity: account.MatrixUserID + "#" + suffix(), Name: name})
	if err != nil {
		return Response{}, errors.New("unavailable", "Echo did not issue a join token", err)
	}
	if !req.Join && !s.running(ctx, req.RoomID, room.RoomID) {
		s.announce(ctx, account.MatrixUserID, req.RoomID, name, room)
	}
	return Response{
		CallID: room.RoomID, LiveKitURL: token.LiveKitURL, Token: token.Token, JoinURL: room.JoinURL,
		E2EEKey: roomKey(s.keys, req.RoomID),
	}, nil
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
		"body":               name + " started an end-to-end encrypted call in Prisme",
		"studio.facile.call": CallContent{CallID: room.RoomID, JoinURL: room.JoinURL, StartedBy: userID},
	}
	eventID, err := s.matrix.SendMessage(ctx, userID, roomID, content)
	if err != nil || eventID == "" {
		return
	}
	record := schemas.Call{RoomID: roomID, CallID: room.RoomID, EventID: eventID, StartedBy: userID, StartedAt: time.Now()}
	err = s.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Create(&record).Error
	if err != nil {
		s.log.Warn("calls: recording the call", slog.String("room_id", roomID), slog.Any("error", err))
	}
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
