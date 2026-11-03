package matrix

import (
	"context"
	"net/http"
	"net/url"
)

// Login is an access token minted for a user, as handed to a Matrix client.
type Login struct {
	UserID      string `json:"user_id"`
	AccessToken string `json:"access_token"`
	DeviceID    string `json:"device_id"`
}

// UserID builds the full id of a local user.
func (c *Client) UserID(localpart string) string {
	return "@" + localpart + ":" + c.cfg.ServerName
}

// Register creates a passwordless user in the appservice namespace. It sets
// inhibit_login, as MSC4190 asks, so registering leaves no stray device; the
// tokens come from Login. An existing user is not an error.
func (c *Client) Register(ctx context.Context, localpart string) error {
	body := map[string]any{"type": "m.login.application_service", "username": localpart, "inhibit_login": true}
	err := c.do(ctx, request{method: http.MethodPost, path: "/_matrix/client/v3/register", body: body}, nil)
	if IsCode(err, "M_USER_IN_USE") {
		return nil
	}
	return err
}

// Login mints a new device and access token for a user of the namespace.
func (c *Client) Login(ctx context.Context, userID, deviceName string) (Login, error) {
	body := map[string]any{
		"type":                        "m.login.application_service",
		"identifier":                  map[string]string{"type": "m.id.user", "user": userID},
		"initial_device_display_name": deviceName,
	}
	var login Login
	err := c.do(ctx, request{method: http.MethodPost, path: "/_matrix/client/v3/login", body: body}, &login)
	return login, err
}

// Logout revokes one access token and deletes its device.
func (c *Client) Logout(ctx context.Context, accessToken string) error {
	err := c.do(ctx, request{method: http.MethodPost, path: "/_matrix/client/v3/logout", token: accessToken}, nil)
	if IsCode(err, "M_UNKNOWN_TOKEN") {
		return nil
	}
	return err
}

// WhoAmI reports whether an access token is still valid.
func (c *Client) WhoAmI(ctx context.Context, accessToken string) bool {
	return c.do(ctx, request{method: http.MethodGet, path: "/_matrix/client/v3/account/whoami", token: accessToken}, nil) == nil
}

// Deactivate closes a user's account for good. Synapse skips interactive auth
// for its appservice, so the call needs nothing but the user id.
func (c *Client) Deactivate(ctx context.Context, userID string) error {
	body := map[string]any{"erase": false}
	path := "/_matrix/client/v3/account/deactivate"
	return c.do(ctx, request{method: http.MethodPost, path: path, asUser: userID, body: body}, nil)
}

// SetProfile sets a user's display name and, when given, avatar (an mxc URI).
func (c *Client) SetProfile(ctx context.Context, userID, displayName, avatarMXC string) error {
	base := "/_matrix/client/v3/profile/" + url.PathEscape(userID)
	if displayName != "" {
		body := map[string]string{"displayname": displayName}
		if err := c.do(ctx, request{method: http.MethodPut, path: base + "/displayname", asUser: userID, body: body}, nil); err != nil {
			return err
		}
	}
	if avatarMXC == "" {
		return nil
	}
	body := map[string]string{"avatar_url": avatarMXC}
	return c.do(ctx, request{method: http.MethodPut, path: base + "/avatar_url", asUser: userID, body: body}, nil)
}

// DirectRooms reads a user's m.direct account data: each partner's user id
// to the rooms of their direct messages. A user who never had one answers
// an empty map.
func (c *Client) DirectRooms(ctx context.Context, userID string) (map[string][]string, error) {
	out := map[string][]string{}
	path := "/_matrix/client/v3/user/" + url.PathEscape(userID) + "/account_data/m.direct"
	err := c.do(ctx, request{method: http.MethodGet, path: path, asUser: userID}, &out)
	if IsCode(err, "M_NOT_FOUND") {
		return map[string][]string{}, nil
	}
	return out, err
}
