package subscription

import "time"

type Subscription struct {
	ID           uint `gorm:"primaryKey"`
	SubscriberID uint
	CreatorID    uint
	StartDate    time.Time
	EndDate      time.Time
	IsActive     bool
	Type         string // monthly or one-time
}
