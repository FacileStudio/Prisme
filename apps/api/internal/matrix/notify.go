package matrix

import (
	"context"
	"net/http"
	"net/url"
	"slices"
)

// DisplayName reads a user's display name, or answers the id when there is none.
func (c *Client) DisplayName(ctx context.Context, userID string) string {
	var out struct {
		DisplayName string `json:"displayname"`
	}
	path := "/_matrix/client/v3/profile/" + url.PathEscape(userID) + "/displayname"
	if err := c.do(ctx, request{method: http.MethodGet, path: path}, &out); err != nil || out.DisplayName == "" {
		return userID
	}
	return out.DisplayName
}

// RoomMuted reports whether a user turned a room's notifications off: a room
// push rule with no action or with dont_notify, which is what the client writes.
func (c *Client) RoomMuted(ctx context.Context, userID, roomID string) bool {
	var out struct {
		Actions []any `json:"actions"`
	}
	path := "/_matrix/client/v3/pushrules/global/room/" + url.PathEscape(roomID)
	if err := c.do(ctx, request{method: http.MethodGet, path: path, asUser: userID}, &out); err != nil {
		return false
	}
	return len(out.Actions) == 0 || slices.Contains(out.Actions, any("dont_notify"))
}
