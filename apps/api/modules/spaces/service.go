// Package spaces provisions Matrix spaces and their channels through the
// appservice. The creator owns each room; the Prisme bot is a second admin so
// it can link channels into a space whoever creates them.
package spaces

import (
	"context"
	stderrors "errors"
	"strings"

	"github.com/FacileStudio/Prisme/apps/api/internal/accounts"
	"github.com/FacileStudio/Prisme/apps/api/internal/matrix"
	"github.com/FacileStudio/Prisme/apps/api/schemas"
	"github.com/FacileStudio/tronc/errors"

	"gorm.io/gorm"
)

const botLocalpart = "prisme"

// Service lists and creates spaces and channels.
type Service struct {
	db       *gorm.DB
	matrix   *matrix.Client
	accounts *accounts.Service
}

// NewService builds the spaces service.
func NewService(db *gorm.DB, client *matrix.Client, accountService *accounts.Service) *Service {
	return &Service{db: db, matrix: client, accounts: accountService}
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
		out = append(out, Space{ID: row.ID, RoomID: row.RoomID, Name: row.Name, Slug: row.Slug, Channels: channels})
	}
	return out, nil
}

// Create makes a space owned by the user, with a #general channel.
func (s *Service) Create(ctx context.Context, userID int64, req CreateSpaceRequest) (Space, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" || len(name) > 80 {
		return Space{}, errors.Invalid("a space needs a name of 1 to 80 characters")
	}
	account, err := s.accounts.Ensure(ctx, userID)
	if err != nil {
		return Space{}, err
	}
	row := schemas.Space{Name: name, CreatedBy: userID, Slug: s.freeSlug(ctx, Slugify(name, "space"))}
	row.RoomID, err = s.createRoom(ctx, account.MatrixUserID, spaceRoom(name, s.bot()))
	if err != nil {
		return Space{}, err
	}
	if err := s.db.WithContext(ctx).Create(&row).Error; err != nil {
		return Space{}, errors.Internal("failed to store the space", err)
	}
	general, err := s.addChannel(ctx, row, account.MatrixUserID, CreateChannelRequest{Name: "general"})
	if err != nil {
		return Space{}, err
	}
	return Space{ID: row.ID, RoomID: row.RoomID, Name: row.Name, Slug: row.Slug, Channels: []Channel{general}}, nil
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
