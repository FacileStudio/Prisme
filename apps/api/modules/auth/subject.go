package auth

import (
	"context"
	"database/sql"
	stderrors "errors"
	"log/slog"
	"net/http"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/FacileStudio/porte"
	"github.com/FacileStudio/tronc/errors"
)

// ActingSubjectHeader carries the Registre subject of the person a service
// call acts for, under Prisme's own OIDC issuer. Only a service-key call
// reads it; a session or bearer request that sends it is unaffected.
const ActingSubjectHeader = "X-Prisme-Acting-Subject"

// maxSubjectLength is Registre's own bound on a subject, in bytes.
const maxSubjectLength = 64

// The refusals the header can cause. Each code is stable and each one is
// final for the data as stored: a caller may stop retrying on any of them.
// A failure that could pass on a retry is a 500 instead.
var (
	subjectInvalid     = refusal{http.StatusBadRequest, "acting_subject_invalid", ActingSubjectHeader + " must carry one Registre subject of 1 to 64 bytes with no control characters"}
	subjectUnsupported = refusal{http.StatusPreconditionFailed, "acting_subject_unsupported", ActingSubjectHeader + " needs Registre sign-in, which this Prisme does not have"}
	subjectTaken       = refusal{http.StatusConflict, "acting_subject_taken", "this Registre subject belongs to another Prisme account"}
	subjectMismatch    = refusal{http.StatusConflict, "acting_subject_mismatch", "the Prisme account for this email holds another Registre subject"}
)

// actingSubject reads the header, reporting whether the call sent it. A sent
// header must hold exactly one valid subject.
func actingSubject(r *http.Request) (string, bool, error) {
	values := r.Header.Values(ActingSubjectHeader)
	if len(values) == 0 {
		return "", false, nil
	}
	if len(values) > 1 || !validSubject(values[0]) {
		return "", true, subjectInvalid
	}
	return values[0], true, nil
}

func validSubject(subject string) bool {
	if subject == "" || len(subject) > maxSubjectLength || !utf8.ValidString(subject) {
		return false
	}
	return !strings.ContainsFunc(subject, unicode.IsControl)
}

// attach resolves the account for email as UserFor does and files subject as
// its identity of the issuer, through the identity store a sign-in saves
// with. An account already holding this subject is left alone. When the
// subject belongs to another account, or the account holds another subject
// of the issuer, it refuses before writing anything: a service key never
// re-points or merges anybody's identity.
func (s serviceKeys) attach(ctx context.Context, email, subject string, key int) (int64, error) {
	if s.issuer == "" {
		return 0, subjectUnsupported
	}
	if err := s.refuseErased(ctx, subject, key); err != nil {
		return 0, err
	}
	if userID, err := s.classify(ctx, email, subject, key); userID != 0 || err != nil {
		return userID, err
	}
	userID, err := UserFor(ctx, s.db, email)
	if err != nil {
		return 0, errors.Internal("failed to resolve the user", err)
	}
	err = s.identities.Save(ctx, porte.StoredIdentity{UserID: userID, Provider: s.issuer, Subject: subject})
	if stderrors.Is(err, porte.ErrNotFound) {
		return s.raced(ctx, email, subject, key)
	}
	if err != nil {
		return 0, errors.Internal("failed to link the acting subject", err)
	}
	s.logger.Info("auth: a service call linked a registre subject to an account",
		slog.Int64("user_id", userID), slog.String("service", serviceActor), slog.Int("service_key", key))
	return userID, nil
}

// classify reads the stored rows for this email and subject: the account id
// when it already holds the subject, a refusal when the rows conflict, and
// 0 with no error when the subject may be linked.
func (s serviceKeys) classify(ctx context.Context, email, subject string, key int) (int64, error) {
	existing, holder, err := s.lookup(ctx, email, subject)
	if err != nil {
		return 0, errors.Internal("failed to resolve the acting subject", err)
	}
	if holder != 0 {
		return s.settled(existing, holder, key)
	}
	return 0, s.refuseSecondSubject(ctx, existing, key)
}

// raced classifies a refusal from the identity store, which means another
// write landed after classify: the stored rows decide, and a refusal they do
// not explain is a 500, so the caller retries.
func (s serviceKeys) raced(ctx context.Context, email, subject string, key int) (int64, error) {
	if userID, err := s.classify(ctx, email, subject, key); userID != 0 || err != nil {
		return userID, err
	}
	return 0, errors.Internal("the identity store refused the acting subject", porte.ErrNotFound)
}

// settled answers a subject some account already holds: the account for
// this email holds it, or another one does and the call is refused.
func (s serviceKeys) settled(existing, holder int64, key int) (int64, error) {
	if holder == existing {
		return existing, nil
	}
	s.refused("a service call named a registre subject another account holds", holder, key)
	return 0, subjectTaken
}

func (s serviceKeys) refuseSecondSubject(ctx context.Context, userID int64, key int) error {
	if userID == 0 {
		return nil
	}
	var held bool
	trimmed := strings.TrimSuffix(s.issuer, "/")
	err := s.db.QueryRowContext(ctx, `SELECT EXISTS (SELECT 1 FROM porte_identities
		WHERE user_id = $1 AND provider IN ($2, $3))`, userID, trimmed, trimmed+"/").Scan(&held)
	if err != nil {
		return errors.Internal("failed to read the account's identities", err)
	}
	if held {
		s.refused("a service call named a registre subject the account does not hold", userID, key)
		return subjectMismatch
	}
	return nil
}

// lookup returns the account with this email and the account holding the
// subject under either spelling of the issuer, the exact one first; 0 for
// none.
func (s serviceKeys) lookup(ctx context.Context, email, subject string) (int64, int64, error) {
	var existing, holder int64
	err := s.db.QueryRowContext(ctx, `SELECT id FROM porte_users WHERE email = $1`, email).Scan(&existing)
	if err != nil && !stderrors.Is(err, sql.ErrNoRows) {
		return 0, 0, err
	}
	trimmed := strings.TrimSuffix(s.issuer, "/")
	err = s.db.QueryRowContext(ctx, `SELECT user_id FROM porte_identities
		WHERE provider IN ($1, $2) AND subject = $3 ORDER BY provider = $4 DESC LIMIT 1`,
		trimmed, trimmed+"/", subject, s.issuer).Scan(&holder)
	if err != nil && !stderrors.Is(err, sql.ErrNoRows) {
		return 0, 0, err
	}
	return existing, holder, nil
}
