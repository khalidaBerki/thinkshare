package payment

import (
	"backend/internal/db"
	"backend/internal/models"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

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

	var eventType string
	var eventData json.RawMessage
	var event stripe.Event

	if os.Getenv("DISABLE_STRIPE_SIGNATURE_CHECK") == "true" {
		log.Printf("[StripeWebhook][TEST] Vérification de signature Stripe désactivée pour les tests")
		var testEvent struct {
			Type string `json:"type"`
			Data struct {
				Object json.RawMessage `json:"object"`
			} `json:"data"`
		}
		if err := json.Unmarshal(payload, &testEvent); err != nil {
			log.Printf("[StripeWebhook][TEST] Erreur parsing event test: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Event test Stripe invalide"})
			return
		}
		eventType = testEvent.Type
		eventData = testEvent.Data.Object
	} else {
		if endpointSecret == "" {
			log.Printf("[StripeWebhook] STRIPE_WEBHOOK_SECRET manquant")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Secret Stripe webhook manquant"})
			return
		}
		event, err = webhook.ConstructEventWithOptions(
			payload, sigHeader, endpointSecret,
			webhook.ConstructEventOptions{IgnoreAPIVersionMismatch: true},
		)
		if err != nil {
			log.Printf("[StripeWebhook] Signature Stripe invalide: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Signature Stripe invalide"})
			return
		}
		eventType = string(event.Type)
		eventData = event.Data.Raw
	}

	log.Printf("[StripeWebhook] Event reçu: %s", eventType)

	switch eventType {
	case "checkout.session.completed":
		var session stripe.CheckoutSession
		if err := json.Unmarshal(eventData, &session); err == nil {
			creatorID := session.Metadata["creator_id"]
			subscriberID := session.Metadata["subscriber_id"]
			log.Printf("[StripeWebhook] checkout.session.completed: creator_id=%s, subscriber_id=%s, session_id=%s", creatorID, subscriberID, session.ID)

			// Vérifier si la subscription existe déjà
			var sub models.Subscription
			err := db.GormDB.Where("creator_id = ? AND subscriber_id = ?", creatorID, subscriberID).First(&sub).Error
			if err != nil {
				// Si elle n'existe pas, on la crée
				sub = models.Subscription{
					CreatorID:    parseUintOrZero(creatorID),
					SubscriberID: parseUintOrZero(subscriberID),
					IsActive:     true,
					Type:         "stripe",
				}
				if session.Subscription != nil {
					sub.StripeSubscriptionID = session.Subscription.ID
				}
				if err := db.GormDB.Create(&sub).Error; err != nil {
					log.Printf("[StripeWebhook][ERROR] Erreur création subscription DB: %v", err)
				} else {
					log.Printf("[StripeWebhook] Subscription créée: creator_id=%d, subscriber_id=%d, is_active=%v", sub.CreatorID, sub.SubscriberID, sub.IsActive)
				}
			} else {
				// Sinon, on l'active
				if err := db.GormDB.Model(&sub).Update("is_active", true).Error; err != nil {
					log.Printf("[StripeWebhook][ERROR] Erreur activation subscription DB: %v", err)
				}
				if session.Subscription != nil {
					db.GormDB.Model(&sub).Update("stripe_subscription_id", session.Subscription.ID)
				}
				log.Printf("[StripeWebhook] Subscription activée: creator_id=%d, subscriber_id=%d, is_active=%v", sub.CreatorID, sub.SubscriberID, true)
			}
		} else {
			log.Printf("[StripeWebhook] Erreur parsing session: %v", err)
		}
	case "customer.subscription.deleted", "customer.subscription.updated":
		var sub stripe.Subscription
		if err := json.Unmarshal(eventData, &sub); err == nil {
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

// Utilitaire pour parser un uint à partir d'une string
func parseUintOrZero(s string) uint {
	u, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0
	}
	return uint(u)
}
