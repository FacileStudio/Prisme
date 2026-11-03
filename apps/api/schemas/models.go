// Package schemas holds Prisme's tables. Identity lives in porte's tables;
// these map a porte user to Matrix and record what Prisme provisioned.
package schemas

import "time"

// Account is the Matrix user of one porte user.
type Account struct {
	UserID       int64  `gorm:"primaryKey;autoIncrement:false"`
	Localpart    string `gorm:"uniqueIndex;not null"`
	MatrixUserID string `gorm:"uniqueIndex;not null"`
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

// Space is a Matrix space created through Prisme.
type Space struct {
	ID        int64  `gorm:"primaryKey"`
	RoomID    string `gorm:"uniqueIndex;not null"`
	Name      string `gorm:"not null"`
	Slug      string `gorm:"uniqueIndex;not null"`
	CreatedBy int64  `gorm:"not null"`
	CreatedAt time.Time
}

// Channel is a room listed as a child of a space.
type Channel struct {
	RoomID    string `gorm:"primaryKey"`
	SpaceID   int64  `gorm:"index;not null"`
	Name      string `gorm:"not null"`
	Alias     string `gorm:"uniqueIndex;not null"`
	CreatedAt time.Time
}

// Agent is a bot account owned by a user, one per (owner, slug).
type Agent struct {
	ID            int64  `gorm:"primaryKey"`
	OwnerID       int64  `gorm:"uniqueIndex:agents_owner_slug;not null"`
	Slug          string `gorm:"uniqueIndex:agents_owner_slug;not null"`
	Name          string `gorm:"not null"`
	MatrixUserID  string `gorm:"uniqueIndex;not null"`
	DeviceID      string
	SealedToken   string
	CreatedAt     time.Time
	DeactivatedAt *time.Time
}
