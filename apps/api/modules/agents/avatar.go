package agents

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log/slog"

	"github.com/FacileStudio/Prisme/apps/api/schemas"
	"github.com/FacileStudio/blob"
)

const (
	blobSize   = 256
	blobPrefix = "blob:"
)

// profile sets the agent's display name and avatar. A caller's picture URL
// wins; without one the agent wears the blob its name draws. Either is
// uploaded only when it differs from what the agent already wears. Failures
// are logged: a missing avatar must not block provisioning.
func (s *Service) profile(ctx context.Context, row *schemas.Agent, avatarURL string) {
	source := avatarURL
	if source == "" {
		source = blobPrefix + row.Name
	}
	avatar := ""
	if row.AvatarURL == "" || row.AvatarSource != source {
		uploaded, err := s.uploadAvatar(ctx, row, avatarURL)
		if err != nil {
			s.deps.Logger.Warn("agent avatar upload failed", slog.String("agent", row.MatrixUserID), slog.Any("error", err))
		}
		avatar = uploaded
	}
	if err := s.deps.Matrix.SetProfile(ctx, row.MatrixUserID, row.Name, avatar); err != nil {
		s.deps.Logger.Warn("agent profile sync failed", slog.String("agent", row.MatrixUserID), slog.Any("error", err))
		return
	}
	if avatar == "" {
		return
	}
	row.AvatarURL, row.AvatarSource = avatar, source
	if err := s.deps.DB.WithContext(ctx).Save(row).Error; err != nil {
		s.deps.Logger.Warn("failed to store the agent avatar", slog.String("agent", row.MatrixUserID), slog.Any("error", err))
	}
}

func (s *Service) uploadAvatar(ctx context.Context, row *schemas.Agent, avatarURL string) (string, error) {
	if avatarURL != "" {
		return s.deps.Matrix.UploadFromURL(ctx, row.MatrixUserID, avatarURL)
	}
	data, err := blobPNG(row.Name)
	if err != nil {
		return "", err
	}
	return s.deps.Matrix.Upload(ctx, row.MatrixUserID, data, "image/png")
}

// blobPNG draws the name's blob centred on a transparent square canvas.
func blobPNG(name string) ([]byte, error) {
	drawn := blob.New(name, blob.Options{}).Pixels(blobSize)
	side := max(drawn.Width, drawn.Height, 1)
	left, top := (side-drawn.Width)/2, (side-drawn.Height)/2
	canvas := image.NewNRGBA(image.Rect(0, 0, side, side))
	for y := range drawn.Height {
		for x := range drawn.Width {
			if pixel := drawn.At(x, y); pixel.Filled {
				canvas.SetNRGBA(left+x, top+y, color.NRGBA{R: pixel.Color.R, G: pixel.Color.G, B: pixel.Color.B, A: 255})
			}
		}
	}
	var out bytes.Buffer
	if err := png.Encode(&out, canvas); err != nil {
		return nil, fmt.Errorf("agents: encode the blob avatar: %w", err)
	}
	return out.Bytes(), nil
}
