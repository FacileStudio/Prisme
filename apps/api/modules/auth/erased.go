package auth

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"net/http"
	"strings"

	"github.com/FacileStudio/tronc/errors"
)

// subjectErased answers a service call naming a subject Prisme erased. It is
// final: Registre never reissues a subject.
var subjectErased = refusal{http.StatusConflict, "acting_subject_erased", "Prisme erased this Registre subject; it never holds an account again"}

// refuseErased refuses a subject Prisme already erased, before any lookup
// or write: Registre never reissues a subject, so the refusal is final.
func (s serviceKeys) refuseErased(ctx context.Context, subject string, key int) error {
	erased, err := Erased(ctx, s.db, s.issuer, subject)
	if err != nil {
		return errors.Internal("failed to read the erased subjects", err)
	}
	if erased {
		s.refused("a service call named a registre subject prisme erased", 0, key)
		return subjectErased
	}
	return nil
}

// Erased reports whether Prisme erased this subject of the issuer, under
// either spelling of the issuer: the tombstone the erasure writes is
// sha256(issuer without trailing slashes + " " + subject). Every path that
// creates or links an account from a subject asks it first.
func Erased(ctx context.Context, db *sql.DB, issuer, subject string) (bool, error) {
	sum := sha256.Sum256([]byte(strings.TrimRight(issuer, "/") + " " + subject))
	var erased bool
	err := db.QueryRowContext(ctx, `SELECT EXISTS (SELECT 1 FROM erased_subjects WHERE subject_hash = $1)`,
		hex.EncodeToString(sum[:])).Scan(&erased)
	return erased, err
}
