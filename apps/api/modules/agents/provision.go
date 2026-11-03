package agents

import (
	"context"
	"log/slog"

	"github.com/FacileStudio/Prisme/apps/api/schemas"
	"github.com/FacileStudio/tronc/errors"
)

func (s *Service) provision(ctx context.Context, owner schemas.Account, row *schemas.Agent, localpart string) error {
	if err := s.deps.Matrix.Register(ctx, localpart); err != nil {
		return errors.Internal("failed to register the agent", err)
	}
	*row = schemas.Agent{OwnerID: owner.UserID, Slug: localpart, MatrixUserID: s.deps.Matrix.UserID(localpart)}
	return nil
}

// migrate moves an agent onto the Matrix account named by localpart: it
// registers it, joins it to every room the old account is in, deactivates the
// old account and points the row at the new one. A room the new account cannot
// enter is logged and skipped; the owner can add it back through SetRooms.
func (s *Service) migrate(ctx context.Context, owner schemas.Account, row *schemas.Agent, localpart string) error {
	if err := s.deps.Matrix.Register(ctx, localpart); err != nil {
		return errors.Internal("failed to register the agent's new account", err)
	}
	next := s.deps.Matrix.UserID(localpart)
	rooms, err := s.deps.Matrix.JoinedRooms(ctx, row.MatrixUserID)
	if err != nil {
		return errors.Internal("failed to list the agent's rooms", err)
	}
	for _, roomID := range rooms {
		if s.place(ctx, row.MatrixUserID, next, roomID) == nil {
			continue
		}
		if err := s.place(ctx, owner.MatrixUserID, next, roomID); err != nil {
			s.deps.Logger.Warn("agent migration skipped a room", slog.String("agent", next),
				slog.String("room", roomID), slog.Any("error", err))
		}
	}
	if err := s.deps.Matrix.Deactivate(ctx, row.MatrixUserID); err != nil {
		return errors.Internal("failed to deactivate the agent's old account", err)
	}
	s.deps.Logger.Info("agent migrated", slog.String("from", row.MatrixUserID), slog.String("to", next))
	row.Slug, row.MatrixUserID, row.DeviceID, row.SealedToken = localpart, next, "", ""
	row.AvatarURL, row.AvatarSource = "", ""
	return nil
}

func (s *Service) login(ctx context.Context, row *schemas.Agent) (Login, error) {
	token := ""
	if openToken, err := s.deps.Box.Open(row.SealedToken); err == nil && s.deps.Matrix.WhoAmI(ctx, openToken) {
		token = openToken
	} else {
		login, err := s.deps.Matrix.Login(ctx, row.MatrixUserID, "Bulle")
		if err != nil {
			return Login{}, errors.Internal("failed to sign the agent in", err)
		}
		row.DeviceID, row.SealedToken, token = login.DeviceID, s.deps.Box.Seal(login.AccessToken), login.AccessToken
	}
	if err := s.deps.DB.WithContext(ctx).Save(row).Error; err != nil {
		return Login{}, errors.Internal("failed to store the agent", err)
	}
	return Login{Homeserver: s.deps.Homeserver, UserID: row.MatrixUserID, AccessToken: token, DeviceID: row.DeviceID}, nil
}
