package spaces

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/FacileStudio/Prisme/apps/api/modules/auth"
	"github.com/FacileStudio/Prisme/apps/api/schemas"
	"github.com/FacileStudio/porte/directory"
)

// person is a snapshot member resolved to a porte user.
type person struct {
	userID int64
	role   string
}

// resolveAll resolves the snapshot's members, dropping those it must skip.
func (r *Registre) resolveAll(ctx context.Context, members []directory.Member) ([]person, error) {
	out := make([]person, 0, len(members))
	for _, member := range members {
		userID, err := r.resolve(ctx, member)
		if err != nil {
			return nil, fmt.Errorf("resolve member %s: %w", member.Subject, err)
		}
		if userID != 0 {
			out = append(out, person{userID: userID, role: member.Role})
		}
	}
	return out, nil
}

// resolve returns the porte user for a snapshot member, or 0 to skip them: the
// user holding the identity (issuer, subject), else the user with this email,
// created the way a service-key call acting for them would create it.
func (r *Registre) resolve(ctx context.Context, member directory.Member) (int64, error) {
	if member.Subject == "" {
		return 0, nil
	}
	var ids []int64
	err := r.spaces.db.WithContext(ctx).Table("porte_identities").
		Where("provider = ? AND subject = ?", r.issuer, member.Subject).Limit(1).Pluck("user_id", &ids).Error
	if err != nil || len(ids) > 0 {
		return first(ids), err
	}
	email := strings.ToLower(strings.TrimSpace(member.Email))
	if email == "" {
		return 0, nil
	}
	return r.adopt(ctx, email, member)
}

// adopt finds or creates the user with this email and attaches the member's
// identity, with Registre's name and picture where the user has none. A user
// who already holds another subject under the issuer is skipped and logged
// rather than given a second one.
func (r *Registre) adopt(ctx context.Context, email string, member directory.Member) (int64, error) {
	sqlDB, err := r.spaces.db.DB()
	if err != nil {
		return 0, err
	}
	userID, err := auth.UserFor(ctx, sqlDB, email)
	if err != nil {
		return 0, err
	}
	db := r.spaces.db.WithContext(ctx)
	var held int64
	if err := db.Table("porte_identities").Where("provider = ? AND user_id = ?", r.issuer, userID).Count(&held).Error; err != nil {
		return 0, err
	}
	if held > 0 {
		r.logger.Warn("directory: skipped a member whose email belongs to another registre identity", slog.Int64("user_id", userID))
		return 0, nil
	}
	if err := db.Exec(`INSERT INTO porte_identities (user_id, provider, subject) VALUES (?, ?, ?)
		ON CONFLICT (provider, subject) DO NOTHING`, userID, r.issuer, member.Subject).Error; err != nil {
		return 0, err
	}
	return userID, db.Exec(`UPDATE porte_users SET
		name = CASE WHEN name = '' THEN ? ELSE name END,
		avatar_url = CASE WHEN avatar_url = '' THEN ? ELSE avatar_url END
		WHERE id = ?`, strings.TrimSpace(member.Name), httpsOnly(member.Picture), userID).Error
}

// accountsIn maps the users who hold a Prisme account among the Matrix ids to
// those ids. The bot and agents hold none.
func (r *Registre) accountsIn(ctx context.Context, matrixIDs []string) (map[int64]string, error) {
	var rows []schemas.Account
	err := r.spaces.db.WithContext(ctx).Where("matrix_user_id IN ?", append(matrixIDs, "")).Find(&rows).Error
	if err != nil {
		return nil, err
	}
	out := make(map[int64]string, len(rows))
	for _, row := range rows {
		out[row.UserID] = row.MatrixUserID
	}
	return out, nil
}

// rooms is the space room followed by its channels.
func (r *Registre) rooms(ctx context.Context, row schemas.Space) ([]string, error) {
	channels, err := r.spaces.channels(ctx, row.ID)
	if err != nil {
		return nil, err
	}
	return roomsOf(Space{RoomID: row.RoomID, Channels: channels}), nil
}

func first(ids []int64) int64 {
	if len(ids) == 0 {
		return 0
	}
	return ids[0]
}

func httpsOnly(picture string) string {
	if strings.HasPrefix(picture, "https://") {
		return picture
	}
	return ""
}
