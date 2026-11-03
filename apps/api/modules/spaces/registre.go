package spaces

import (
	"context"
	"log/slog"
	"time"

	"github.com/FacileStudio/Prisme/apps/api/schemas"
	"github.com/FacileStudio/porte/directory"
)

const (
	retryFirst = 10 * time.Second
	retryMax   = 5 * time.Minute
)

// Registre keeps Prisme's spaces equal to the shared spaces Registre owns.
// issuer is the configured OIDC_ISSUER, the provider column of every porte
// identity Registre's subjects live under.
type Registre struct {
	spaces *Service
	issuer string
	logger *slog.Logger
}

// NewRegistre wires the directory sync over the spaces service.
func NewRegistre(service *Service, issuer string, logger *slog.Logger) *Registre {
	return &Registre{spaces: service, issuer: issuer, logger: logger}
}

// Run imports Prisme's unlinked spaces into Registre, retrying with backoff
// until it succeeds, and only then starts the syncer: a sync before the import
// would remove every member Registre has not heard of yet. It returns when ctx
// ends.
func (r *Registre) Run(ctx context.Context, syncer *directory.Syncer) {
	delay := retryFirst
	for {
		err := r.Import(ctx, syncer.Client)
		if err == nil {
			break
		}
		r.logger.Warn("directory: import into registre failed, retrying", slog.Any("error", err), slog.Duration("retry_in", delay))
		if !wait(ctx, delay) {
			return
		}
		delay = min(delay*2, retryMax)
	}
	syncer.Run(ctx)
}

// Import hands Registre every unlinked space with its members. Registre
// records the import once per client, so a second call after a restart is a
// no-op.
func (r *Registre) Import(ctx context.Context, client directory.Client) error {
	var rows []schemas.Space
	if err := r.spaces.db.WithContext(ctx).Where("registre_space_id IS NULL").Order("created_at, id").Find(&rows).Error; err != nil {
		return err
	}
	spaces := []directory.ImportSpace{}
	for _, row := range rows {
		members, err := r.importMembers(ctx, row)
		if err != nil {
			return err
		}
		if len(members) > 0 {
			spaces = append(spaces, directory.ImportSpace{Name: row.Name, Members: members})
		}
	}
	_, err := client.Import(ctx, spaces)
	return err
}

// importMember is one person joined to a space, with their subject under the
// issuer when they have one.
type importMember struct {
	UserID  int64
	Email   string
	Subject string
}

// importMembers lists the people joined to the space: its creator as owner,
// everyone else as member, named by subject when known and by email otherwise.
// Agents and the Prisme bot hold no account row and never reach Registre.
func (r *Registre) importMembers(ctx context.Context, row schemas.Space) ([]directory.ImportMember, error) {
	joined, err := r.spaces.matrix.JoinedMembers(ctx, r.spaces.bot(), row.RoomID)
	if err != nil {
		return nil, err
	}
	var people []importMember
	err = r.spaces.db.WithContext(ctx).Table("accounts a").
		Select("a.user_id, u.email, COALESCE(i.subject, '') AS subject").
		Joins("JOIN porte_users u ON u.id = a.user_id").
		Joins("LEFT JOIN porte_identities i ON i.user_id = a.user_id AND i.provider = ?", r.issuer).
		Where("a.matrix_user_id IN ?", append(joined, "")).Order("a.user_id").Scan(&people).Error
	if err != nil {
		return nil, err
	}
	members := make([]directory.ImportMember, 0, len(people))
	for _, person := range people {
		member := directory.ImportMember{Subject: person.Subject, Role: directory.RoleMember}
		if person.UserID == row.CreatedBy {
			member.Role = directory.RoleOwner
		}
		if member.Subject == "" {
			member.Email = person.Email
		}
		members = append(members, member)
	}
	return members, nil
}

// wait sleeps for d or until ctx ends, reporting whether it slept fully.
func wait(ctx context.Context, d time.Duration) bool {
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return false
	case <-timer.C:
		return true
	}
}
