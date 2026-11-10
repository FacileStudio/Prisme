package erasure

import (
	"context"
	"database/sql"
	"fmt"
)

const personMatrixIDs = `SELECT matrix_user_id FROM accounts WHERE user_id = $1
	UNION SELECT matrix_user_id FROM agents WHERE owner_id = $1`

var rowStatements = []string{
	`INSERT INTO erased_subjects (subject_hash, erased_at)
		SELECT encode(sha256(convert_to(rtrim(provider, '/') || ' ' || subject, 'UTF8')), 'hex'), now()
		FROM porte_identities WHERE user_id = $1
		ON CONFLICT (subject_hash) DO NOTHING`,
	`UPDATE calls SET started_by = '' WHERE started_by IN (` + personMatrixIDs + `)`,
	`UPDATE spaces SET created_by = 0 WHERE created_by = $1`,
	`DELETE FROM matrix_sessions WHERE user_id = $1`,
	`DELETE FROM agents WHERE owner_id = $1`,
	`DELETE FROM accounts WHERE user_id = $1`,
}

// eraseRows is the kit's EraseData hook: inside the transaction that then
// deletes the porte_users row, whose cascade takes the identities, sessions,
// login codes, magic links and push devices. It records the erased subjects,
// lifts the person off shared rows (a running call, the spaces they created,
// which stay for their other members) and deletes what is theirs alone.
func eraseRows(ctx context.Context, tx *sql.Tx, userID int64) error {
	for _, statement := range rowStatements {
		if _, err := tx.ExecContext(ctx, statement, userID); err != nil {
			return fmt.Errorf("erasure: %w", err)
		}
	}
	return nil
}

// EraseUnlinked erases an account that holds no Registre subject, found by its
// email: the one kind of account a deletion notice cannot reach, because the
// notice carries only the subject. An account linked to a subject is refused;
// Registre's notice erases that one. It reports whether an account was found.
func (e *Eraser) EraseUnlinked(ctx context.Context, email string) (bool, error) {
	userID, err := e.unlinkedAccount(ctx, email)
	if err != nil || userID == 0 {
		return false, err
	}
	if err := e.eraseMatrix(ctx, userID); err != nil {
		return true, err
	}
	tx, err := e.deps.SQL.BeginTx(ctx, nil)
	if err != nil {
		return true, err
	}
	defer func() { _ = tx.Rollback() }()
	if err := eraseRows(ctx, tx, userID); err != nil {
		return true, err
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM porte_users WHERE id = $1`, userID); err != nil {
		return true, err
	}
	return true, tx.Commit()
}

func (e *Eraser) unlinkedAccount(ctx context.Context, email string) (int64, error) {
	var rows []struct {
		ID     int64
		Linked bool
	}
	err := e.deps.DB.WithContext(ctx).Raw(`SELECT u.id,
		EXISTS (SELECT 1 FROM porte_identities i WHERE i.user_id = u.id AND i.provider <> 'local') AS linked
		FROM porte_users u WHERE lower(u.email) = lower(?)`, email).Scan(&rows).Error
	if err != nil || len(rows) == 0 {
		return 0, err
	}
	if rows[0].Linked {
		return 0, fmt.Errorf("erasure: the account is linked to a registre identity; registre's deletion notice erases it")
	}
	return rows[0].ID, nil
}
