package unfurl

import (
	"context"
	"encoding/json"
	"html"
	"net/url"
	"regexp"
	"strings"
)

var (
	metaTag   = regexp.MustCompile(`(?is)<meta\s[^>]*>`)
	attribute = regexp.MustCompile(`(?is)([a-z:-]+)\s*=\s*(?:"([^"]*)"|'([^']*)')`)
	titleTag  = regexp.MustCompile(`(?is)<title[^>]*>(.*?)</title>`)
	iframeSrc = regexp.MustCompile(`(?is)<iframe\s[^>]*src\s*=\s*"([^"]+)"`)
	headEnd   = regexp.MustCompile(`(?i)</head>`)
)

type oembedAnswer struct {
	Title        string `json:"title"`
	ThumbnailURL string `json:"thumbnail_url"`
	HTML         string `json:"html"`
	DurationMS   int64  `json:"duration_ms"`
}

func (c *Client) oembed(ctx context.Context, link *url.URL) (Preview, bool) {
	endpoint := url.URL{Scheme: "https", Host: link.Host, Path: "/oembed"}
	endpoint.RawQuery = url.Values{"url": {link.String()}, "format": {"json"}}.Encode()
	body, kind, ok := c.get(ctx, endpoint.String(), "application/json", maxEmbedJSON)
	if !ok || !strings.Contains(kind, "json") {
		return Preview{}, false
	}
	var answer oembedAnswer
	if json.Unmarshal(body, &answer) != nil || answer.Title == "" {
		return Preview{}, false
	}
	preview := Preview{Title: answer.Title, Image: answer.ThumbnailURL, DurationMS: answer.DurationMS}
	if match := iframeSrc.FindStringSubmatch(answer.HTML); match != nil {
		preview.EmbedURL = html.UnescapeString(match[1])
	}
	return preview, true
}

func (c *Client) page(ctx context.Context, link *url.URL) (Preview, bool) {
	body, kind, ok := c.get(ctx, link.String(), "text/html", maxPageBytes)
	if !ok || !strings.Contains(kind, "html") {
		return Preview{}, false
	}
	head := string(body)
	if loc := headEnd.FindStringIndex(head); loc != nil {
		head = head[:loc[0]]
	}
	preview := openGraph(head)
	if preview.Title == "" {
		if match := titleTag.FindStringSubmatch(head); match != nil {
			preview.Title = html.UnescapeString(match[1])
		}
	}
	return preview, preview.Title != ""
}

func openGraph(head string) Preview {
	var preview Preview
	for _, tag := range metaTag.FindAllString(head, -1) {
		attrs := attributes(tag)
		name := strings.ToLower(attrs["property"])
		if name == "" {
			name = strings.ToLower(attrs["name"])
		}
		value := attrs["content"]
		switch name {
		case "og:title":
			preview.Title = value
		case "og:description", "description":
			if preview.Description == "" || name == "og:description" {
				preview.Description = value
			}
		case "og:image":
			preview.Image = value
		}
	}
	return preview
}

func attributes(tag string) map[string]string {
	attrs := map[string]string{}
	for _, match := range attribute.FindAllStringSubmatch(tag, -1) {
		value := match[2]
		if value == "" {
			value = match[3]
		}
		attrs[strings.ToLower(match[1])] = html.UnescapeString(value)
	}
	return attrs
}

func (c *Client) specificPage(ctx context.Context, link *url.URL) (Preview, bool) {
	found, ok := c.page(ctx, link)
	if !ok || link.Path == "" || link.Path == "/" {
		return found, ok
	}
	root := url.URL{Scheme: "https", Host: link.Host, Path: "/"}
	key := "root:" + link.Host
	landing, cached := c.cache.get(key)
	if !cached {
		landing, _ = c.page(ctx, &root)
		c.cache.put(key, landing)
	}
	if landing.Title != "" && landing.Title == found.Title {
		return Preview{}, false
	}
	return found, true
}
