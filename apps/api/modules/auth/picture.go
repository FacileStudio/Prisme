package auth

import (
	"context"
	"database/sql"
	"strings"

	"github.com/FacileStudio/porte"
	"github.com/FacileStudio/porte/pg"
	"github.com/FacileStudio/tronc/errors"
)

// pictureUsers copies Registre's picture claim into the avatar porte/pg
// stores. porte/pg writes only Claims.AvatarURL, which porte fills when an
// AvatarStore is wired; Prisme wires none because Matrix keeps the image.
// It is also the one door porte's sign-in and bearer provisioning create an
// account through, so it refuses a subject Prisme already erased.
type pictureUsers struct {
	*pg.UserStore
	db *sql.DB
}

// UpsertFromOIDC refuses an erased subject, stores the https picture claim
// as the avatar, then upserts.
func (u pictureUsers) UpsertFromOIDC(ctx context.Context, claims porte.Claims) (int64, error) {
	erased, err := Erased(ctx, u.db, claims.Provider, claims.Subject)
	if err != nil {
		return 0, errors.Internal("failed to read the erased subjects", err)
	}
	if erased {
		return 0, errors.Forbidden("this Registre account was erased from Prisme")
	}
	if claims.AvatarURL == "" && strings.HasPrefix(claims.Picture, "https://") {
		claims.AvatarURL = claims.Picture
	}
	return u.UserStore.UpsertFromOIDC(ctx, claims)
}
