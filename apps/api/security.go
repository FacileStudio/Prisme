package main

import (
	"log/slog"
	"net/url"

	"github.com/FacileStudio/tronc/middleware"
	"github.com/FacileStudio/tronc/spa"
)

var callFeatures = []string{"camera", "microphone", "display-capture", "speaker-selection"}

// origins returns the http(s) origin of rawURL and its websocket twin, or
// nothing when rawURL is empty or has no host.
func origins(rawURL string) []string {
	parsed, err := url.Parse(rawURL)
	if err != nil || parsed.Host == "" {
		return nil
	}
	secure := parsed.Scheme == "https" || parsed.Scheme == "wss"
	http, ws := "http", "ws"
	if secure {
		http, ws = "https", "wss"
	}
	return []string{http + "://" + parsed.Host, ws + "://" + parsed.Host}
}

// securityConfig extends tronc's default security headers with what Prisme
// loads: the homeserver and the LiveKit server in connect-src, blob and https
// media, the call permissions, and the inline scripts adapter-static writes
// into the client's pages. A client directory that cannot be hashed is logged
// and leaves the default script policy, which only reports.
func securityConfig(clientDir, homeserver, liveKit string, logger *slog.Logger) middleware.SecurityConfig {
	connect := append(origins(homeserver), origins(liveKit)...)
	sources := map[string][]string{
		"connect-src": connect,
		"media-src":   {"blob:", "https:"},
		"worker-src":  {"blob:"},
	}
	if spa.Available(clientDir) {
		hashes, err := spa.InlineScriptHashes(clientDir)
		if err != nil {
			logger.Warn("failed to hash the client's inline scripts", slog.Any("error", err))
		} else {
			sources["script-src"] = hashes
		}
	}
	return middleware.SecurityConfig{CSPSources: sources, Allow: callFeatures}
}
