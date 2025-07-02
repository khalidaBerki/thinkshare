package models

import "time"

type Subscription struct {
	ID                   uint `gorm:"primaryKey"`
	SubscriberID         uint
	CreatorID            uint
	StartDate            time.Time
	EndDate              time.Time
	IsActive             bool
	Type                 string
	StripeSubscriptionID string // ID Stripe de la subscription pour suivi
}
