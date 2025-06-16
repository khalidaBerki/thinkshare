package user

import (
	"backend/internal/message"
	"backend/internal/post"
	"backend/internal/subscription"
	"time"
)

type User struct {
	ID           uint   `gorm:"primaryKey"`
	Name         string `gorm:"uniqueIndex"`
	FirstName    string `gorm:"uniqueIndex"`
	Username     string `gorm:"uniqueIndex"`
	Email        string `gorm:"uniqueIndex"`
	PasswordHash string
	Role         string
	CreatedAt    time.Time

	Posts         []post.Post                 `gorm:"foreignKey:CreatorID"`
	Subscriptions []subscription.Subscription `gorm:"foreignKey:SubscriberID"`
	MessagesSent  []message.Message           `gorm:"foreignKey:SenderID"`
	MessagesRecv  []message.Message           `gorm:"foreignKey:ReceiverID"`
}

func (User) TableName() string {
	return "users"
}
