package payment

import (
	"backend/internal/db"
	"backend/internal/models"
	"encoding/json"
	"io/ioutil"
	"log"
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
	payload, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("[StripeWebhook] Erreur lecture body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Lecture du corps échouée"})
		return
	}

	sigHeader := c.GetHeader("Stripe-Signature")
	endpointSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")
	if endpointSecret == "" {
		log.Printf("[StripeWebhook] STRIPE_WEBHOOK_SECRET manquant")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Secret Stripe webhook manquant"})
		return
	}
	event, err := webhook.ConstructEventWithOptions(
		payload, sigHeader, endpointSecret,
		webhook.ConstructEventOptions{IgnoreAPIVersionMismatch: true},
	)
	if err != nil {
		log.Printf("[StripeWebhook] Signature Stripe invalide: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Signature Stripe invalide"})
		return
	}

	log.Printf("[StripeWebhook] Event reçu: %s", event.Type)

	switch event.Type {
	case "checkout.session.completed":
		var session stripe.CheckoutSession
		if err := json.Unmarshal(event.Data.Raw, &session); err == nil {
			creatorID := session.Metadata["creator_id"]
			subscriberID := session.Metadata["subscriber_id"]
			log.Printf("[StripeWebhook] checkout.session.completed: creator_id=%s, subscriber_id=%s", creatorID, subscriberID)
			db.GormDB.Model(&models.Subscription{}).
				Where("creator_id = ? AND subscriber_id = ?", creatorID, subscriberID).
				Update("is_active", true)
			if session.Subscription != nil {
				db.GormDB.Model(&models.Subscription{}).
					Where("creator_id = ? AND subscriber_id = ?", creatorID, subscriberID).
					Update("stripe_subscription_id", *session.Subscription)
			}
		} else {
			log.Printf("[StripeWebhook] Erreur parsing session: %v", err)
		}
	case "customer.subscription.deleted", "customer.subscription.updated":
		var sub stripe.Subscription
		if err := json.Unmarshal(event.Data.Raw, &sub); err == nil {
			stripeSubID := sub.ID
			var localSub models.Subscription
			if err := db.GormDB.Where("stripe_subscription_id = ?", stripeSubID).First(&localSub).Error; err == nil {
				if sub.Status == "canceled" || sub.Status == "incomplete_expired" || sub.Status == "unpaid" {
					db.GormDB.Model(&localSub).Update("is_active", false)
					log.Printf("[StripeWebhook] Abonnement désactivé: %s", stripeSubID)
				} else if sub.Status == "active" {
					db.GormDB.Model(&localSub).Update("is_active", true)
					log.Printf("[StripeWebhook] Abonnement réactivé: %s", stripeSubID)
				}
			}
		} else {
			log.Printf("[StripeWebhook] Erreur parsing subscription: %v", err)
		}
	}
	c.Status(http.StatusOK)
}
