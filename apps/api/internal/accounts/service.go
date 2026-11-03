// Package accounts gives every porte user exactly one Matrix account, created
// on first use and named after their email.
package accounts

import (
	"context"
	stderrors "errors"
	"log/slog"

	"github.com/FacileStudio/Prisme/apps/api/internal/matrix"
	"github.com/FacileStudio/Prisme/apps/api/schemas"
	"github.com/FacileStudio/tronc/errors"

	"gorm.io/gorm"
)

const maxCandidates = 50

// Profile is what porte knows about a user.
type Profile struct {
	Email     string
	Name      string
	AvatarURL string
}

// Service provisions and looks up Matrix accounts.
type Service struct {
	db     *gorm.DB
	matrix *matrix.Client
	logger *slog.Logger
}

// NewService builds the account service.
func NewService(db *gorm.DB, client *matrix.Client, logger *slog.Logger) *Service {
	return &Service{db: db, matrix: client, logger: logger}
}

// Ensure returns the user's Matrix account, registering it on first call.
func (s *Service) Ensure(ctx context.Context, userID int64) (schemas.Account, error) {
	var account schemas.Account
	err := s.db.WithContext(ctx).First(&account, "user_id = ?", userID).Error
	if err == nil {
		return account, nil
	}
	if !stderrors.Is(err, gorm.ErrRecordNotFound) {
		return account, errors.Internal("failed to read the account", err)
	}
	profile, err := s.Profile(ctx, userID)
	if err != nil {
		return account, err
	}
	account, err = s.create(ctx, userID, Localpart(profile.Email))
	if err != nil {
		return account, err
	}
	s.SyncProfile(ctx, account, profile)
	return account, nil
}

// Profile reads the user's name, email and picture from porte's table.
func (s *Service) Profile(ctx context.Context, userID int64) (Profile, error) {
	var profile Profile
	err := s.db.WithContext(ctx).Table("porte_users").Select("email, name, avatar_url").
		Where("id = ?", userID).Take(&profile).Error
	if err != nil {
		return profile, errors.Internal("failed to read the user", err)
	}
	return profile, nil
}

// SyncProfile pushes the porte name and picture onto the Matrix profile. A
// failure is logged, not returned: a missing avatar must not block sign-in.
func (s *Service) SyncProfile(ctx context.Context, account schemas.Account, profile Profile) {
	name := profile.Name
	if name == "" {
		name = account.Localpart
	}
	avatar := ""
	if profile.AvatarURL != "" {
		uploaded, err := s.matrix.UploadFromURL(ctx, account.MatrixUserID, profile.AvatarURL)
		if err != nil {
			s.logger.Warn("avatar upload failed", slog.String("user", account.MatrixUserID), slog.Any("error", err))
		}
		avatar = uploaded
	}
	if err := s.matrix.SetProfile(ctx, account.MatrixUserID, name, avatar); err != nil {
		s.logger.Warn("profile sync failed", slog.String("user", account.MatrixUserID), slog.Any("error", err))
	}
}

func (s *Service) create(ctx context.Context, userID int64, base string) (schemas.Account, error) {
	for n := 1; n <= maxCandidates; n++ {
		localpart := Candidate(base, n)
		account := schemas.Account{UserID: userID, Localpart: localpart, MatrixUserID: s.matrix.UserID(localpart)}
		err := s.db.WithContext(ctx).Create(&account).Error
		if stderrors.Is(err, gorm.ErrDuplicatedKey) {
			if existing, found := s.existing(ctx, userID); found {
				return existing, nil
			}
			continue
		}
		if err != nil {
			return account, errors.Internal("failed to store the account", err)
		}
		if err := s.matrix.Register(ctx, localpart); err != nil {
			s.db.WithContext(ctx).Delete(&account)
			return account, errors.Internal("failed to register the Matrix account", err)
		}
		return account, nil
	}
	return schemas.Account{}, errors.Conflict("no free Matrix name for this user")
}

func (s *Service) existing(ctx context.Context, userID int64) (schemas.Account, bool) {
	var account schemas.Account
	err := s.db.WithContext(ctx).First(&account, "user_id = ?", userID).Error
	return account, err == nil
}
