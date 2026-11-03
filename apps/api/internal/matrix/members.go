package matrix

import (
	"context"
	"net/http"
	"net/url"
	"slices"
)

// Invite invites a user into a room, as asUser. Inviting someone already
// joined is not an error.
func (c *Client) Invite(ctx context.Context, asUser, roomID, userID string) error {
	path := "/_matrix/client/v3/rooms/" + url.PathEscape(roomID) + "/invite"
	err := c.do(ctx, request{method: http.MethodPost, path: path, asUser: asUser, body: map[string]string{"user_id": userID}}, nil)
	if IsCode(err, "M_FORBIDDEN") && c.IsJoined(ctx, userID, roomID) {
		return nil
	}
	return err
}

// Join joins a room as the given user.
func (c *Client) Join(ctx context.Context, asUser, roomID string) error {
	path := "/_matrix/client/v3/join/" + url.PathEscape(roomID)
	return c.do(ctx, request{method: http.MethodPost, path: path, asUser: asUser, body: map[string]any{}}, nil)
}

// Leave leaves a room as the given user.
func (c *Client) Leave(ctx context.Context, asUser, roomID string) error {
	path := "/_matrix/client/v3/rooms/" + url.PathEscape(roomID) + "/leave"
	return c.do(ctx, request{method: http.MethodPost, path: path, asUser: asUser, body: map[string]any{}}, nil)
}

// JoinedRooms lists the rooms a user is joined to.
func (c *Client) JoinedRooms(ctx context.Context, userID string) ([]string, error) {
	var out struct {
		Rooms []string `json:"joined_rooms"`
	}
	err := c.do(ctx, request{method: http.MethodGet, path: "/_matrix/client/v3/joined_rooms", asUser: userID}, &out)
	return out.Rooms, err
}

// IsJoined reports whether a user is joined to a room.
func (c *Client) IsJoined(ctx context.Context, userID, roomID string) bool {
	rooms, err := c.JoinedRooms(ctx, userID)
	if err != nil {
		return false
	}
	return slices.Contains(rooms, roomID)
}
