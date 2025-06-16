package message

import "time"

type Message struct {
	ID         uint `gorm:"primaryKey"`
	SenderID   uint
	ReceiverID uint
	Content    string
	Timestamp  time.Time
}