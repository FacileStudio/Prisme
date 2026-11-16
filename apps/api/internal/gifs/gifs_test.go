package gifs

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMediaURL(t *testing.T) {
	cases := map[string]bool{
		"https://media.giphy.com/media/x/giphy.gif":    true,
		"https://media3.giphy.com/media/x/200w.webp":   true,
		"https://i.giphy.com/x.gif":                    true,
		"http://media.giphy.com/media/x/giphy.gif":     false,
		"https://media.giphy.com.evil.test/x.gif":      false,
		"https://evil.media.giphy.com/x.gif":           false,
		"https://user@media.giphy.com/x.gif":           false,
		"https://media.giphy.com:8443/x.gif":           false,
		"https://127.0.0.1/x.gif":                      false,
		"https://api.giphy.com/v1/gifs/search?api_key": false,
	}
	for raw, want := range cases {
		if got := MediaURL(raw); got != want {
			t.Errorf("MediaURL(%q) = %v, want %v", raw, got, want)
		}
	}
}

func TestFetchRefusesOtherHosts(t *testing.T) {
	_, err := New("k").Fetch(context.Background(), "https://example.com/x.gif")
	if !errors.Is(err, ErrHost) {
		t.Fatalf("Fetch = %v, want ErrHost", err)
	}
}

func TestSearchParsesAndDropsForeignMedia(t *testing.T) {
	var query string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query = r.URL.RawQuery
		_, _ = w.Write([]byte(`{"data":[
			{"id":"a","title":"cat","images":{
				"fixed_width_downsampled":{"url":"https://media1.giphy.com/a/200w_d.gif","width":"200","height":"150"},
				"fixed_width":{"url":"https://media1.giphy.com/a/200w.gif","width":"200","height":"150"},
				"downsized_medium":{"url":"https://media1.giphy.com/a/giphy.gif","width":"480","height":"360"}}},
			{"id":"b","title":"bad","images":{
				"fixed_width_downsampled":{"url":"https://evil.test/b.gif"},
				"fixed_width":{"url":"https://evil.test/b.gif"},
				"downsized_medium":{"url":"https://evil.test/b.gif"}}}
		],"pagination":{"offset":0,"count":2}}`))
	}))
	defer server.Close()
	page, err := New("secret").WithBase(server.URL).Search(context.Background(), "cat", "fr", 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(page.GIFs) != 1 || page.GIFs[0].ID != "a" || page.GIFs[0].Full.Width != 480 || page.Next != 2 {
		t.Fatalf("page = %+v", page)
	}
	if query == "" || !contains(query, "rating=g") || !contains(query, "q=cat") || !contains(query, "lang=fr") {
		t.Fatalf("query = %q", query)
	}
}

func contains(s, sub string) bool { return strings.Contains(s, sub) }
