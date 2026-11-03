package auth

import (
	"context"
	"strings"

	"github.com/FacileStudio/porte"
	"github.com/FacileStudio/porte/pg"
)

// pictureUsers copies Registre's picture claim into the avatar porte/pg
// stores. porte/pg writes only Claims.AvatarURL, which porte fills when an
// AvatarStore is wired; Prisme wires none because Matrix keeps the image.
type pictureUsers struct {
	*pg.UserStore
}

// UpsertFromOIDC stores the https picture claim as the avatar, then upserts.
func (u pictureUsers) UpsertFromOIDC(ctx context.Context, claims porte.Claims) (int64, error) {
	if claims.AvatarURL == "" && strings.HasPrefix(claims.Picture, "https://") {
		claims.AvatarURL = claims.Picture
	}
	return u.UserStore.UpsertFromOIDC(ctx, claims)
}
