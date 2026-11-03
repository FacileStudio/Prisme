// Package spaces provisions Matrix spaces and their channels through the
// appservice. The creator owns each room; the Prisme bot is a second admin so
// it can link channels into a space whoever creates them.
package spaces

import (
	"context"
	stderrors "errors"
	"log/slog"
	"strings"

	"github.com/FacileStudio/Prisme/apps/api/internal/accounts"
	"github.com/FacileStudio/Prisme/apps/api/internal/matrix"
	"github.com/FacileStudio/Prisme/apps/api/schemas"
	"github.com/FacileStudio/tronc/errors"

	"gorm.io/gorm"
)

const botLocalpart = "prisme"

// Followers moves a user's agents along with the user. The agents module
// implements it; spaces cannot import agents, which imports spaces.
type Followers interface {
	Follow(ctx context.Context, ownerID int64, roomIDs []string) error
	Unfollow(ctx context.Context, ownerID int64, roomIDs []string) error
}

// Deps holds the services spaces rely on. Followers may be nil.
type Deps struct {
	DB        *gorm.DB
	Matrix    *matrix.Client
	Accounts  *accounts.Service
	Followers Followers
	Logger    *slog.Logger
}

// Service lists and creates spaces and channels.
type Service struct {
	db        *gorm.DB
	matrix    *matrix.Client
	accounts  *accounts.Service
	followers Followers
	logger    *slog.Logger
}

// NewService builds the spaces service.
func NewService(deps Deps) *Service {
	return &Service{db: deps.DB, matrix: deps.Matrix, accounts: deps.Accounts, followers: deps.Followers, logger: deps.Logger}
}

// List returns the spaces the user is joined to, with their channels.
func (s *Service) List(ctx context.Context, userID int64) ([]Space, error) {
	account, err := s.accounts.Ensure(ctx, userID)
	if err != nil {
		return nil, err
	}
	joined, err := s.matrix.JoinedRooms(ctx, account.MatrixUserID)
	if err != nil {
		return nil, errors.Internal("failed to list Matrix rooms", err)
	}
	var rows []schemas.Space
	if err := s.db.WithContext(ctx).Where("room_id IN ?", append(joined, "")).Order("name").Find(&rows).Error; err != nil {
		return nil, errors.Internal("failed to list spaces", err)
	}
	out := make([]Space, 0, len(rows))
	for _, row := range rows {
		channels, err := s.channels(ctx, row.ID)
		if err != nil {
			return nil, err
		}
		out = append(out, toSpace(row, channels))
	}
	return out, nil
}

// Create makes a space owned by the user, with a #general channel, and brings
// the user's agents into both.
func (s *Service) Create(ctx context.Context, userID int64, req CreateSpaceRequest) (Space, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" || len(name) > 80 {
		return Space{}, errors.Invalid("a space needs a name of 1 to 80 characters")
	}
	account, err := s.accounts.Ensure(ctx, userID)
	if err != nil {
		return Space{}, err
	}
	space, err := s.provision(ctx, account, name, nil)
	if err != nil {
		return Space{}, err
	}
	s.follow(ctx, userID, roomsOf(space))
	return space, nil
}

// Get returns one space the user is joined to.
func (s *Service) Get(ctx context.Context, userID, spaceID int64) (schemas.Space, schemas.Account, error) {
	var row schemas.Space
	err := s.db.WithContext(ctx).First(&row, "id = ?", spaceID).Error
	if stderrors.Is(err, gorm.ErrRecordNotFound) {
		return row, schemas.Account{}, errors.NotFound("no such space")
	}
	if err != nil {
		return row, schemas.Account{}, errors.Internal("failed to read the space", err)
	}
	account, err := s.accounts.Ensure(ctx, userID)
	if err != nil {
		return row, account, err
	}
	if !s.matrix.IsJoined(ctx, account.MatrixUserID, row.RoomID) {
		return row, account, errors.NotFound("no such space")
	}
	return row, account, nil
}

func (s *Service) bot() string { return s.matrix.UserID(botLocalpart) }

func (s *Service) createRoom(ctx context.Context, creator string, body matrix.CreateRoomRequest) (string, error) {
	roomID, err := s.matrix.CreateRoom(ctx, creator, body)
	if err != nil {
		return "", errors.Internal("failed to create the Matrix room", err)
	}
	if err := s.matrix.Join(ctx, s.bot(), roomID); err != nil {
		return "", errors.Internal("the Prisme bot could not join the room", err)
	}
	return roomID, nil
}

func (s *Service) freeSlug(ctx context.Context, base string) string {
	for n := 1; ; n++ {
		slug := accounts.Candidate(base, n)
		var count int64
		s.db.WithContext(ctx).Model(&schemas.Space{}).Where("slug = ?", slug).Count(&count)
		if count == 0 {
			return slug
		}
	}
}
