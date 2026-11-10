package agents

import (
	"context"
	stderrors "errors"
	"strings"

	"github.com/FacileStudio/Prisme/apps/api/internal/accounts"
	"github.com/FacileStudio/Prisme/apps/api/schemas"
	"github.com/FacileStudio/tronc/errors"

	"gorm.io/gorm"
)

const (
	maxCandidates = 50
	noFreeName    = "no free Matrix name for this agent"
)

// identity is who an agent name resolves to: the owner's account, their
// active agent with that name (an empty row when there is none) and the
// localpart that agent should have, empty when every candidate is taken.
type identity struct {
	owner     schemas.Account
	row       schemas.Agent
	localpart string
}

func (s *Service) identify(ctx context.Context, ownerID int64, name, slug string) (identity, error) {
	owner, err := s.deps.Accounts.Ensure(ctx, ownerID)
	if err != nil {
		return identity{}, err
	}
	row, err := s.find(ctx, ownerID, name)
	if err != nil {
		return identity{}, err
	}
	localpart := chooseLocalpart(slug, owner.Localpart, func(candidate string) bool {
		return s.taken(ctx, row.ID, candidate)
	})
	return identity{owner: owner, row: row, localpart: localpart}, nil
}

// find returns the owner's active agent whose trimmed name matches,
// case-insensitively, or an empty row.
func (s *Service) find(ctx context.Context, ownerID int64, name string) (schemas.Agent, error) {
	var row schemas.Agent
	err := s.deps.DB.WithContext(ctx).
		Where("owner_id = ? AND lower(name) = ? AND deactivated_at IS NULL", ownerID, strings.ToLower(strings.TrimSpace(name))).
		Order("id").First(&row).Error
	if err != nil && !stderrors.Is(err, gorm.ErrRecordNotFound) {
		return row, errors.Internal("failed to read the agent", err)
	}
	return row, nil
}

// taken reports whether a localpart belongs to a human account, to any agent
// row other than self, or to an account Synapse deactivated, which an erased
// person's agents leave behind: Matrix never reissues an id.
func (s *Service) taken(ctx context.Context, selfID int64, localpart string) bool {
	db := s.deps.DB.WithContext(ctx)
	var humans, agents int64
	if err := db.Model(&schemas.Account{}).Where("localpart = ?", localpart).Count(&humans).Error; err != nil || humans > 0 {
		return true
	}
	userID := s.deps.Matrix.UserID(localpart)
	err := db.Model(&schemas.Agent{}).
		Where("id <> ? AND (slug = ? OR matrix_user_id = ?)", selfID, localpart, userID).
		Count(&agents).Error
	if err != nil || agents > 0 {
		return true
	}
	retired, err := s.deps.Matrix.Retired(ctx, userID)
	return err != nil || retired
}

// chooseLocalpart picks an agent's Matrix localpart: the slug itself, else the
// slug suffixed with the owner's localpart, then numbered variants of that. It
// returns "" when every candidate is taken.
func chooseLocalpart(slug, ownerLocalpart string, taken func(string) bool) string {
	if !taken(slug) {
		return slug
	}
	base := slug + "-" + ownerLocalpart
	for n := 1; n <= maxCandidates; n++ {
		if candidate := accounts.Candidate(base, n); !taken(candidate) {
			return candidate
		}
	}
	return ""
}
