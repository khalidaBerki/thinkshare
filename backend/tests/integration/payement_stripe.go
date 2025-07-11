package integration

import (
	"backend/internal/db"
	"backend/internal/models"
	"testing"
	"time"
)

func TestStripeSubscriptionWebhookCreatesOrActivatesSubscription(t *testing.T) {
	subscriberID := uint(1001) // à adapter selon ta base de test
	creatorID := uint(7)       // créateur premium

	// Nettoyage avant test
	db.GormDB.Where("subscriber_id = ? AND creator_id = ?", subscriberID, creatorID).Delete(&models.Subscription{})

	// Simule une création de subscription Stripe (comme le webhook)
	sub := models.Subscription{
		SubscriberID: subscriberID,
		CreatorID:    creatorID,
		IsActive:     true,
		Type:         "stripe",
		StartDate:    time.Now(),
		EndDate:      time.Now().AddDate(0, 1, 0),
	}
	if err := db.GormDB.Create(&sub).Error; err != nil {
		t.Fatalf("Erreur création subscription Stripe: %v", err)
	}

	// Vérifie que la subscription existe et est active
	var found models.Subscription
	err := db.GormDB.Where("subscriber_id = ? AND creator_id = ? AND is_active = ?", subscriberID, creatorID, true).First(&found).Error
	if err != nil {
		t.Fatalf("Subscription Stripe non trouvée ou inactive: %v", err)
	}
	if found.Type != "stripe" {
		t.Errorf("Type attendu 'stripe', obtenu: %s", found.Type)
	}
}
