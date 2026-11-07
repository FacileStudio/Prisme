package matrix

import (
	"context"
	"net/http"
)

// Deactivate closes a user's account for good. Synapse skips interactive auth
// for its appservice, so the call needs nothing but the user id.
func (c *Client) Deactivate(ctx context.Context, userID string) error {
	return c.deactivate(ctx, userID, false)
}

// Erase closes a user's account for good and asks Synapse to erase it: the
// profile is dropped and the user's messages are hidden from anyone who has
// not already seen them. It is how a GDPR deletion reaches Matrix.
func (c *Client) Erase(ctx context.Context, userID string) error {
	return c.deactivate(ctx, userID, true)
}

func (c *Client) deactivate(ctx context.Context, userID string, erase bool) error {
	body := map[string]any{"erase": erase}
	path := "/_matrix/client/v3/account/deactivate"
	return c.do(ctx, request{method: http.MethodPost, path: path, asUser: userID, body: body}, nil)
}
