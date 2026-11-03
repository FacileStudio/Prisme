package agents

import (
	"context"
	stderrors "errors"
	"fmt"
	"slices"

	"github.com/FacileStudio/Prisme/apps/api/schemas"
	"github.com/FacileStudio/tronc/errors"
)

func (s *Service) active(ctx context.Context, ownerID int64) ([]schemas.Agent, error) {
	var rows []schemas.Agent
	err := s.deps.DB.WithContext(ctx).Where("owner_id = ? AND deactivated_at IS NULL", ownerID).Order("name").Find(&rows).Error
	if err != nil {
		return nil, errors.Internal("failed to list agents", err)
	}
	return rows, nil
}

// Follow joins every active agent of the owner to the rooms, which the owner
// must already be in, so agents reach a new room without waiting for Bulle.
func (s *Service) Follow(ctx context.Context, ownerID int64, roomIDs []string) error {
	rows, err := s.active(ctx, ownerID)
	if err != nil || len(rows) == 0 {
		return err
	}
	owner, err := s.deps.Accounts.Ensure(ctx, ownerID)
	if err != nil {
		return err
	}
	var failures []error
	for _, row := range rows {
		for _, roomID := range unique(roomIDs) {
			failures = append(failures, s.place(ctx, owner.MatrixUserID, row.MatrixUserID, roomID))
		}
	}
	return stderrors.Join(failures...)
}

// Unfollow makes every active agent of the owner leave those of the rooms it
// is joined to.
func (s *Service) Unfollow(ctx context.Context, ownerID int64, roomIDs []string) error {
	rows, err := s.active(ctx, ownerID)
	if err != nil {
		return err
	}
	var failures []error
	for _, row := range rows {
		joined, err := s.deps.Matrix.JoinedRooms(ctx, row.MatrixUserID)
		if err != nil {
			failures = append(failures, fmt.Errorf("list the rooms of %s: %w", row.MatrixUserID, err))
			continue
		}
		for _, roomID := range unique(roomIDs) {
			if slices.Contains(joined, roomID) {
				failures = append(failures, s.deps.Matrix.Leave(ctx, row.MatrixUserID, roomID))
			}
		}
	}
	return stderrors.Join(failures...)
}
