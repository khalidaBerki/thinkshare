package integration

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"os"
	"testing"

	"backend/internal/db"
	"backend/internal/models"
	"backend/internal/payment"

	"github.com/gin-gonic/gin"
)

// Test complet du flux Stripe côté backend (simulateur end-to-end)
func TestStripeEndToEndBackendFlow(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/webhook", payment.StripeWebhookHandler)

	subscriberID := uint(1003) // à adapter
	creatorID := uint(7)

	// Nettoyage
	db.GormDB.Where("subscriber_id = ? AND creator_id = ?", subscriberID, creatorID).Delete(&models.Subscription{})

	// 1. Simule le paiement Stripe (webhook checkout.session.completed)
	session := map[string]interface{}{
		"metadata": map[string]string{
			"creator_id":    "7",
			"subscriber_id": "1003",
		},
		"id": "cs_test_456",
	}
	event := map[string]interface{}{
		"type": "checkout.session.completed",
		"data": map[string]interface{}{
			"object": session,
		},
	}
	payload, _ := json.Marshal(event)

	os.Setenv("STRIPE_WEBHOOK_SECRET", "test")
	os.Setenv("DISABLE_STRIPE_SIGNATURE_CHECK", "true")

	req := httptest.NewRequest("POST", "/webhook", bytes.NewBuffer(payload))
	req.Header.Set("Stripe-Signature", "test")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("Webhook Stripe HTTP code attendu 200, obtenu %d", w.Code)
	}

	// 2. Vérifie la subscription créée/active
	var found models.Subscription
	err := db.GormDB.Where("subscriber_id = ? AND creator_id = ? AND is_active = ?", subscriberID, creatorID, true).First(&found).Error
	if err != nil {
		t.Fatalf("Subscription Stripe non trouvée ou inactive après webhook: %v", err)
	}
	if found.Type != "stripe" {
		t.Errorf("Type attendu 'stripe', obtenu: %s", found.Type)
	}

	// 3. Simule expiration (désactivation) via webhook Stripe
	subscriptionEvent := map[string]interface{}{
		"type": "customer.subscription.deleted",
		"data": map[string]interface{}{
			"object": map[string]interface{}{
				"id":     "stripe_sub_id_test",
				"status": "canceled",
			},
		},
	}
	found.StripeSubscriptionID = "stripe_sub_id_test"
	db.GormDB.Save(&found)

	payload2, _ := json.Marshal(subscriptionEvent)
	req2 := httptest.NewRequest("POST", "/webhook", bytes.NewBuffer(payload2))
	req2.Header.Set("Stripe-Signature", "test")
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)

	if w2.Code != 200 {
		t.Fatalf("Webhook Stripe HTTP code attendu 200 pour désactivation, obtenu %d", w2.Code)
	}

	var found2 models.Subscription
	err = db.GormDB.Where("subscriber_id = ? AND creator_id = ?", subscriberID, creatorID).First(&found2).Error
	if err != nil {
		t.Fatalf("Subscription Stripe non trouvée après désactivation: %v", err)
	}
	if found2.IsActive {
		t.Errorf("Subscription devrait être inactive après désactivation Stripe")
	}
}
