// Package nuage lets a Prisme user pick a file from their Nuage drive and
// share it into a chat. Prisme acts as the user, so Nuage decides what they
// can see; the file stays hosted on Nuage and the chat carries a view link.
package nuage

import (
	"context"
	stderrors "errors"
	"log/slog"
	"strings"

	"github.com/FacileStudio/Prisme/apps/api/internal/nuage"
	"github.com/FacileStudio/porte/exchange"
	"github.com/FacileStudio/tronc/errors"
)

const maxQuery = 200

// Service answers the picker for one Nuage deployment, or answers unavailable
// on every route when Prisme has no Nuage.
type Service struct {
	client *nuage.Client
	log    *slog.Logger
}

// NewService builds the picker service. A nil client turns the picker off.
func NewService(client *nuage.Client, logger *slog.Logger) *Service {
	return &Service{client: client, log: logger}
}

// Spaces lists the shared spaces the user belongs to in Nuage.
func (s *Service) Spaces(ctx context.Context, userID int64) (SpacesResponse, error) {
	if s.client == nil {
		return SpacesResponse{}, off()
	}
	spaces, err := s.client.Spaces(ctx, userID)
	if err != nil {
		return SpacesResponse{}, s.fail(err)
	}
	return SpacesResponse{Spaces: nonNil(spaces)}, nil
}

// List returns what sits directly inside a folder, a space or the user's root.
func (s *Service) List(ctx context.Context, userID int64, place nuage.Place) (ListResponse, error) {
	if s.client == nil {
		return ListResponse{}, off()
	}
	folders, files, err := s.client.List(ctx, userID, place)
	if err != nil {
		return ListResponse{}, s.fail(err)
	}
	return ListResponse{Folders: nonNil(folders), Files: nonNil(files)}, nil
}

// Search finds the user's files by name.
func (s *Service) Search(ctx context.Context, userID int64, query string) (SearchResponse, error) {
	query = strings.TrimSpace(query)
	if query == "" || len(query) > maxQuery {
		return SearchResponse{}, errors.Invalid("q must be 1 to 200 characters")
	}
	if s.client == nil {
		return SearchResponse{}, off()
	}
	files, err := s.client.Search(ctx, userID, query)
	if err != nil {
		return SearchResponse{}, s.fail(err)
	}
	return SearchResponse{Files: nonNil(files)}, nil
}

// Share returns a view link to one of the user's files.
func (s *Service) Share(ctx context.Context, userID int64, req ShareRequest) (nuage.Shared, error) {
	if req.FileID <= 0 {
		return nuage.Shared{}, errors.Invalid("file_id is required")
	}
	if s.client == nil {
		return nuage.Shared{}, off()
	}
	shared, err := s.client.Share(ctx, userID, req.FileID)
	if err != nil {
		return nuage.Shared{}, s.fail(err)
	}
	return shared, nil
}

func off() error { return errors.Unavailable("Nuage is not connected to this Prisme") }

func (s *Service) fail(err error) error {
	var refused *exchange.RefusedError
	switch {
	case stderrors.Is(err, nuage.ErrNotFound):
		return errors.NotFound("this file or folder is not in your Nuage")
	case stderrors.Is(err, exchange.ErrReconnect):
		return errors.Failed("sign in to Prisme again to reach your Nuage")
	case stderrors.Is(err, nuage.ErrRefused), stderrors.As(err, &refused):
		s.log.Warn("nuage refused prisme's delegated token", slog.Any("error", err))
		return errors.Forbidden("Nuage refused Prisme's access to your drive")
	}
	s.log.Warn("nuage picker call failed", slog.Any("error", err))
	return errors.Unavailable("Nuage did not answer")
}

func nonNil[T any](items []T) []T {
	if items == nil {
		return []T{}
	}
	return items
}
