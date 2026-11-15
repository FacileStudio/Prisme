// Package capsule relays sealed messages to Capsule's public paste API. The
// browser encrypts before anything reaches Prisme, so this package only ever
// carries ciphertext and never sees a key.
package capsule

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"time"
)

// MaxCiphertext is Capsule's default paste size limit.
const MaxCiphertext = 16 << 20

// ErrNotFound is Capsule's answer for a capsule that was opened, expired or
// never existed.
var ErrNotFound = errors.New("capsule: not found")

var idPattern = regexp.MustCompile(`^cap_[0-9a-f]{16}$`)

var expiries = map[string]bool{"1h": true, "24h": true, "7d": true, "30d": true}

// ValidID reports whether id has the shape Capsule gives its capsules.
func ValidID(id string) bool { return idPattern.MatchString(id) }

// ValidExpiry reports whether Capsule accepts expiresIn.
func ValidExpiry(expiresIn string) bool { return expiries[expiresIn] }

// Client talks to one Capsule deployment.
type Client struct {
	baseURL string
	http    *http.Client
}

// Sealed is a freshly created capsule.
type Sealed struct {
	ID        string     `json:"id"`
	ExpiresAt *time.Time `json:"expires_at"`
}

// Meta is what Capsule says about a capsule before anyone opens it.
type Meta struct {
	ID        string     `json:"id"`
	Exists    bool       `json:"exists"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

type createRequest struct {
	Content       string `json:"content"`
	BurnAfterRead bool   `json:"burn_after_read"`
	ExpiresIn     string `json:"expires_in"`
}

// New builds a client for the Capsule deployment at baseURL.
func New(baseURL string) *Client {
	return &Client{baseURL: baseURL, http: &http.Client{Timeout: 15 * time.Second}}
}

// BaseURL is the public origin of the Capsule deployment.
func (c *Client) BaseURL() string { return c.baseURL }

// Seal stores ciphertext as a capsule that burns on its first read.
func (c *Client) Seal(ctx context.Context, ciphertext, expiresIn string) (Sealed, error) {
	var sealed Sealed
	body := createRequest{Content: ciphertext, BurnAfterRead: true, ExpiresIn: expiresIn}
	err := c.send(ctx, http.MethodPost, "/api/pastes", body, &sealed)
	return sealed, err
}

// Meta reads whether a capsule can still be opened.
func (c *Client) Meta(ctx context.Context, id string) (Meta, error) {
	var meta Meta
	err := c.send(ctx, http.MethodGet, "/api/pastes/"+url.PathEscape(id), nil, &meta)
	return meta, err
}

// Open returns a capsule's ciphertext, which burns it.
func (c *Client) Open(ctx context.Context, id string) (string, error) {
	var out struct {
		Content string `json:"content"`
	}
	err := c.send(ctx, http.MethodPost, "/api/pastes/"+url.PathEscape(id)+"/content", struct{}{}, &out)
	return out.Content, err
}

func (c *Client) send(ctx context.Context, method, path string, body, out any) error {
	var payload io.Reader
	if body != nil {
		encoded, err := json.Marshal(body)
		if err != nil {
			return err
		}
		payload = bytes.NewReader(encoded)
	}
	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, payload)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("capsule: %s %s: %w", method, "/api/pastes", err)
	}
	defer func() { _ = resp.Body.Close() }()
	data, _ := io.ReadAll(io.LimitReader(resp.Body, MaxCiphertext+(1<<20)))
	if resp.StatusCode == http.StatusNotFound {
		return ErrNotFound
	}
	if resp.StatusCode >= 300 {
		return fmt.Errorf("capsule: %s /api/pastes answered %d", method, resp.StatusCode)
	}
	return json.Unmarshal(data, out)
}
