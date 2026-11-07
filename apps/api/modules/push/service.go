// Package push wakes the Android app when someone writes to its user. The app
// hands its Firebase token to POST /push/devices. Synapse gives the appservice
// every event; for each new message the notifier finds who it is for (the
// other person in a direct message, and anyone it mentions), skips muted
// rooms, and sends to their devices through Firebase. Synapse cannot do this
// itself: it never runs push rules for users an appservice owns, and every
// Prisme user is one. The desktop app and the browser notify from their own
// live sync instead.
package push

import (
	"context"
	"crypto/subtle"
	"log/slog"
	"strings"

	"github.com/FacileStudio/Prisme/apps/api/internal/matrix"
	"github.com/FacileStudio/tronc/errors"

	"gorm.io/gorm"
)

const maxTokenBytes = 4096

// Deps is what the push service is built from. A nil FCM or an empty HSToken
// turns Android push off.
type Deps struct {
	DB      *gorm.DB
	Matrix  *matrix.Client
	FCM     *FCM
	HSToken string
	Logger  *slog.Logger
}

// Service stores Android devices and notifies them of new messages.
type Service struct {
	deps   Deps
	sender sender
	people people
}

type sender interface {
	send(ctx context.Context, token string, m message) error
}

type people interface {
	JoinedMembers(ctx context.Context, asUser, roomID string) ([]string, error)
	RoomMuted(ctx context.Context, userID, roomID string) bool
	DisplayName(ctx context.Context, userID string) string
	RoomName(ctx context.Context, asUser, roomID string) string
}

// NewService builds the push service.
func NewService(deps Deps) *Service {
	service := &Service{deps: deps, people: deps.Matrix}
	if deps.FCM != nil {
		service.sender = deps.FCM
	}
	return service
}

// Config says whether this server can push to the Android app.
func (s *Service) Config() Config {
	return Config{Native: s.sender != nil && s.deps.HSToken != ""}
}

// Register stores the Firebase token of this porte session's Android app,
// moving it over when the phone was signed in as someone else.
func (s *Service) Register(ctx context.Context, userID, porteSessionID int64, token string) error {
	if !s.Config().Native {
		return errors.Unavailable("push is not configured on this server")
	}
	token = strings.TrimSpace(token)
	if token == "" || len(token) > maxTokenBytes {
		return errors.Invalid("token must be a Firebase registration token")
	}
	err := s.deps.DB.WithContext(ctx).Exec(`INSERT INTO push_devices (token, user_id, porte_session_id) VALUES (?, ?, ?)
		ON CONFLICT (token) DO UPDATE SET user_id = EXCLUDED.user_id, porte_session_id = EXCLUDED.porte_session_id`,
		token, userID, porteSessionID).Error
	if err != nil {
		return errors.Internal("failed to store the device", err)
	}
	return nil
}

// Authorized reports whether a request carries Synapse's hs_token.
func (s *Service) Authorized(header string) bool {
	token, ok := strings.CutPrefix(header, "Bearer ")
	return ok && s.deps.HSToken != "" && subtle.ConstantTimeCompare([]byte(token), []byte(s.deps.HSToken)) == 1
}

func (s *Service) devices(ctx context.Context, matrixUserIDs []string) []device {
	var rows []device
	err := s.deps.DB.WithContext(ctx).Raw(`SELECT d.token, a.matrix_user_id FROM push_devices d
		JOIN accounts a ON a.user_id = d.user_id
		JOIN porte_sessions p ON p.id = d.porte_session_id AND (p.expires_at IS NULL OR p.expires_at > now())
		WHERE a.matrix_user_id IN ?`, matrixUserIDs).Scan(&rows).Error
	if err != nil {
		s.deps.Logger.Warn("failed to read push devices", slog.Any("error", err))
	}
	return rows
}

func (s *Service) forget(ctx context.Context, token string) {
	if err := s.deps.DB.WithContext(ctx).Exec(`DELETE FROM push_devices WHERE token = ?`, token).Error; err != nil {
		s.deps.Logger.Warn("failed to forget a push device", slog.Any("error", err))
	}
}
