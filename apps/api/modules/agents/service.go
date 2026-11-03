package agents

import (
	"context"
	"fmt"
	"strings"

	"github.com/FacileStudio/Prisme/apps/api/internal/accounts"
	"github.com/FacileStudio/Prisme/apps/api/internal/matrix"
	"github.com/FacileStudio/Prisme/apps/api/internal/seal"
	"github.com/FacileStudio/Prisme/apps/api/modules/spaces"
	"github.com/FacileStudio/Prisme/apps/api/schemas"
	"github.com/FacileStudio/tronc/errors"

	"gorm.io/gorm"
)

// Deps holds the services agents rely on.
type Deps struct {
	DB         *gorm.DB
	Matrix     *matrix.Client
	Accounts   *accounts.Service
	Box        *seal.Box
	Homeserver string
}

// Service manages agents.
type Service struct {
	deps Deps
}

// NewService creates a new agent service.
func NewService(deps Deps) *Service { return &Service{deps: deps} }

func resolveSlug(name string) (string, string) {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" || len(trimmed) > 80 {
		return "", "an agent needs a name of 1 to 80 characters"
	}
	slug := spaces.Slugify(trimmed, "agent")
	if slug == "" || slug == "prisme" || slug == "admin" || slug == "matrix" {
		return "", fmt.Sprintf("the name %q is reserved", name)
	}
	return slug, ""
}

func (s *Service) checkTaken(ctx context.Context, slug string, ownerID int64) string {
	var count int64
	if s.deps.DB.WithContext(ctx).Model(&schemas.Account{}).Where("localpart = ?", slug).Count(&count).Error != nil {
		return "failed to check account availability"
	}
	if count > 0 {
		return fmt.Sprintf("the name %q is already taken", slug)
	}
	var existing schemas.Agent
	err := s.deps.DB.WithContext(ctx).First(&existing, "slug = ?", slug).Error
	if err == nil {
		if existing.OwnerID != ownerID {
			return fmt.Sprintf("the name %q is already taken", slug)
		}
		if existing.DeactivatedAt != nil {
			return fmt.Sprintf("the name %q was deactivated and cannot be reused", slug)
		}
	}
	return ""
}

func (s *Service) CheckAvailable(ctx context.Context, ownerID int64, name string) (AvailabilityResponse, error) {
	slug, reason := resolveSlug(name)
	if reason != "" {
		return AvailabilityResponse{Name: name, Available: false, Reason: reason}, nil
	}
	if reason := s.checkTaken(ctx, slug, ownerID); reason != "" {
		return AvailabilityResponse{Name: name, Slug: slug, Available: false, Reason: reason}, nil
	}
	return AvailabilityResponse{Name: name, Slug: slug, Available: true}, nil
}

func (s *Service) Create(ctx context.Context, ownerID int64, req CreateRequest) (Login, error) {
	slug, reason := resolveSlug(req.Name)
	if reason != "" {
		return Login{}, errors.Invalid(reason)
	}
	if reason := s.checkTaken(ctx, slug, ownerID); reason != "" {
		return Login{}, errors.Conflict(reason)
	}
	owner, err := s.deps.Accounts.Ensure(ctx, ownerID)
	if err != nil {
		return Login{}, err
	}
	row, err := s.resolveRow(ctx, owner, slug)
	if err != nil {
		return Login{}, err
	}
	row.Name = strings.TrimSpace(req.Name)
	login, err := s.login(ctx, &row)
	if err != nil {
		return Login{}, err
	}
	s.profile(ctx, row, req.AvatarURL)
	return login, nil
}

func (s *Service) profile(ctx context.Context, row schemas.Agent, avatarURL string) {
	avatar := ""
	if avatarURL != "" {
		avatar, _ = s.deps.Matrix.UploadFromURL(ctx, row.MatrixUserID, avatarURL)
	}
	_ = s.deps.Matrix.SetProfile(ctx, row.MatrixUserID, row.Name, avatar)
}
