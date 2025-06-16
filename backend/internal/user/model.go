package user

import (
	"backend/internal/message"
	"backend/internal/post"
	"backend/internal/subscription"
	"time"
)

type User struct {
	ID           uint                  `gorm:"primaryKey" json:"id"`
	Username     string                `gorm:"uniqueIndex" json:"username"`
	FullName     string                `json:"full_name"`
	Name         string                `gorm:"uniqueIndex" json:"name"`
	FirstName    string                `gorm:"uniqueIndex" json:"first_name"`
	Bio          string                `json:"bio"`
	AvatarURL    string                `json:"avatar_url"`
	Email        string                `gorm:"uniqueIndex" json:"email"`
	PasswordHash string                `json:"-"`
	Role         string                `json:"role"`
	CreatedAt    time.Time             `json:"created_at"`
	Posts         []post.Post                 `gorm:"foreignKey:CreatorID" json:"posts,omitempty"`
	Subscriptions []subscription.Subscription `gorm:"foreignKey:SubscriberID" json:"subscriptions,omitempty"`
	MessagesSent  []message.Message           `gorm:"foreignKey:SenderID" json:"messages_sent,omitempty"`
	MessagesRecv  []message.Message           `gorm:"foreignKey:ReceiverID" json:"messages_recv,omitempty"`
}

func (User) TableName() string {
	return "users"
}

type UpdateUserInput struct {
	FullName  string `json:"full_name"`
	Bio       string `json:"bio"`
	AvatarURL string `json:"avatar_url"`
}
