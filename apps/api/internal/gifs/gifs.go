// Package gifs searches GIPHY from the server, so a Prisme user's address never
// reaches GIPHY and the client's CSP stays on its own origin. Media is fetched
// through Fetch, which only follows https links to GIPHY's media hosts.
package gifs

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// MaxMedia caps the bytes Fetch reads for one GIF.
const MaxMedia = 8 << 20

// PageSize is how many GIFs one search or trending call returns.
const PageSize = 24

// ErrHost is returned for a media URL that is not on GIPHY's media hosts.
var ErrHost = errors.New("gifs: not a GIPHY media url")

// Image is one rendition of a GIF.
type Image struct {
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

// GIF is one search result: a small preview for the grid and the rendition
// that is sent.
type GIF struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Preview Image  `json:"preview"`
	Full    Image  `json:"full"`
}

// Page is one page of results and the offset of the next one.
type Page struct {
	GIFs []GIF `json:"gifs"`
	Next int   `json:"next"`
}

// Media is the body of a fetched GIF.
type Media struct {
	Body        []byte
	ContentType string
}

var mediaHost = regexp.MustCompile(`^(i|media|media[0-9])\.giphy\.com$`)

// Client calls GIPHY's REST API with one key.
type Client struct {
	key   string
	base  string
	api   *http.Client
	media *http.Client
}

// New builds a client for GIPHY's public API.
func New(key string) *Client {
	return &Client{
		key:   key,
		base:  "https://api.giphy.com",
		api:   &http.Client{Timeout: 10 * time.Second},
		media: &http.Client{Timeout: 15 * time.Second, CheckRedirect: sameHosts},
	}
}

// WithBase points the client at another API origin, for tests.
func (c *Client) WithBase(base string) *Client {
	c.base = strings.TrimRight(base, "/")
	return c
}

// Search finds GIFs for a query, rated G, in the given language.
func (c *Client) Search(ctx context.Context, query, lang string, offset int) (Page, error) {
	params := url.Values{"q": {query}, "lang": {lang}}
	return c.page(ctx, "/v1/gifs/search", params, offset)
}

// Trending lists the GIFs GIPHY trends today.
func (c *Client) Trending(ctx context.Context, offset int) (Page, error) {
	return c.page(ctx, "/v1/gifs/trending", url.Values{}, offset)
}

// Fetch downloads one GIPHY media file.
func (c *Client) Fetch(ctx context.Context, raw string) (Media, error) {
	if !MediaURL(raw) {
		return Media{}, ErrHost
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, raw, nil)
	if err != nil {
		return Media{}, err
	}
	resp, err := c.media.Do(req)
	if err != nil {
		return Media{}, err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return Media{}, fmt.Errorf("gifs: media answered %d", resp.StatusCode)
	}
	kind := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(kind, "image/") && !strings.HasPrefix(kind, "video/") {
		return Media{}, fmt.Errorf("gifs: media is %q", kind)
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, MaxMedia+1))
	if err != nil {
		return Media{}, err
	}
	if len(body) > MaxMedia {
		return Media{}, errors.New("gifs: media is too large")
	}
	return Media{Body: body, ContentType: kind}, nil
}

// MediaURL reports whether raw is an https link on GIPHY's media hosts.
func MediaURL(raw string) bool {
	parsed, err := url.Parse(raw)
	if err != nil || parsed.Scheme != "https" || parsed.User != nil || parsed.Port() != "" {
		return false
	}
	return mediaHost.MatchString(strings.ToLower(parsed.Hostname()))
}

func sameHosts(req *http.Request, _ []*http.Request) error {
	if !MediaURL(req.URL.String()) {
		return ErrHost
	}
	return nil
}

func (c *Client) page(ctx context.Context, path string, params url.Values, offset int) (Page, error) {
	params.Set("api_key", c.key)
	params.Set("limit", strconv.Itoa(PageSize))
	params.Set("offset", strconv.Itoa(max(offset, 0)))
	params.Set("rating", "g")
	params.Set("bundle", "messaging_non_clips")
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.base+path+"?"+params.Encode(), nil)
	if err != nil {
		return Page{}, err
	}
	resp, err := c.api.Do(req)
	if err != nil {
		return Page{}, err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return Page{}, fmt.Errorf("gifs: GIPHY answered %d", resp.StatusCode)
	}
	var list listing
	if err := json.NewDecoder(io.LimitReader(resp.Body, 4<<20)).Decode(&list); err != nil {
		return Page{}, err
	}
	return toPage(list), nil
}
