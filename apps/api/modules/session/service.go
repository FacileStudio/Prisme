// Package session hands a signed-in browser its Matrix login: one device per
// porte session, reused across reloads and revoked on logout.
package session

import (
	"context"
	stderrors "errors"

	"github.com/FacileStudio/Prisme/apps/api/internal/accounts"
	"github.com/FacileStudio/Prisme/apps/api/internal/matrix"
	"github.com/FacileStudio/Prisme/apps/api/internal/seal"
	"github.com/FacileStudio/Prisme/apps/api/schemas"
	"github.com/FacileStudio/tronc/errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Deps is what the session service is built from.
type Deps struct {
	DB           *gorm.DB
	Matrix       *matrix.Client
	Accounts     *accounts.Service
	Box          *seal.Box
	Homeserver   string
	CallsEnabled bool
}

// Service mints and revokes browser Matrix logins.
type Service struct {
	deps Deps
}

// NewService builds the session service.
func NewService(deps Deps) *Service { return &Service{deps: deps} }

// Get returns the user and a working Matrix login for this porte session.
func (s *Service) Get(ctx context.Context, userID, porteSessionID int64) (Response, error) {
	account, err := s.deps.Accounts.Ensure(ctx, userID)
	if err != nil {
		return Response{}, err
	}
	profile, err := s.deps.Accounts.Profile(ctx, userID)
	if err != nil {
		return Response{}, err
	}
	login, err := s.login(ctx, account, profile, porteSessionID)
	if err != nil {
		return Response{}, err
	}
	return Response{
		User: User{
			ID: userID, Email: profile.Email, Name: profile.Name,
			AvatarURL: profile.AvatarURL, MatrixUserID: account.MatrixUserID,
		},
		Matrix:       MatrixLogin{Homeserver: s.deps.Homeserver, UserID: login.UserID, AccessToken: login.AccessToken, DeviceID: login.DeviceID},
		CallsEnabled: s.deps.CallsEnabled,
	}, nil
}

func (s *Service) login(ctx context.Context, account schemas.Account, profile accounts.Profile, porteSessionID int64) (matrix.Login, error) {
	if login, ok := s.stored(ctx, account, porteSessionID); ok {
		return login, nil
	}
	login, err := s.deps.Matrix.Login(ctx, account.MatrixUserID, "Prisme web")
	if err != nil {
		return login, errors.Internal("failed to sign in to Matrix", err)
	}
	row := schemas.MatrixSession{
		PorteSessionID: porteSessionID, UserID: account.UserID,
		DeviceID: login.DeviceID, SealedToken: s.deps.Box.Seal(login.AccessToken),
	}
	err = s.deps.DB.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Create(&row).Error
	if err != nil {
		return login, errors.Internal("failed to store the Matrix login", err)
	}
	s.deps.Accounts.SyncProfile(ctx, account, profile)
	return login, nil
}

func (s *Service) stored(ctx context.Context, account schemas.Account, porteSessionID int64) (matrix.Login, bool) {
	var row schemas.MatrixSession
	err := s.deps.DB.WithContext(ctx).First(&row, "porte_session_id = ? AND user_id = ?", porteSessionID, account.UserID).Error
	if err != nil {
		return matrix.Login{}, false
	}
	token, err := s.deps.Box.Open(row.SealedToken)
	if err != nil || !s.deps.Matrix.WhoAmI(ctx, token) {
		return matrix.Login{}, false
	}
	return matrix.Login{UserID: account.MatrixUserID, AccessToken: token, DeviceID: row.DeviceID}, true
}

// Revoke logs the Matrix device of a porte session out and forgets it.
func (s *Service) Revoke(ctx context.Context, porteSessionID int64) error {
	var row schemas.MatrixSession
	err := s.deps.DB.WithContext(ctx).First(&row, "porte_session_id = ?", porteSessionID).Error
	if stderrors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}
	if err != nil {
		return err
	}
	if token, openErr := s.deps.Box.Open(row.SealedToken); openErr == nil {
		if err := s.deps.Matrix.Logout(ctx, token); err != nil {
			return err
		}
	}
	return s.deps.DB.WithContext(ctx).Delete(&row).Error
}
