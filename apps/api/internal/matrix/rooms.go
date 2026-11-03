package matrix

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"net/url"
)

// CreateRoomRequest is the body of POST /createRoom, trimmed to what Prisme
// sends.
type CreateRoomRequest struct {
	Name            string         `json:"name,omitempty"`
	Topic           string         `json:"topic,omitempty"`
	Alias           string         `json:"room_alias_name,omitempty"`
	Preset          string         `json:"preset,omitempty"`
	Visibility      string         `json:"visibility,omitempty"`
	Invite          []string       `json:"invite,omitempty"`
	IsDirect        bool           `json:"is_direct,omitempty"`
	CreationContent map[string]any `json:"creation_content,omitempty"`
	InitialState    []StateEvent   `json:"initial_state,omitempty"`
	PowerLevels     map[string]any `json:"power_level_content_override,omitempty"`
}

// StateEvent is one state event for initial_state.
type StateEvent struct {
	Type     string `json:"type"`
	StateKey string `json:"state_key"`
	Content  any    `json:"content"`
}

// CreateRoom creates a room as the given user and returns its id.
func (c *Client) CreateRoom(ctx context.Context, asUser string, body CreateRoomRequest) (string, error) {
	var out struct {
		RoomID string `json:"room_id"`
	}
	err := c.do(ctx, request{method: http.MethodPost, path: "/_matrix/client/v3/createRoom", asUser: asUser, body: body}, &out)
	return out.RoomID, err
}

// SendState sets a state event in a room, as asUser.
func (c *Client) SendState(ctx context.Context, asUser, roomID string, event StateEvent) error {
	path := "/_matrix/client/v3/rooms/" + url.PathEscape(roomID) + "/state/" +
		url.PathEscape(event.Type) + "/" + url.PathEscape(event.StateKey)
	return c.do(ctx, request{method: http.MethodPut, path: path, asUser: asUser, body: event.Content}, nil)
}

// SendMessage posts an m.room.message as asUser and returns the event id.
func (c *Client) SendMessage(ctx context.Context, asUser, roomID string, content any) (string, error) {
	path := "/_matrix/client/v3/rooms/" + url.PathEscape(roomID) + "/send/m.room.message/" + txnID()
	var out struct {
		EventID string `json:"event_id"`
	}
	err := c.do(ctx, request{method: http.MethodPut, path: path, asUser: asUser, body: content}, &out)
	return out.EventID, err
}

// RoomName reads a room's m.room.name, as a member. A room without one
// answers an empty name.
func (c *Client) RoomName(ctx context.Context, asUser, roomID string) string {
	var out struct {
		Name string `json:"name"`
	}
	path := "/_matrix/client/v3/rooms/" + url.PathEscape(roomID) + "/state/m.room.name/"
	if err := c.do(ctx, request{method: http.MethodGet, path: path, asUser: asUser}, &out); err != nil {
		return ""
	}
	return out.Name
}

func txnID() string {
	random := make([]byte, 12)
	_, _ = rand.Read(random)
	return "prisme." + hex.EncodeToString(random)
}
