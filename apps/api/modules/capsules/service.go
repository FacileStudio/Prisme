// Package capsules lets a Prisme user seal a message as a Capsule that burns on
// its first read. The browser encrypts and decrypts; Prisme relays ciphertext
// to Capsule and never holds a key.
package capsules

import (
	"context"
	stderrors "errors"
	"log/slog"

	"github.com/FacileStudio/Prisme/apps/api/internal/capsule"
	"github.com/FacileStudio/tronc/errors"
)

// Service relays capsules to one Capsule deployment.
type Service struct {
	capsule *capsule.Client
	log     *slog.Logger
}

// NewService builds the capsule relay.
func NewService(client *capsule.Client, logger *slog.Logger) *Service {
	return &Service{capsule: client, log: logger}
}

// Seal stores ciphertext as a burn-after-read capsule.
func (s *Service) Seal(ctx context.Context, req SealRequest) (SealResponse, error) {
	if req.Ciphertext == "" {
		return SealResponse{}, errors.Invalid("ciphertext is required")
	}
	if !capsule.ValidExpiry(req.ExpiresIn) {
		return SealResponse{}, errors.Invalid("expires_in must be one of: 1h, 24h, 7d, 30d")
	}
	sealed, err := s.capsule.Seal(ctx, req.Ciphertext, req.ExpiresIn)
	if err != nil {
		return SealResponse{}, s.unavailable(err)
	}
	return SealResponse{ID: sealed.ID, URL: s.capsule.BaseURL() + "/" + sealed.ID, ExpiresAt: sealed.ExpiresAt}, nil
}

// Meta reports whether a capsule can still be opened.
func (s *Service) Meta(ctx context.Context, id string) (MetaResponse, error) {
	if !capsule.ValidID(id) {
		return MetaResponse{}, errors.Invalid("not a capsule id")
	}
	meta, err := s.capsule.Meta(ctx, id)
	if stderrors.Is(err, capsule.ErrNotFound) {
		return MetaResponse{ID: id}, nil
	}
	if err != nil {
		return MetaResponse{}, s.unavailable(err)
	}
	return MetaResponse{ID: id, Exists: meta.Exists, ExpiresAt: meta.ExpiresAt}, nil
}

// Open returns a capsule's ciphertext and burns it.
func (s *Service) Open(ctx context.Context, id string) (OpenResponse, error) {
	if !capsule.ValidID(id) {
		return OpenResponse{}, errors.Invalid("not a capsule id")
	}
	ciphertext, err := s.capsule.Open(ctx, id)
	if stderrors.Is(err, capsule.ErrNotFound) {
		return OpenResponse{}, errors.NotFound("this capsule was already opened or has expired")
	}
	if err != nil {
		return OpenResponse{}, s.unavailable(err)
	}
	return OpenResponse{Ciphertext: ciphertext}, nil
}

func (s *Service) unavailable(err error) error {
	s.log.Warn("capsule relay failed", slog.Any("error", err))
	return errors.Unavailable("Capsule did not answer")
}
