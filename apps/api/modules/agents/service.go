package agents

import (
	"context"
	"fmt"
	"log/slog"
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
	Logger     *slog.Logger
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

// CheckAvailable reports whether the owner can create an agent with this name,
// and the localpart Create would give it.
func (s *Service) CheckAvailable(ctx context.Context, ownerID int64, name string) (AvailabilityResponse, error) {
	slug, reason := resolveSlug(name)
	if reason != "" {
		return AvailabilityResponse{Name: name, Available: false, Reason: reason}, nil
	}
	who, err := s.identify(ctx, ownerID, name, slug)
	if err != nil {
		return AvailabilityResponse{}, err
	}
	if who.localpart == "" {
		return AvailabilityResponse{Name: name, Slug: slug, Available: false, Reason: noFreeName}, nil
	}
	return AvailabilityResponse{Name: name, Slug: who.localpart, Available: true}, nil
}

// Create returns the owner's agent with this name, provisioning it on first
// call and moving it onto its wanted Matrix id when it sits elsewhere.
func (s *Service) Create(ctx context.Context, ownerID int64, req CreateRequest) (Login, error) {
	slug, reason := resolveSlug(req.Name)
	if reason != "" {
		return Login{}, errors.Invalid(reason)
	}
	who, err := s.identify(ctx, ownerID, req.Name, slug)
	if err != nil {
		return Login{}, err
	}
	if who.localpart == "" {
		return Login{}, errors.Conflict(noFreeName)
	}
	row := who.row
	if err := s.settle(ctx, who.owner, &row, who.localpart); err != nil {
		return Login{}, err
	}
	row.Name = strings.TrimSpace(req.Name)
	login, err := s.login(ctx, &row)
	if err != nil {
		return Login{}, err
	}
	s.profile(ctx, &row, req.AvatarURL)
	login.OwnerUserID = who.owner.MatrixUserID
	return login, nil
}

func (s *Service) settle(ctx context.Context, owner schemas.Account, row *schemas.Agent, localpart string) error {
	if row.ID == 0 {
		return s.provision(ctx, owner, row, localpart)
	}
	if row.MatrixUserID != s.deps.Matrix.UserID(localpart) {
		return s.migrate(ctx, owner, row, localpart)
	}
	return nil
}
