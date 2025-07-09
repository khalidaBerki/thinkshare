package payment

import "time"

type Payment struct {
	ID             uint `gorm:"primaryKey"`
	UserID         uint
	Amount         float64
	Type           string // subscription or post
	Status         string
	Date           time.Time
	SubscriptionID *uint
	PostID         *uint
}
