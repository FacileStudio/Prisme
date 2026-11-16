// Package nuage reads a user's Nuage drive and shares one of its files, acting
// as that user with a token Registre exchanged for Nuage's audience. Nuage
// checks every call against the user's own access, so Prisme never sees a
// file the user could not open in Nuage itself.
package nuage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// ErrNotFound is Nuage's answer for a file or folder the user cannot reach.
var ErrNotFound = errors.New("nuage: not found")

// ErrRefused is Nuage refusing the delegated token itself.
var ErrRefused = errors.New("nuage: refused")

const maxBody = 4 << 20

// Tokens mints a Nuage token for a user and drops it after Nuage refused it.
type Tokens struct {
	Token  func(ctx context.Context, userID int64) (string, error)
	Forget func(userID int64)
}

// Client talks to one Nuage deployment on behalf of its users.
type Client struct {
	baseURL string
	tokens  Tokens
	http    *http.Client
}

// New builds a client for the Nuage deployment at baseURL.
func New(baseURL string, tokens Tokens) *Client {
	return &Client{baseURL: strings.TrimRight(baseURL, "/"), tokens: tokens, http: &http.Client{Timeout: 15 * time.Second}}
}

// ShareURL is the public page of a share token.
func (c *Client) ShareURL(token string) string {
	return c.baseURL + "/s/" + url.PathEscape(token)
}

type call struct {
	method string
	path   string
	body   any
	out    any
}

func (c *Client) get(ctx context.Context, userID int64, path string, out any) error {
	return c.send(ctx, userID, call{method: http.MethodGet, path: path, out: out})
}

func (c *Client) send(ctx context.Context, userID int64, req call) error {
	status, err := c.attempt(ctx, userID, req)
	if status == http.StatusUnauthorized {
		c.tokens.Forget(userID)
		status, err = c.attempt(ctx, userID, req)
	}
	if status == http.StatusUnauthorized || status == http.StatusForbidden {
		return ErrRefused
	}
	return err
}

func (c *Client) attempt(ctx context.Context, userID int64, req call) (int, error) {
	token, err := c.tokens.Token(ctx, userID)
	if err != nil {
		return 0, err
	}
	httpReq, err := newRequest(ctx, req.method, c.baseURL+req.path, req.body)
	if err != nil {
		return 0, err
	}
	httpReq.Header.Set("Authorization", "Bearer "+token)
	resp, err := c.http.Do(httpReq)
	if err != nil {
		return 0, fmt.Errorf("nuage: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()
	return resp.StatusCode, decode(resp, req.out)
}

func newRequest(ctx context.Context, method, target string, body any) (*http.Request, error) {
	if body == nil {
		return http.NewRequestWithContext(ctx, method, target, nil)
	}
	raw, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, method, target, strings.NewReader(string(raw)))
	if err == nil {
		req.Header.Set("Content-Type", "application/json")
	}
	return req, err
}

func decode(resp *http.Response, out any) error {
	switch {
	case resp.StatusCode == http.StatusNotFound:
		return ErrNotFound
	case resp.StatusCode < 200 || resp.StatusCode > 299:
		return fmt.Errorf("nuage: status %d", resp.StatusCode)
	}
	if err := json.NewDecoder(io.LimitReader(resp.Body, maxBody)).Decode(out); err != nil {
		return fmt.Errorf("nuage: unreadable answer: %w", err)
	}
	return nil
}
