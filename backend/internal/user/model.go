package user

import (
	"backend/internal/message"
	"backend/internal/post"
	"backend/internal/subscription"
	"time"
)

// User représente le modèle complet d'un utilisateur (en base de données)
type User struct {
	ID            uint                        `gorm:"primaryKey" json:"id" example:"1"`
	Username      string                      `gorm:"uniqueIndex" json:"username" example:"haithemdev"`
	FullName      string                      `json:"full_name" example:"Haithem Hammami"`
	Name          string                      `gorm:"uniqueIndex" json:"name" example:"Hammami"`
	FirstName     string                      `gorm:"uniqueIndex" json:"first_name" example:"Haithem"`
	Bio           string                      `json:"bio" example:"Étudiant à l’EEMI et dev fullstack"`
	AvatarURL     string                      `json:"avatar_url" example:"https://cdn.thinkshare/avatar.jpg"`
	Email         string                      `gorm:"uniqueIndex" json:"email" example:"haithem@example.com"`
	PasswordHash  string                      `json:"-"`
	Role          string                      `json:"role" example:"user"`
	CreatedAt     time.Time                   `json:"created_at" example:"2024-01-01T15:04:05Z"`
	Posts         []post.Post                 `gorm:"foreignKey:CreatorID" json:"posts,omitempty"`
	Subscriptions []subscription.Subscription `gorm:"foreignKey:SubscriberID" json:"subscriptions,omitempty"`
	MessagesSent  []message.Message           `gorm:"foreignKey:SenderID" json:"messages_sent,omitempty"`
	MessagesRecv  []message.Message           `gorm:"foreignKey:ReceiverID" json:"messages_recv,omitempty"`
}

// ProfileDTO est une version simplifiée de User, envoyée au client (sans email, password, etc.)
type ProfileDTO struct {
	ID        uint   `json:"id" example:"1"`
	FullName  string `json:"full_name" example:"Haithem Hammami"`
	Bio       string `json:"bio" example:"Étudiant à l’EEMI et dev fullstack"`
	AvatarURL string `json:"avatar_url" example:"https://cdn.thinkshare/avatar.jpg"`
}

// TableName permet de forcer le nom de la table "users"
func (User) TableName() string {
	return "users"
}

// UpdateUserInput représente les données reçues lors d'une modification du profil
type UpdateUserInput struct {
	FullName  string `json:"full_name" example:"Haithem Hammami"`
	Bio       string `json:"bio" example:"Développeur Go, passionné par l'éducation"`
	AvatarURL string `json:"avatar_url" example:"https://cdn.thinkshare/avatar.jpg"`
}
