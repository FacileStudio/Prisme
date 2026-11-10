package erasure

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

// ErrUnproven is a deletion notice for a subject no identity row holds while
// some account holds no subject at all: the person may be one of those, so
// the notice must not be answered as done.
var ErrUnproven = errors.New("erasure: cannot prove the subject never held an account")

const tombstoneHash = `encode(sha256(convert_to(rtrim(?, '/') || ' ' || ?, 'UTF8')), 'hex')`

// owner returns the account holding the identity under either spelling of the
// issuer, the exact one first, or 0.
func (e *Eraser) owner(ctx context.Context, provider, subject string) (int64, error) {
	var ids []int64
	trimmed := strings.TrimSuffix(provider, "/")
	err := e.deps.DB.WithContext(ctx).Raw(`SELECT user_id FROM porte_identities
		WHERE provider IN (?, ?) AND subject = ? ORDER BY provider = ? DESC LIMIT 1`,
		trimmed, trimmed+"/", subject, provider).Scan(&ids).Error
	if err != nil || len(ids) == 0 {
		return 0, err
	}
	return ids[0], nil
}

// proveUnknown answers nil for a subject Prisme already erased, or when every
// account is linked to a subject of this issuer, so an unknown subject cannot
// be anyone here. Otherwise it is ErrUnproven and the notice stays pending.
func (e *Eraser) proveUnknown(ctx context.Context, provider, subject string) error {
	var erased int64
	err := e.deps.DB.WithContext(ctx).Table("erased_subjects").
		Where("subject_hash = "+tombstoneHash, provider, subject).Count(&erased).Error
	if err != nil || erased > 0 {
		return err
	}
	unlinked, err := e.unlinked(ctx, provider)
	if err != nil || unlinked == 0 {
		return err
	}
	return fmt.Errorf("%w: %d accounts carry no subject of this issuer", ErrUnproven, unlinked)
}

// unlinked counts the accounts with no identity under either spelling of the
// issuer: people a service key or an early sign-in created by email alone.
func (e *Eraser) unlinked(ctx context.Context, provider string) (int64, error) {
	trimmed := strings.TrimSuffix(provider, "/")
	var count int64
	err := e.deps.DB.WithContext(ctx).Raw(`SELECT count(*) FROM porte_users u WHERE NOT EXISTS (
		SELECT 1 FROM porte_identities i WHERE i.user_id = u.id AND i.provider IN (?, ?))`, trimmed, trimmed+"/").
		Scan(&count).Error
	return count, err
}
