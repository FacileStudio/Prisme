// Package unfurl serves link previews for links to other Facile apps, so a
// message carrying a sheet, a doc, a call or a deck shows a card.
package unfurl

import (
	"context"
	stderrors "errors"

	"github.com/FacileStudio/Prisme/apps/api/internal/unfurl"
	"github.com/FacileStudio/tronc/errors"
)

// Service answers previews from one unfurl client.
type Service struct {
	client *unfurl.Client
}

// NewService builds the preview service.
func NewService(client *unfurl.Client) *Service {
	return &Service{client: client}
}

// Preview returns the card for raw, or invalid_argument for a link Prisme
// does not preview.
func (s *Service) Preview(ctx context.Context, raw string) (Preview, error) {
	preview, err := s.client.Preview(ctx, raw)
	if stderrors.Is(err, unfurl.ErrNotAllowed) {
		return Preview{}, errors.Invalid("only https links to Feuilles, Docs, Echo, Clap, Diapo and Nuage are previewed")
	}
	return preview, err
}
