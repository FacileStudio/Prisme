// Package users is the directory of people who have signed in to Prisme and
// of their agents. Synapse leaves appservice users out of its own user
// directory, and every Prisme account is one, so search is answered here.
package users

import (
	"context"
	"database/sql"
	"strings"

	"github.com/FacileStudio/tronc/errors"

	"gorm.io/gorm"
)

const maxResults = 20

const searchQuery = `SELECT user_id, name, avatar_url, agent FROM (
	SELECT a.matrix_user_id AS user_id, COALESCE(NULLIF(u.name, ''), a.localpart) AS name,
		a.avatar_url, false AS agent
	FROM accounts a JOIN porte_users u ON u.id = a.user_id
	WHERE lower(u.name) LIKE @pattern OR lower(u.email) LIKE @pattern OR a.matrix_user_id LIKE @pattern
	UNION ALL
	SELECT g.matrix_user_id, g.name, g.avatar_url, true
	FROM agents g
	WHERE g.deactivated_at IS NULL AND (lower(g.name) LIKE @pattern OR g.matrix_user_id LIKE @pattern)
) directory ORDER BY lower(name), user_id LIMIT @limit`

const agentsQuery = `SELECT g.matrix_user_id AS user_id, g.name, a.matrix_user_id AS owner_user_id
	FROM agents g JOIN accounts a ON a.user_id = g.owner_id
	WHERE g.deactivated_at IS NULL ORDER BY lower(g.name), g.matrix_user_id`

// Service searches the directory.
type Service struct {
	db *gorm.DB
}

// NewService builds the directory service.
func NewService(db *gorm.DB) *Service { return &Service{db: db} }

// Search matches the query against names, emails and Matrix ids of people and
// active agents. An empty query lists the first entries alphabetically.
func (s *Service) Search(ctx context.Context, query string) ([]User, error) {
	pattern := "%" + escapeLike(strings.ToLower(strings.TrimSpace(query))) + "%"
	out := []User{}
	err := s.db.WithContext(ctx).
		Raw(searchQuery, sql.Named("pattern", pattern), sql.Named("limit", maxResults)).Scan(&out).Error
	if err != nil {
		return nil, errors.Internal("failed to search users", err)
	}
	return out, nil
}

// Agents lists every active agent with its owner's Matrix id.
func (s *Service) Agents(ctx context.Context) ([]Agent, error) {
	out := []Agent{}
	if err := s.db.WithContext(ctx).Raw(agentsQuery).Scan(&out).Error; err != nil {
		return nil, errors.Internal("failed to list agents", err)
	}
	return out, nil
}

func escapeLike(value string) string {
	return strings.NewReplacer(`\`, `\\`, "%", `\%`, "_", `\_`).Replace(value)
}
