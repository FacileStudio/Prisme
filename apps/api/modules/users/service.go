// Package users is the directory of people who have signed in to Prisme.
// Synapse leaves appservice users out of its own user directory, and every
// Prisme human is one, so search is answered here.
package users

import (
	"context"
	"strings"

	"github.com/FacileStudio/tronc/errors"

	"gorm.io/gorm"
)

const maxResults = 20

// User is one directory entry.
type User struct {
	UserID string `json:"user_id"`
	Name   string `json:"name"`
}

// Service searches the directory.
type Service struct {
	db *gorm.DB
}

// NewService builds the directory service.
func NewService(db *gorm.DB) *Service { return &Service{db: db} }

// Search matches the query against names, emails and Matrix ids. An empty
// query lists the first people alphabetically.
func (s *Service) Search(ctx context.Context, query string) ([]User, error) {
	pattern := "%" + escapeLike(strings.ToLower(strings.TrimSpace(query))) + "%"
	out := []User{}
	err := s.db.WithContext(ctx).Table("accounts a").
		Select("a.matrix_user_id AS user_id, COALESCE(NULLIF(u.name, ''), a.localpart) AS name").
		Joins("JOIN porte_users u ON u.id = a.user_id").
		Where("lower(u.name) LIKE ? OR lower(u.email) LIKE ? OR a.matrix_user_id LIKE ?", pattern, pattern, pattern).
		Order("name").Limit(maxResults).Scan(&out).Error
	if err != nil {
		return nil, errors.Internal("failed to search users", err)
	}
	return out, nil
}

func escapeLike(value string) string {
	return strings.NewReplacer(`\`, `\\`, "%", `\%`, "_", `\_`).Replace(value)
}
