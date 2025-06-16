package auth

import (
	"time"
)

type AuthToken struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      // Clé étrangère vers User
	Token     string    `gorm:"uniqueIndex"`
	ExpiresAt time.Time
	CreatedAt time.Time
}
