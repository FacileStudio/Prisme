// Package matrix is Prisme's application-service client for Synapse.
//
// Every call authenticates with the appservice token. A call made for a user
// passes that user's id and Synapse runs it as them (identity assertion), so
// Prisme never stores a password and only keeps the access tokens it hands to
// browsers and agents.
package matrix

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// Config names the homeserver and the appservice credential.
type Config struct {
	BaseURL    string
	ServerName string
	ASToken    string
}

// Client talks to Synapse's client-server API as the application service.
type Client struct {
	cfg  Config
	http *http.Client
}

// Error is a Matrix error response: the HTTP status and the errcode.
type Error struct {
	Status  int
	Code    string `json:"errcode"`
	Message string `json:"error"`
}

// Error formats the status, errcode and message.
func (e *Error) Error() string {
	return fmt.Sprintf("matrix: %d %s: %s", e.Status, e.Code, e.Message)
}

// IsCode reports whether err is a Matrix error carrying the given errcode.
func IsCode(err error, code string) bool {
	matrixErr, ok := err.(*Error)
	return ok && matrixErr.Code == code
}

// New builds a client over the given homeserver.
func New(cfg Config) *Client {
	return &Client{cfg: cfg, http: &http.Client{Timeout: 30 * time.Second}}
}

// ServerName is the homeserver's domain, the part after the colon in an id.
func (c *Client) ServerName() string { return c.cfg.ServerName }

// request is one call: who it runs as and what it sends.
type request struct {
	method string
	path   string
	asUser string
	token  string
	body   any
	query  url.Values
}

func (c *Client) do(ctx context.Context, req request, out any) error {
	httpReq, err := c.build(ctx, req)
	if err != nil {
		return err
	}
	resp, err := c.http.Do(httpReq)
	if err != nil {
		return fmt.Errorf("matrix: %s %s: %w", req.method, req.path, err)
	}
	defer func() { _ = resp.Body.Close() }()
	data, err := io.ReadAll(io.LimitReader(resp.Body, 4<<20))
	if err != nil {
		return err
	}
	if resp.StatusCode >= 300 {
		matrixErr := &Error{Status: resp.StatusCode}
		_ = json.Unmarshal(data, matrixErr)
		return matrixErr
	}
	if out == nil || len(data) == 0 {
		return nil
	}
	return json.Unmarshal(data, out)
}

func (c *Client) build(ctx context.Context, req request) (*http.Request, error) {
	query := req.query
	if query == nil {
		query = url.Values{}
	}
	if req.asUser != "" {
		query.Set("user_id", req.asUser)
	}
	target := c.cfg.BaseURL + req.path
	if len(query) > 0 {
		target += "?" + query.Encode()
	}
	var body io.Reader
	if req.body != nil {
		encoded, err := json.Marshal(req.body)
		if err != nil {
			return nil, err
		}
		body = bytes.NewReader(encoded)
	}
	httpReq, err := http.NewRequestWithContext(ctx, req.method, target, body)
	if err != nil {
		return nil, err
	}
	token := req.token
	if token == "" {
		token = c.cfg.ASToken
	}
	httpReq.Header.Set("Authorization", "Bearer "+token)
	httpReq.Header.Set("Content-Type", "application/json")
	return httpReq, nil
}
