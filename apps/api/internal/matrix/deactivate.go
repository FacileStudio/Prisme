package matrix

import (
	"context"
	"net/http"
)

// Deactivate closes a user's account for good. Synapse skips interactive auth
// for its appservice, so the call needs nothing but the user id. A GDPR
// erasure goes through EraseUser instead.
func (c *Client) Deactivate(ctx context.Context, userID string) error {
	body := map[string]any{"erase": false}
	path := "/_matrix/client/v3/account/deactivate"
	return c.do(ctx, request{method: http.MethodPost, path: path, asUser: userID, body: body}, nil)
}
