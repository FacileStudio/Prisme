// Package gifs is the GIF picker's search, run from the server through GIPHY so
// no user's address or session reaches GIPHY. It is off when no key is set.
package gifs

import (
	"context"
	stderrors "errors"
	"log/slog"
	"net/url"
	"strings"

	"github.com/FacileStudio/Prisme/apps/api/internal/gifs"
	"github.com/FacileStudio/tronc/errors"
)

const maxQuery = 100

// Service searches GIPHY and relays its media.
type Service struct {
	client *gifs.Client
	log    *slog.Logger
}

// NewService builds the GIF search. A nil client turns it off.
func NewService(client *gifs.Client, logger *slog.Logger) *Service {
	return &Service{client: client, log: logger}
}

// Enabled reports whether a GIPHY key is configured.
func (s *Service) Enabled() bool { return s != nil && s.client != nil }

// Search finds GIFs for a query; an empty query lists the trending ones.
func (s *Service) Search(ctx context.Context, query, lang string, offset int) (PageResponse, error) {
	if !s.Enabled() {
		return PageResponse{}, errors.Unavailable("GIF search is not configured")
	}
	query = strings.TrimSpace(query)
	if len([]rune(query)) > maxQuery {
		return PageResponse{}, errors.Invalid("q is too long")
	}
	var page gifs.Page
	var err error
	if query == "" {
		page, err = s.client.Trending(ctx, offset)
	} else {
		page, err = s.client.Search(ctx, query, language(lang), offset)
	}
	if err != nil {
		s.log.Warn("GIPHY search failed", slog.Any("error", err))
		return PageResponse{}, errors.Unavailable("GIPHY did not answer")
	}
	out := PageResponse{GIFs: make([]GIF, 0, len(page.GIFs)), Next: page.Next}
	for _, g := range page.GIFs {
		out.GIFs = append(out.GIFs, GIF{ID: g.ID, Title: g.Title, Preview: image(g.Preview), Full: image(g.Full)})
	}
	return out, nil
}

// Media fetches one GIPHY media file named by its URL.
func (s *Service) Media(ctx context.Context, raw string) (gifs.Media, error) {
	if !s.Enabled() {
		return gifs.Media{}, errors.Unavailable("GIF search is not configured")
	}
	media, err := s.client.Fetch(ctx, raw)
	if stderrors.Is(err, gifs.ErrHost) {
		return gifs.Media{}, errors.Invalid("not a GIPHY media url")
	}
	if err != nil {
		s.log.Warn("GIPHY media failed", slog.Any("error", err))
		return gifs.Media{}, errors.Unavailable("GIPHY did not answer")
	}
	return media, nil
}

// MediaPath is the Prisme path that serves a GIPHY media URL.
func MediaPath(raw string) string {
	return "/gifs/media?u=" + url.QueryEscape(raw)
}

func language(lang string) string {
	lang = strings.ToLower(strings.TrimSpace(lang))
	switch lang {
	case "fr", "es", "de":
		return lang
	}
	return "en"
}
