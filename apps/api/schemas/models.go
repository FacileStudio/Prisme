// Package schemas holds Prisme's tables. Identity lives in porte's tables;
// these map a porte user to Matrix and record what Prisme provisioned.
package schemas

import "time"

// Account is the Matrix user of one porte user.
type Account struct {
	UserID       int64  `gorm:"primaryKey;autoIncrement:false"`
	Localpart    string `gorm:"uniqueIndex;not null"`
	MatrixUserID string `gorm:"uniqueIndex;not null"`
	AvatarURL    string `gorm:"not null;default:''"`
	CreatedAt    time.Time
}

// MatrixSession is the Matrix device handed to one porte browser session, so a
// reload reuses it and a logout revokes it.
type MatrixSession struct {
	PorteSessionID int64 `gorm:"primaryKey;autoIncrement:false"`
	UserID         int64 `gorm:"index;not null"`
	DeviceID       string
	SealedToken    string `gorm:"not null"`
	CreatedAt      time.Time
}

// Space is a Matrix space created through Prisme. RegistreSpaceID links it to
// the Registre space it mirrors; a linked space follows Registre's name and
// members.
type Space struct {
	ID              int64   `gorm:"primaryKey"`
	RoomID          string  `gorm:"uniqueIndex;not null"`
	Name            string  `gorm:"not null"`
	Slug            string  `gorm:"uniqueIndex;not null"`
	CreatedBy       int64   `gorm:"not null"`
	RegistreSpaceID *string `gorm:"uniqueIndex"`
	CreatedAt       time.Time
}

// Channel is a room listed as a child of a space.
type Channel struct {
	RoomID    string `gorm:"primaryKey"`
	SpaceID   int64  `gorm:"index;not null"`
	Name      string `gorm:"not null"`
	Alias     string `gorm:"uniqueIndex;not null"`
	CreatedAt time.Time
}

// Agent is a bot account owned by a user. Slug is its Matrix localpart.
// AvatarURL is the mxc URI it wears and AvatarSource what that image was made
// from: the caller's picture URL, or "blob:" and the name for a drawn blob.
type Agent struct {
	ID            int64  `gorm:"primaryKey"`
	OwnerID       int64  `gorm:"index;not null"`
	Slug          string `gorm:"uniqueIndex;not null"`
	Name          string `gorm:"not null"`
	MatrixUserID  string `gorm:"uniqueIndex;not null"`
	DeviceID      string
	SealedToken   string
	AvatarURL     string `gorm:"not null;default:''"`
	AvatarSource  string `gorm:"not null;default:''"`
	CreatedAt     time.Time
	DeactivatedAt *time.Time
}

// Call is the running call of a Matrix room: the message that announced it
// and who started it. Its row is deleted when the call ends.
type Call struct {
	RoomID    string `gorm:"primaryKey"`
	CallID    string `gorm:"not null"`
	EventID   string `gorm:"not null"`
	StartedBy string `gorm:"not null"`
	StartedAt time.Time
}
