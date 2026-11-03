package agents

import (
	"context"
	"slices"

	"github.com/FacileStudio/Prisme/apps/api/schemas"
	"github.com/FacileStudio/tronc/errors"
)

// directs lists the owner's direct messages with the agent that the owner is
// still in. A failed read keeps nothing extra rather than failing the call.
func (s *Service) directs(ctx context.Context, owner, agent string, ownerRooms []string) []string {
	direct, err := s.deps.Matrix.DirectRooms(ctx, owner)
	if err != nil {
		s.deps.Logger.Warn("agents: reading the owner's direct messages", "owner", owner, "error", err)
		return nil
	}
	var out []string
	for _, roomID := range unique(direct[agent]) {
		if slices.Contains(ownerRooms, roomID) {
			out = append(out, roomID)
		}
	}
	return out
}

// managed keeps the rooms Prisme created as spaces or channels, the only
// rooms SetRooms may make an agent leave.
func (s *Service) managed(ctx context.Context, roomIDs []string) ([]string, error) {
	if len(roomIDs) == 0 {
		return nil, nil
	}
	var known []string
	err := s.deps.DB.WithContext(ctx).Model(&schemas.Channel{}).Where("room_id IN ?", roomIDs).Pluck("room_id", &known).Error
	if err == nil {
		var spaces []string
		err = s.deps.DB.WithContext(ctx).Model(&schemas.Space{}).Where("room_id IN ?", roomIDs).Pluck("room_id", &spaces).Error
		known = append(known, spaces...)
	}
	if err != nil {
		return nil, errors.Internal("failed to read Prisme's rooms", err)
	}
	return known, nil
}
