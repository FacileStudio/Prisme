// Package unfurl reads the public preview of a link to another Facile app:
// its oEmbed answer when the app serves one, else the Open Graph tags of its
// page. Only https links to the allowlisted suite hosts are fetched, without
// cookies, and a redirect off the allowlist ends the fetch.
package unfurl

import (
	"context"
	"errors"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// ErrNotAllowed means the link is not an https link to an allowlisted app.
var ErrNotAllowed = errors.New("unfurl: link is not to an allowlisted app")

// SuiteHosts maps each production host Prisme previews to its app. Capsule is
// absent on purpose: fetching a capsule link must never be what burns it.
func SuiteHosts() map[string]string {
	return map[string]string{
		"feuilles.facile.studio": "feuilles",
		"docs.facile.studio":     "docs",
		"echo.facile.studio":     "echo",
		"clap.facile.studio":     "clap",
		"diapo.facile.studio":    "diapo",
		"nuage.facile.studio":    "nuage",
	}
}

const (
	fetchTimeout = 5 * time.Second
	maxPageBytes = 512 << 10
	maxEmbedJSON = 64 << 10
	maxRedirects = 3
)

// Preview is what a link card shows. Every URL in it points at an allowlisted
// host; App alone is set when the app published nothing about the link.
type Preview struct {
	URL         string `json:"url"`
	App         string `json:"app"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Image       string `json:"image,omitempty"`
	EmbedURL    string `json:"embed_url,omitempty"`
	DurationMS  int64  `json:"duration_ms,omitempty"`
}

// Client fetches previews for one set of allowlisted hosts.
type Client struct {
	hosts map[string]string
	http  *http.Client
	cache *cache
}

// New builds a client for hosts, which maps a host name to its app name.
func New(hosts map[string]string) *Client {
	dialer := &net.Dialer{Timeout: fetchTimeout, Control: refusePrivate}
	transport := &http.Transport{
		DialContext:           dialer.DialContext,
		TLSHandshakeTimeout:   fetchTimeout,
		ResponseHeaderTimeout: fetchTimeout,
		MaxIdleConns:          16,
		IdleConnTimeout:       time.Minute,
	}
	return newClient(hosts, &http.Client{Transport: transport})
}

func newClient(hosts map[string]string, client *http.Client) *Client {
	c := &Client{hosts: hosts, cache: newCache(512, 10*time.Minute)}
	client.Timeout = fetchTimeout
	client.CheckRedirect = c.checkRedirect
	c.http = client
	return c
}

// Allowed resolves raw to its canonical form, without the fragment that may
// carry a key, and the app it belongs to.
func (c *Client) Allowed(raw string) (*url.URL, string, error) {
	u, err := url.Parse(strings.TrimSpace(raw))
	if err != nil || u.Scheme != "https" || u.User != nil || u.Port() != "" || u.Opaque != "" {
		return nil, "", ErrNotAllowed
	}
	app, ok := c.hosts[strings.ToLower(u.Hostname())]
	if !ok {
		return nil, "", ErrNotAllowed
	}
	u.Host = strings.ToLower(u.Hostname())
	u.Fragment, u.RawFragment = "", ""
	return u, app, nil
}

// Preview returns the card for raw. A link the app says nothing about, or an
// app that does not answer, still gets a card naming the app.
func (c *Client) Preview(ctx context.Context, raw string) (Preview, error) {
	u, app, err := c.Allowed(raw)
	if err != nil {
		return Preview{}, err
	}
	key := u.String()
	if cached, ok := c.cache.get(key); ok {
		return cached, nil
	}
	preview := Preview{URL: key, App: app}
	if embed, ok := c.oembed(ctx, u); ok {
		preview = c.merge(preview, embed)
	} else if page, ok := c.specificPage(ctx, u); ok {
		preview = c.merge(preview, page)
	}
	c.cache.put(key, preview)
	return preview, nil
}

func (c *Client) merge(base, found Preview) Preview {
	base.Title = clip(found.Title, 200)
	base.Description = clip(found.Description, 300)
	base.Image = c.sameSuite(found.Image)
	base.EmbedURL = c.embedURL(found.EmbedURL)
	base.DurationMS = found.DurationMS
	if strings.EqualFold(base.Title, base.App) {
		base.Title = ""
	}
	return base
}

func (c *Client) sameSuite(raw string) string {
	if raw == "" {
		return ""
	}
	u, _, err := c.Allowed(raw)
	if err != nil {
		return ""
	}
	return u.String()
}

func (c *Client) embedURL(raw string) string {
	u, _, err := c.Allowed(raw)
	if err != nil || !strings.HasPrefix(u.Path, "/embed/") {
		return ""
	}
	return u.String()
}
