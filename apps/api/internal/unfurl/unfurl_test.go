package unfurl

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
)

func suiteServer(t *testing.T, handler http.HandlerFunc) *Client {
	t.Helper()
	server := httptest.NewTLSServer(handler)
	t.Cleanup(server.Close)
	transport := server.Client().Transport.(*http.Transport).Clone()
	transport.TLSClientConfig.ServerName = "example.com"
	transport.DialContext = func(ctx context.Context, network, _ string) (net.Conn, error) {
		return (&net.Dialer{}).DialContext(ctx, network, server.Listener.Addr().String())
	}
	return newClient(SuiteHosts(), &http.Client{Transport: transport})
}

func TestAllowedRefusesEverythingButSuiteHTTPS(t *testing.T) {
	client := New(SuiteHosts())
	refused := []string{
		"http://docs.facile.studio/d/1",
		"https://127.0.0.1/",
		"https://[::1]/",
		"https://user@docs.facile.studio/",
		"https://docs.facile.studio:8443/",
		"https://capsule.facile.studio/cap_1",
		"https://evil.example/",
		"https://docs.facile.studio.evil.example/",
		"javascript:alert(1)",
	}
	for _, raw := range refused {
		if _, _, err := client.Allowed(raw); err == nil {
			t.Errorf("%s was allowed", raw)
		}
	}
	u, app, err := client.Allowed("https://Echo.facile.studio/room/team#key=secret")
	if err != nil || app != "echo" || u.String() != "https://echo.facile.studio/room/team" {
		t.Fatalf("got %v %q %v, want the echo link without its fragment", u, app, err)
	}
}

func TestPreviewPrefersOEmbedAndKeepsOnlySuiteURLs(t *testing.T) {
	client := suiteServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/oembed" || r.URL.Query().Get("url") != "https://clap.facile.studio/v/tok" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"version":"1.0","type":"rich","title":"Launch cut","duration_ms":42000,` +
			`"thumbnail_url":"https://cdn.evil.example/p.jpg",` +
			`"html":"<iframe src=\"https://clap.facile.studio/embed/tok\"></iframe>"}`))
	})
	preview, err := client.Preview(context.Background(), "https://clap.facile.studio/v/tok")
	if err != nil {
		t.Fatal(err)
	}
	if preview.Title != "Launch cut" || preview.EmbedURL != "https://clap.facile.studio/embed/tok" || preview.DurationMS != 42000 {
		t.Fatalf("unexpected preview %+v", preview)
	}
	if preview.Image != "" {
		t.Fatalf("an off-suite thumbnail was kept: %q", preview.Image)
	}
}

func TestPreviewFallsBackToOpenGraph(t *testing.T) {
	client := suiteServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/oembed" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if r.URL.Path == "/" {
			_, _ = w.Write([]byte(`<head><meta property="og:title" content="Docs · Write together"></head>`))
			return
		}
		_, _ = w.Write([]byte(`<html><head><title>Docs</title>` +
			`<meta property="og:title" content="Q3 &amp; plan">` +
			`<meta property="og:image" content="https://docs.facile.studio/s/1/cover.png">` +
			`<meta name="description" content="The roadmap"></head><body>secret body</body></html>`))
	})
	preview, err := client.Preview(context.Background(), "https://docs.facile.studio/s/1")
	if err != nil {
		t.Fatal(err)
	}
	if preview.Title != "Q3 & plan" || preview.Description != "The roadmap" || preview.Image != "https://docs.facile.studio/s/1/cover.png" {
		t.Fatalf("unexpected preview %+v", preview)
	}
}

func TestPreviewStopsAtARedirectOffTheSuite(t *testing.T) {
	client := suiteServer(t, func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "https://sso.facile.studio/login", http.StatusFound)
	})
	preview, err := client.Preview(context.Background(), "https://feuilles.facile.studio/s/private")
	if err != nil {
		t.Fatal(err)
	}
	if preview.App != "feuilles" || preview.Title != "" {
		t.Fatalf("a redirect off the suite produced %+v", preview)
	}
}

func TestPublicIPRefusesInternalRanges(t *testing.T) {
	for _, raw := range []string{"127.0.0.1", "10.1.2.3", "172.18.0.4", "192.168.1.1", "169.254.169.254", "::1", "fd00::1", "0.0.0.0"} {
		if publicIP(net.ParseIP(raw)) {
			t.Errorf("%s counted as public", raw)
		}
	}
	if !publicIP(net.ParseIP("104.21.3.4")) {
		t.Error("a public address was refused")
	}
}

func TestPreviewIgnoresTheLandingTagsAnSPAServesEverywhere(t *testing.T) {
	client := suiteServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/oembed" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		_, _ = w.Write([]byte(`<head><meta property="og:title" content="Docs · Write together">` +
			`<meta property="og:image" content="https://docs.facile.studio/og-image.png"></head>`))
	})
	preview, err := client.Preview(context.Background(), "https://docs.facile.studio/d/private-id")
	if err != nil {
		t.Fatal(err)
	}
	if preview.Title != "" || preview.Image != "" {
		t.Fatalf("the landing page's tags leaked into a private link card: %+v", preview)
	}
}
