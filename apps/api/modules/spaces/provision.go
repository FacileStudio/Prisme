package spaces

import (
	"context"
	"log/slog"

	"github.com/FacileStudio/Prisme/apps/api/schemas"
	"github.com/FacileStudio/tronc/errors"
)

// provision creates a space and its #general channel, owned by the account.
// A non-nil registreID links the space to that Registre space.
func (s *Service) provision(ctx context.Context, account schemas.Account, name string, registreID *string) (Space, error) {
	row := schemas.Space{Name: name, CreatedBy: account.UserID, Slug: s.freeSlug(ctx, Slugify(name, "space")), RegistreSpaceID: registreID}
	roomID, err := s.createRoom(ctx, account.MatrixUserID, spaceRoom(name, s.bot()))
	if err != nil {
		return Space{}, err
	}
	row.RoomID = roomID
	if err := s.db.WithContext(ctx).Create(&row).Error; err != nil {
		return Space{}, errors.Internal("failed to store the space", err)
	}
	general, err := s.addChannel(ctx, row, account.MatrixUserID, CreateChannelRequest{Name: "general"})
	if err != nil {
		return Space{}, err
	}
	return toSpace(row, []Channel{general}), nil
}

// follow brings the user's agents into the rooms. It logs a failure rather
// than returning it: Bulle's next provisioning pass repairs what is missing.
func (s *Service) follow(ctx context.Context, userID int64, roomIDs []string) {
	if s.followers == nil {
		return
	}
	if err := s.followers.Follow(ctx, userID, roomIDs); err != nil {
		s.logger.Warn("agents could not follow their owner", slog.Int64("owner", userID), slog.Any("error", err))
	}
}

// unfollow takes the user's agents out of the rooms, logging a failure.
func (s *Service) unfollow(ctx context.Context, userID int64, roomIDs []string) {
	if s.followers == nil {
		return
	}
	if err := s.followers.Unfollow(ctx, userID, roomIDs); err != nil {
		s.logger.Warn("agents could not leave with their owner", slog.Int64("owner", userID), slog.Any("error", err))
	}
}

func toSpace(row schemas.Space, channels []Channel) Space {
	return Space{
		ID: row.ID, RoomID: row.RoomID, Name: row.Name, Slug: row.Slug,
		Managed: row.RegistreSpaceID != nil, Channels: channels,
	}
}

func roomsOf(space Space) []string {
	rooms := []string{space.RoomID}
	for _, channel := range space.Channels {
		rooms = append(rooms, channel.RoomID)
	}
	return rooms
}
