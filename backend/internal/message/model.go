package message

import (
	"time"
)

// MessageStatus représente l'état d'un message
type MessageStatus string

const (
	StatusUnread   MessageStatus = "UNREAD"
	StatusRead     MessageStatus = "READ"
	StatusArchived MessageStatus = "ARCHIVED"
	StatusDeleted  MessageStatus = "DELETED"
)

// Message représente un message privé entre deux utilisateurs
type Message struct {
	ID         uint          `gorm:"primaryKey"`
	SenderID   uint          `gorm:"not null"`
	ReceiverID uint          `gorm:"not null"`
	Content    string        `gorm:"type:text;not null"`
	Status     MessageStatus `gorm:"default:'UNREAD'"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  *time.Time `gorm:"index"`
}

// DTO pour la création d’un message (reçu via JSON)
type CreateMessageInput struct {
	ReceiverID uint   `json:"receiver_id" binding:"required"`
	Content    string `json:"content" binding:"required"`
}

// DTO pour l'affichage enrichi d’un message
type MessageDTO struct {
	ID        uint          `json:"id"`
	Content   string        `json:"content"`
	Status    MessageStatus `json:"status"`
	CreatedAt time.Time     `json:"created_at"`

	// Infos utilisateur enrichies
	Sender   *UserInfo `json:"sender"`
	Receiver *UserInfo `json:"receiver"`
}

// Représente les infos de base d’un utilisateur (utilisé dans les DTOs)
type UserInfo struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	AvatarURL string `json:"avatar_url"`
}

// MessagePreviewDTO = structure pour la liste des conversations
type MessagePreviewDTO struct {
	ConversationID string    `json:"conversation_id"` // clé virtuelle: "user1-user2"
	LastMessage    string    `json:"last_message"`
	Timestamp      time.Time `json:"timestamp"`
	UnreadCount    int       `json:"unread_count"`
	OtherUser      *UserInfo `json:"other_user"` // celui avec qui je parle
}

type MessagePreviewRaw struct {
	LastMessage    string    `json:"last_message"`
	OtherUserID    uint      `json:"other_user_id"`
	OtherUsername  string    `json:"other_username"`
	OtherAvatarURL string    `json:"other_avatar_url"`
	CreatedAt      time.Time `json:"created_at"`
	UnreadCount    int       `json:"unread_count"`
}
