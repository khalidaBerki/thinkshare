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

// Simule un appel Stripe webhook avec metadata correcte
func TestStripeWebhookHandler_CreatesSubscription(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/webhook", payment.StripeWebhookHandler)

	subscriberID := uint(1002) // à adapter
	creatorID := uint(7)

	// Nettoyage
	db.GormDB.Where("subscriber_id = ? AND creator_id = ?", subscriberID, creatorID).Delete(&models.Subscription{})

	// Simule un event Stripe checkout.session.completed
	session := map[string]interface{}{
		"metadata": map[string]string{
			"creator_id":    "7",
			"subscriber_id": "1002",
		},
		"id": "cs_test_123",
	}
	event := map[string]interface{}{
		"type": "checkout.session.completed",
		"data": map[string]interface{}{
			"object": session,
		},
	}
	payload, _ := json.Marshal(event)

	// Appel HTTP
	req := httptest.NewRequest("POST", "/webhook", bytes.NewBuffer(payload))
	req.Header.Set("Stripe-Signature", "test") // ignoré car on ne vérifie pas la signature ici
	w := httptest.NewRecorder()

	os.Setenv("STRIPE_WEBHOOK_SECRET", "test")
	os.Setenv("DISABLE_STRIPE_SIGNATURE_CHECK", "true")
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("Webhook Stripe HTTP code attendu 200, obtenu %d", w.Code)
	}

	// Vérifie la subscription
	var found models.Subscription
	err := db.GormDB.Where("subscriber_id = ? AND creator_id = ? AND is_active = ?", subscriberID, creatorID, true).First(&found).Error
	if err != nil {
		t.Fatalf("Subscription Stripe non trouvée ou inactive après webhook: %v", err)
	}
}
