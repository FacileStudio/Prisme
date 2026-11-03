package agents

import (
	"context"
	stderrors "errors"

	"github.com/FacileStudio/Prisme/apps/api/schemas"
	"github.com/FacileStudio/tronc/errors"

	"gorm.io/gorm"
)

func (s *Service) resolveRow(ctx context.Context, owner schemas.Account, slug string) (schemas.Agent, error) {
	row, err := s.find(ctx, slug)
	if err != nil || row.ID != 0 {
		return row, err
	}
	return s.provision(ctx, owner, slug)
}

func (s *Service) find(ctx context.Context, slug string) (schemas.Agent, error) {
	var row schemas.Agent
	err := s.deps.DB.WithContext(ctx).First(&row, "slug = ?", slug).Error
	if err != nil && !stderrors.Is(err, gorm.ErrRecordNotFound) {
		return row, errors.Internal("failed to read the agent", err)
	}
	return row, nil
}

func (s *Service) provision(ctx context.Context, owner schemas.Account, slug string) (schemas.Agent, error) {
	localpart := slug
	if err := s.deps.Matrix.Register(ctx, localpart); err != nil {
		return schemas.Agent{}, errors.Internal("failed to register the agent", err)
	}
	row := schemas.Agent{
		OwnerID:      owner.UserID,
		Slug:         slug,
		MatrixUserID: s.deps.Matrix.UserID(localpart),
	}
	return row, nil
}

func (s *Service) login(ctx context.Context, row *schemas.Agent) (Login, error) {
	token := ""
	if openToken, err := s.deps.Box.Open(row.SealedToken); err == nil && s.deps.Matrix.WhoAmI(ctx, openToken) {
		token = openToken
	} else {
		login, err := s.deps.Matrix.Login(ctx, row.MatrixUserID, "Bulle")
		if err != nil {
			return Login{}, errors.Internal("failed to sign the agent in", err)
		}
		row.DeviceID, row.SealedToken, token = login.DeviceID, s.deps.Box.Seal(login.AccessToken), login.AccessToken
	}
	if err := s.deps.DB.WithContext(ctx).Save(row).Error; err != nil {
		return Login{}, errors.Internal("failed to store the agent", err)
	}
	return Login{Homeserver: s.deps.Homeserver, UserID: row.MatrixUserID, AccessToken: token, DeviceID: row.DeviceID}, nil
}
