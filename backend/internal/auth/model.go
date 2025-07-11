package auth

import (
	"time"
)

type AuthToken struct {
	ID        uint   `gorm:"primaryKey"`
	UserID    uint   // Clé étrangère vers User
	Token     string `gorm:"uniqueIndex"`
	ExpiresAt time.Time
	CreatedAt time.Time
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type TokenResponse struct {
	Token  string `json:"token"`
	UserID uint   `json:"user_id"`
}
