package matrix

import (
	"context"
	"net/http"
	"net/url"
)

const (
	adminPrefix    = "/_synapse/admin"
	mediaPageLimit = "1000"
	maxMediaPages  = 100
)

// AdminReady proves the appservice token passes Synapse's admin check. It is
// read-only and answers nothing about any person. The token is an admin
// credential only while the appservice sender is a server admin, which the
// Synapse image's entrypoint arranges (deploy/synapse/render.py).
func (c *Client) AdminReady(ctx context.Context) error {
	return c.do(ctx, request{method: http.MethodGet, path: adminPrefix + "/v1/background_updates/enabled"}, nil)
}

// Retired reports whether the homeserver holds a deactivated account under
// this id. Matrix never hands a deactivated id to anyone else, so such an id
// is taken for good even when Prisme holds no row for it.
func (c *Client) Retired(ctx context.Context, userID string) (bool, error) {
	var out struct {
		Deactivated bool `json:"deactivated"`
	}
	err := c.do(ctx, request{method: http.MethodGet, path: adminPrefix + "/v2/users/" + url.PathEscape(userID)}, &out)
	if IsCode(err, "M_NOT_FOUND") {
		return false, nil
	}
	return out.Deactivated, err
}

// EraseUser deactivates an account and asks Synapse to erase it: the profile
// goes, the user leaves every room, their membership events are redacted and
// their messages are pruned for anyone who was not in the room when they were
// sent. An id the homeserver never registered is already gone.
func (c *Client) EraseUser(ctx context.Context, userID string) error {
	path := adminPrefix + "/v1/deactivate/" + url.PathEscape(userID)
	err := c.do(ctx, request{method: http.MethodPost, path: path, body: map[string]bool{"erase": true}}, nil)
	if IsCode(err, "M_NOT_FOUND") {
		return nil
	}
	return err
}

// DeleteUserMedia deletes every file the user uploaded to the media
// repository, page by page, until none is left.
func (c *Client) DeleteUserMedia(ctx context.Context, userID string) error {
	path := adminPrefix + "/v1/users/" + url.PathEscape(userID) + "/media"
	for range maxMediaPages {
		var out struct {
			Deleted []string `json:"deleted_media"`
		}
		query := url.Values{"limit": {mediaPageLimit}}
		err := c.do(ctx, request{method: http.MethodDelete, path: path, query: query, body: map[string]any{}}, &out)
		if IsCode(err, "M_NOT_FOUND") || (err == nil && len(out.Deleted) == 0) {
			return nil
		}
		if err != nil {
			return err
		}
	}
	return &Error{Status: http.StatusInternalServerError, Code: "M_TOO_MUCH_MEDIA", Message: "media left after the last page"}
}

// UserRooms lists the rooms a user is joined to, as the server admin sees it,
// which also works for an account that can no longer sign in.
func (c *Client) UserRooms(ctx context.Context, userID string) ([]string, error) {
	var out struct {
		Rooms []string `json:"joined_rooms"`
	}
	path := adminPrefix + "/v1/users/" + url.PathEscape(userID) + "/joined_rooms"
	err := c.do(ctx, request{method: http.MethodGet, path: path}, &out)
	return out.Rooms, err
}
