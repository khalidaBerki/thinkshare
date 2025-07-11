package payment

import (
	"backend/internal/db"
	"backend/internal/models"
	"encoding/json"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/webhook"
)

// StripeWebhookHandler gère les notifications Stripe
func StripeWebhookHandler(c *gin.Context) {
	const MaxBodyBytes = int64(65536)
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxBodyBytes)
	payload, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Lecture du corps échouée"})
		return
	}

	sigHeader := c.GetHeader("Stripe-Signature")
	endpointSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")
	event, err := webhook.ConstructEvent(payload, sigHeader, endpointSecret)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Signature Stripe invalide"})
		return
	}

	switch event.Type {
	case "checkout.session.completed":
		var session stripe.CheckoutSession
		if err := json.Unmarshal(event.Data.Raw, &session); err == nil {
			// Récupère les infos de metadata
			creatorID := session.Metadata["creator_id"]
			subscriberID := session.Metadata["subscriber_id"]
			// Active l'abonnement en base
			db.GormDB.Model(&models.Subscription{}).
				Where("creator_id = ? AND subscriber_id = ?", creatorID, subscriberID).
				Update("is_active", true)
			// Stocke l'ID Stripe subscription pour suivi
			if session.Subscription != nil {
				db.GormDB.Model(&models.Subscription{}).
					Where("creator_id = ? AND subscriber_id = ?", creatorID, subscriberID).
					Update("stripe_subscription_id", *session.Subscription)
			}
		}
	case "customer.subscription.deleted", "customer.subscription.updated":
		var sub stripe.Subscription
		if err := json.Unmarshal(event.Data.Raw, &sub); err == nil {
			stripeSubID := sub.ID
			// Recherche l'abonnement local lié à cette subscription Stripe
			var localSub models.Subscription
			if err := db.GormDB.Where("stripe_subscription_id = ?", stripeSubID).First(&localSub).Error; err == nil {
				if sub.Status == "canceled" || sub.Status == "incomplete_expired" || sub.Status == "unpaid" {
					// Désactive l'abonnement local
					db.GormDB.Model(&localSub).Update("is_active", false)
				} else if sub.Status == "active" {
					// Réactive si besoin
					db.GormDB.Model(&localSub).Update("is_active", true)
				}
			}
		}
	}
	c.Status(http.StatusOK)
}
