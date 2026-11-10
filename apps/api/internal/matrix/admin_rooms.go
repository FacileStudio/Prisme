package matrix

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
)

// RoomEvent is one current state event as the admin API returns it.
type RoomEvent struct {
	Type     string          `json:"type"`
	StateKey string          `json:"state_key"`
	Sender   string          `json:"sender"`
	Content  json.RawMessage `json:"content"`
}

// RoomMembers lists the users joined to a room, as the server admin sees it.
func (c *Client) RoomMembers(ctx context.Context, roomID string) ([]string, error) {
	var out struct {
		Members []string `json:"members"`
	}
	err := c.do(ctx, request{method: http.MethodGet, path: adminRoom(roomID) + "/members"}, &out)
	return out.Members, err
}

// RoomState returns a room's current state events.
func (c *Client) RoomState(ctx context.Context, roomID string) ([]RoomEvent, error) {
	var out struct {
		State []RoomEvent `json:"state"`
	}
	err := c.do(ctx, request{method: http.MethodGet, path: adminRoom(roomID) + "/state"}, &out)
	return out.State, err
}

// MakeRoomAdmin raises userID to the power of the room's most powerful joined
// local member, sent as that member. It is how a room keeps an admin when its
// only one is about to be erased.
func (c *Client) MakeRoomAdmin(ctx context.Context, roomID, userID string) error {
	body := map[string]string{"user_id": userID}
	return c.do(ctx, request{method: http.MethodPost, path: adminRoom(roomID) + "/make_room_admin", body: body}, nil)
}

// PurgeRoom removes every local member from a room and deletes the room and
// its history from the database. Purging a room that is already gone succeeds.
func (c *Client) PurgeRoom(ctx context.Context, roomID string) error {
	body := map[string]bool{"purge": true, "block": false}
	return c.do(ctx, request{method: http.MethodDelete, path: adminRoom(roomID), body: body}, nil)
}

func adminRoom(roomID string) string {
	return adminPrefix + "/v1/rooms/" + url.PathEscape(roomID)
}
