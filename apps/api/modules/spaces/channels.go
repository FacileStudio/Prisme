package spaces

import (
	"context"
	"strings"

	"github.com/FacileStudio/Prisme/apps/api/internal/accounts"
	"github.com/FacileStudio/Prisme/apps/api/internal/matrix"
	"github.com/FacileStudio/Prisme/apps/api/schemas"
	"github.com/FacileStudio/tronc/errors"
)

// CreateChannel adds a channel to a space the user is joined to.
func (s *Service) CreateChannel(ctx context.Context, userID, spaceID int64, req CreateChannelRequest) (Channel, error) {
	req.Name = strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(req.Name), "#"))
	if req.Name == "" || len(req.Name) > 80 {
		return Channel{}, errors.Invalid("a channel needs a name of 1 to 80 characters")
	}
	space, account, err := s.Get(ctx, userID, spaceID)
	if err != nil {
		return Channel{}, err
	}
	return s.addChannel(ctx, space, account.MatrixUserID, req)
}

func (s *Service) addChannel(ctx context.Context, space schemas.Space, creator string, req CreateChannelRequest) (Channel, error) {
	localAlias := s.freeAlias(ctx, "prisme-"+space.Slug+"-"+Slugify(req.Name, "channel"))
	body := channelRoom(req, localAlias, space.RoomID, s.matrix.ServerName(), s.bot())
	roomID, err := s.createRoom(ctx, creator, body)
	if err != nil {
		return Channel{}, err
	}
	child := matrix.StateEvent{Type: "m.space.child", StateKey: roomID, Content: map[string]any{
		"via": []string{s.matrix.ServerName()}, "suggested": req.Name == "general",
	}}
	if err := s.matrix.SendState(ctx, s.bot(), space.RoomID, child); err != nil {
		return Channel{}, errors.Internal("failed to link the channel into the space", err)
	}
	row := schemas.Channel{RoomID: roomID, SpaceID: space.ID, Name: req.Name, Alias: "#" + localAlias + ":" + s.matrix.ServerName()}
	if err := s.db.WithContext(ctx).Create(&row).Error; err != nil {
		return Channel{}, errors.Internal("failed to store the channel", err)
	}
	return Channel{RoomID: row.RoomID, Name: row.Name, Alias: row.Alias}, nil
}

func (s *Service) channels(ctx context.Context, spaceID int64) ([]Channel, error) {
	var rows []schemas.Channel
	if err := s.db.WithContext(ctx).Where("space_id = ?", spaceID).Order("name").Find(&rows).Error; err != nil {
		return nil, errors.Internal("failed to list channels", err)
	}
	out := make([]Channel, 0, len(rows))
	for _, row := range rows {
		out = append(out, Channel{RoomID: row.RoomID, Name: row.Name, Alias: row.Alias})
	}
	return out, nil
}

func (s *Service) freeAlias(ctx context.Context, base string) string {
	for n := 1; ; n++ {
		local := accounts.Candidate(base, n)
		var count int64
		s.db.WithContext(ctx).Model(&schemas.Channel{}).Where("alias = ?", "#"+local+":"+s.matrix.ServerName()).Count(&count)
		if count == 0 {
			return local
		}
	}
}

// admins lists the Prisme bot as a room admin. Room version 12 gives the
// creator unlimited power and refuses them in the users map.
func admins(bot string) map[string]any {
	return map[string]any{bot: 100}
}

func spaceRoom(name, bot string) matrix.CreateRoomRequest {
	return matrix.CreateRoomRequest{
		Name:            name,
		Preset:          "private_chat",
		Invite:          []string{bot},
		CreationContent: map[string]any{"type": "m.space"},
		PowerLevels:     map[string]any{"users": admins(bot), "invite": 0, "events_default": 100},
		InitialState: []matrix.StateEvent{
			{Type: "m.room.history_visibility", Content: map[string]string{"history_visibility": "shared"}},
		},
	}
}

func channelRoom(req CreateChannelRequest, alias, spaceID, server, bot string) matrix.CreateRoomRequest {
	return matrix.CreateRoomRequest{
		Name:        req.Name,
		Topic:       req.Topic,
		Alias:       alias,
		Preset:      "private_chat",
		Invite:      []string{bot},
		PowerLevels: map[string]any{"users": admins(bot), "invite": 0},
		InitialState: []matrix.StateEvent{
			{Type: "m.room.join_rules", Content: restrictedTo(spaceID)},
			{Type: "m.room.history_visibility", Content: map[string]string{"history_visibility": "shared"}},
			{Type: "m.space.parent", StateKey: spaceID, Content: map[string]any{"via": []string{server}, "canonical": true}},
		},
	}
}

func restrictedTo(spaceID string) map[string]any {
	allow := []map[string]string{{"type": "m.room_membership", "room_id": spaceID}}
	return map[string]any{"join_rule": "restricted", "allow": allow}
}
