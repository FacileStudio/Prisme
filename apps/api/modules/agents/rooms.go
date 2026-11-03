package agents

import (
	"context"
	stderrors "errors"
	"slices"
	"strings"
	"time"

	"github.com/FacileStudio/Prisme/apps/api/schemas"
	"github.com/FacileStudio/tronc/errors"

	"gorm.io/gorm"
)

// List returns the owner's active agents and the rooms each is joined to.
func (s *Service) List(ctx context.Context, ownerID int64) ([]Agent, error) {
	rows, err := s.active(ctx, ownerID)
	if err != nil {
		return nil, err
	}
	out := make([]Agent, 0, len(rows))
	for _, row := range rows {
		rooms, err := s.deps.Matrix.JoinedRooms(ctx, row.MatrixUserID)
		if err != nil {
			rooms = []string{}
		}
		out = append(out, Agent{UserID: row.MatrixUserID, Name: row.Name, RoomIDs: rooms})
	}
	return out, nil
}

// SetRooms makes the agent a member of exactly the listed channels. Every
// room must be one the owner is joined to; the owner invites the agent, and
// the agent joins through the appservice. Only Prisme's spaces and channels
// are ever left: the owner's direct messages with the agent are kept, and
// rejoined when the agent is missing from one.
func (s *Service) SetRooms(ctx context.Context, ownerID int64, userID string, roomIDs []string) (RoomsResponse, error) {
	row, err := s.owned(ctx, ownerID, userID)
	if err != nil {
		return RoomsResponse{}, err
	}
	owner, ownerRooms, err := s.ownerOf(ctx, ownerID, roomIDs)
	if err != nil {
		return RoomsResponse{}, err
	}
	current, err := s.deps.Matrix.JoinedRooms(ctx, row.MatrixUserID)
	if err != nil {
		return RoomsResponse{}, errors.Internal("failed to list the agent's rooms", err)
	}
	wanted := append(unique(roomIDs), s.directs(ctx, owner.MatrixUserID, row.MatrixUserID, ownerRooms)...)
	for _, roomID := range missing(wanted, current) {
		if err := s.place(ctx, owner.MatrixUserID, row.MatrixUserID, roomID); err != nil {
			return RoomsResponse{}, err
		}
	}
	leaving, err := s.managed(ctx, missing(current, wanted))
	if err != nil {
		return RoomsResponse{}, err
	}
	for _, roomID := range leaving {
		if err := s.deps.Matrix.Leave(ctx, row.MatrixUserID, roomID); err != nil {
			return RoomsResponse{}, errors.Internal("the agent could not leave "+roomID, err)
		}
	}
	return RoomsResponse{RoomIDs: unique(roomIDs)}, nil
}

func (s *Service) ownerOf(ctx context.Context, ownerID int64, roomIDs []string) (schemas.Account, []string, error) {
	owner, err := s.deps.Accounts.Ensure(ctx, ownerID)
	if err != nil {
		return owner, nil, err
	}
	ownerRooms, err := s.deps.Matrix.JoinedRooms(ctx, owner.MatrixUserID)
	if err != nil {
		return owner, nil, errors.Internal("failed to list the owner's rooms", err)
	}
	if outside := missing(roomIDs, ownerRooms); len(outside) > 0 {
		return owner, nil, errors.Forbidden("not a member of: " + strings.Join(outside, ", "))
	}
	return owner, ownerRooms, nil
}

func (s *Service) place(ctx context.Context, owner, agent, roomID string) error {
	if err := s.deps.Matrix.Invite(ctx, owner, roomID, agent); err != nil {
		return errors.Internal("could not invite the agent to "+roomID, err)
	}
	if err := s.deps.Matrix.Join(ctx, agent, roomID); err != nil {
		return errors.Internal("the agent could not join "+roomID, err)
	}
	return nil
}

func (s *Service) owned(ctx context.Context, ownerID int64, userID string) (schemas.Agent, error) {
	var row schemas.Agent
	err := s.deps.DB.WithContext(ctx).
		First(&row, "owner_id = ? AND matrix_user_id = ? AND deactivated_at IS NULL", ownerID, userID).Error
	if stderrors.Is(err, gorm.ErrRecordNotFound) {
		return row, errors.NotFound("no such agent")
	}
	if err != nil {
		return row, errors.Internal("failed to read the agent", err)
	}
	return row, nil
}

// Deactivate closes one of the owner's agents for good.
func (s *Service) Deactivate(ctx context.Context, ownerID int64, userID string) error {
	row, err := s.owned(ctx, ownerID, userID)
	if err != nil {
		return err
	}
	if err := s.deps.Matrix.Deactivate(ctx, row.MatrixUserID); err != nil {
		return errors.Internal("failed to deactivate the agent", err)
	}
	now := time.Now()
	row.DeactivatedAt, row.SealedToken, row.DeviceID = &now, "", ""
	return s.deps.DB.WithContext(ctx).Save(&row).Error
}

func missing(want, have []string) []string {
	var out []string
	for _, id := range unique(want) {
		if !slices.Contains(have, id) {
			out = append(out, id)
		}
	}
	return out
}

func unique(ids []string) []string {
	out := []string{}
	for _, id := range ids {
		if id != "" && !slices.Contains(out, id) {
			out = append(out, id)
		}
	}
	return out
}
