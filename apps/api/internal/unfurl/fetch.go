package unfurl

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

func (c *Client) checkRedirect(req *http.Request, via []*http.Request) error {
	if len(via) >= maxRedirects {
		return fmt.Errorf("unfurl: more than %d redirects", maxRedirects)
	}
	if _, _, err := c.Allowed(req.URL.String()); err != nil {
		return err
	}
	return nil
}

func (c *Client) get(ctx context.Context, target, accept string, limit int64) ([]byte, string, bool) {
	ctx, cancel := context.WithTimeout(ctx, fetchTimeout)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, target, nil)
	if err != nil {
		return nil, "", false
	}
	req.Header.Set("Accept", accept)
	req.Header.Set("User-Agent", "Prisme link preview")
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, "", false
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return nil, "", false
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, limit))
	if err != nil {
		return nil, "", false
	}
	return body, resp.Header.Get("Content-Type"), true
}
