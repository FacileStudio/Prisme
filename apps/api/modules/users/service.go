// Package users is the directory of people who have signed in to Prisme and
// of their agents. Synapse leaves appservice users out of its own user
// directory, and every Prisme account is one, so search is answered here.
//
// The directory shows a caller only the people and agents it shares a space
// with, plus its own agents. Anyone can sign in through Registre, so a list of
// every account would hand a stranger the whole user base.
package users

import (
	"context"
	"database/sql"
	"maps"
	"slices"
	"strings"

	"github.com/FacileStudio/Prisme/apps/api/schemas"
	"github.com/FacileStudio/tronc/errors"

	"gorm.io/gorm"
)

const maxResults = 200

const searchQuery = `SELECT user_id, name, avatar_url, agent FROM (
	SELECT a.matrix_user_id AS user_id, COALESCE(NULLIF(u.name, ''), a.localpart) AS name,
		a.avatar_url, false AS agent, a.user_id AS owner_id
	FROM accounts a JOIN porte_users u ON u.id = a.user_id
	WHERE lower(u.name) LIKE @pattern OR lower(u.email) LIKE @pattern OR a.matrix_user_id LIKE @pattern
	UNION ALL
	SELECT g.matrix_user_id, g.name, g.avatar_url, true, g.owner_id
	FROM agents g
	WHERE g.deactivated_at IS NULL AND (lower(g.name) LIKE @pattern OR g.matrix_user_id LIKE @pattern)
) directory WHERE user_id IN @peers OR (agent AND owner_id = @caller)
ORDER BY lower(name), user_id LIMIT @limit`

const agentsQuery = `SELECT g.matrix_user_id AS user_id, g.name, a.matrix_user_id AS owner_user_id
	FROM agents g JOIN accounts a ON a.user_id = g.owner_id
	WHERE g.deactivated_at IS NULL AND a.matrix_user_id IN @peers
		AND (g.matrix_user_id IN @peers OR g.owner_id = @caller)
	ORDER BY lower(g.name), g.matrix_user_id`

// Rooms is what the directory reads from Synapse to find the caller's space
// peers. *matrix.Client is one.
type Rooms interface {
	JoinedRooms(ctx context.Context, userID string) ([]string, error)
	JoinedMembers(ctx context.Context, asUser, roomID string) ([]string, error)
}

// Service searches the directory.
type Service struct {
	db    *gorm.DB
	rooms Rooms
}

// NewService builds the directory service.
func NewService(db *gorm.DB, rooms Rooms) *Service { return &Service{db: db, rooms: rooms} }

// Search matches the query against names, emails and Matrix ids of the
// people and active agents callerID may see. An empty query lists the first
// of them alphabetically.
func (s *Service) Search(ctx context.Context, callerID int64, query string) ([]User, error) {
	out := []User{}
	peers, err := s.peers(ctx, callerID)
	if err != nil || len(peers) == 0 {
		return out, err
	}
	pattern := "%" + escapeLike(strings.ToLower(strings.TrimSpace(query))) + "%"
	err = s.db.WithContext(ctx).Raw(searchQuery, sql.Named("pattern", pattern), sql.Named("peers", peers),
		sql.Named("caller", callerID), sql.Named("limit", maxResults)).Scan(&out).Error
	if err != nil {
		return nil, errors.Internal("failed to search users", err)
	}
	return out, nil
}

// Agents lists the active agents callerID may see, each with its owner's
// Matrix id: its own, and those of people it shares a space with.
func (s *Service) Agents(ctx context.Context, callerID int64) ([]Agent, error) {
	out := []Agent{}
	peers, err := s.peers(ctx, callerID)
	if err != nil || len(peers) == 0 {
		return out, err
	}
	err = s.db.WithContext(ctx).Raw(agentsQuery, sql.Named("peers", peers), sql.Named("caller", callerID)).Scan(&out).Error
	if err != nil {
		return nil, errors.Internal("failed to list agents", err)
	}
	return out, nil
}

// peers returns the Matrix ids joined to any Prisme space callerID is joined
// to, the caller's own included. A caller without a Matrix account yet has no
// peers.
func (s *Service) peers(ctx context.Context, callerID int64) ([]string, error) {
	var accounts []schemas.Account
	if err := s.db.WithContext(ctx).Where("user_id = ?", callerID).Limit(1).Find(&accounts).Error; err != nil {
		return nil, errors.Internal("failed to read the account", err)
	}
	if len(accounts) == 0 {
		return nil, nil
	}
	caller := accounts[0].MatrixUserID
	joined, err := s.rooms.JoinedRooms(ctx, caller)
	if err != nil {
		return nil, errors.Internal("failed to list Matrix rooms", err)
	}
	var spaces []string
	err = s.db.WithContext(ctx).Model(&schemas.Space{}).Where("room_id IN ?", append(joined, "")).Pluck("room_id", &spaces).Error
	if err != nil {
		return nil, errors.Internal("failed to list spaces", err)
	}
	seen := map[string]struct{}{caller: {}}
	for _, space := range spaces {
		members, err := s.rooms.JoinedMembers(ctx, caller, space)
		if err != nil {
			return nil, errors.Internal("failed to list space members", err)
		}
		for _, member := range members {
			seen[member] = struct{}{}
		}
	}
	return slices.Sorted(maps.Keys(seen)), nil
}

func escapeLike(value string) string {
	return strings.NewReplacer(`\`, `\\`, "%", `\%`, "_", `\_`).Replace(value)
}
