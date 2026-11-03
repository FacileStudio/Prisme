// Package agents gives a user bot accounts (Bulle assistants) in the agent
// namespace, @agent-<owner>-<name>, and places them in the owner's rooms.
package agents

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	stderrors "errors"
	"strings"

	"github.com/FacileStudio/Prisme/apps/api/internal/accounts"
	"github.com/FacileStudio/Prisme/apps/api/internal/matrix"
	"github.com/FacileStudio/Prisme/apps/api/internal/seal"
	"github.com/FacileStudio/Prisme/apps/api/modules/spaces"
	"github.com/FacileStudio/Prisme/apps/api/schemas"
	"github.com/FacileStudio/tronc/errors"

	"gorm.io/gorm"
)

// Deps is what the agent service is built from.
type Deps struct {
	DB         *gorm.DB
	Matrix     *matrix.Client
	Accounts   *accounts.Service
	Box        *seal.Box
	Homeserver string
}

// Service creates, places and deactivates agents.
type Service struct {
	deps Deps
}

// NewService builds the agent service.
func NewService(deps Deps) *Service { return &Service{deps: deps} }

// Create returns the owner's agent with that name, creating it on first call,
// and always refreshes its display name and avatar.
func (s *Service) Create(ctx context.Context, ownerID int64, req CreateRequest) (Login, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" || len(name) > 80 {
		return Login{}, errors.Invalid("an agent needs a name of 1 to 80 characters")
	}
	owner, err := s.deps.Accounts.Ensure(ctx, ownerID)
	if err != nil {
		return Login{}, err
	}
	slug := spaces.Slugify(name, "agent")
	row, err := s.find(ctx, ownerID, slug)
	if err != nil {
		return Login{}, err
	}
	if row.ID == 0 || row.DeactivatedAt != nil {
		row, err = s.provision(ctx, owner, row, slug)
		if err != nil {
			return Login{}, err
		}
	}
	row.Name = name
	login, err := s.login(ctx, &row)
	if err != nil {
		return Login{}, err
	}
	s.profile(ctx, row, req.AvatarURL)
	return login, nil
}

func (s *Service) find(ctx context.Context, ownerID int64, slug string) (schemas.Agent, error) {
	var row schemas.Agent
	err := s.deps.DB.WithContext(ctx).First(&row, "owner_id = ? AND slug = ?", ownerID, slug).Error
	if err != nil && !stderrors.Is(err, gorm.ErrRecordNotFound) {
		return row, errors.Internal("failed to read the agent", err)
	}
	return row, nil
}

func (s *Service) provision(ctx context.Context, owner schemas.Account, row schemas.Agent, slug string) (schemas.Agent, error) {
	localpart := "agent-" + owner.Localpart + "-" + slug
	if row.DeactivatedAt != nil {
		localpart += "-" + randomSuffix()
	}
	if err := s.deps.Matrix.Register(ctx, localpart); err != nil {
		return row, errors.Internal("failed to register the agent", err)
	}
	row.OwnerID, row.Slug, row.MatrixUserID = owner.UserID, slug, s.deps.Matrix.UserID(localpart)
	row.DeactivatedAt, row.SealedToken, row.DeviceID = nil, "", ""
	return row, nil
}

func (s *Service) login(ctx context.Context, row *schemas.Agent) (Login, error) {
	if token, err := s.deps.Box.Open(row.SealedToken); err == nil && s.deps.Matrix.WhoAmI(ctx, token) {
		return s.save(ctx, row, token)
	}
	login, err := s.deps.Matrix.Login(ctx, row.MatrixUserID, "Bulle")
	if err != nil {
		return Login{}, errors.Internal("failed to sign the agent in", err)
	}
	row.DeviceID, row.SealedToken = login.DeviceID, s.deps.Box.Seal(login.AccessToken)
	return s.save(ctx, row, login.AccessToken)
}

func (s *Service) save(ctx context.Context, row *schemas.Agent, token string) (Login, error) {
	if err := s.deps.DB.WithContext(ctx).Save(row).Error; err != nil {
		return Login{}, errors.Internal("failed to store the agent", err)
	}
	return Login{Homeserver: s.deps.Homeserver, UserID: row.MatrixUserID, AccessToken: token, DeviceID: row.DeviceID}, nil
}

func (s *Service) profile(ctx context.Context, row schemas.Agent, avatarURL string) {
	avatar := ""
	if avatarURL != "" {
		avatar, _ = s.deps.Matrix.UploadFromURL(ctx, row.MatrixUserID, avatarURL)
	}
	_ = s.deps.Matrix.SetProfile(ctx, row.MatrixUserID, row.Name, avatar)
}

func randomSuffix() string {
	random := make([]byte, 3)
	_, _ = rand.Read(random)
	return hex.EncodeToString(random)
}
