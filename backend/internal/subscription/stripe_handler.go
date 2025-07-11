package subscription

import (
	"backend/internal/payment"
	"backend/internal/user"
	"net/http"
	"os"
	"strconv"

	"log"

	"github.com/gin-gonic/gin"
)

// SubscribePaidStripeHandler : Crée une session Stripe pour l'abonnement mensuel
func SubscribePaidStripeHandler(c *gin.Context) {
	var input SubscriptionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Entrée invalide"})
		return
	}

	subscriberID := c.GetInt("user_id")
	if uint(subscriberID) == input.CreatorID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Vous ne pouvez pas vous abonner à vous-même"})
		return
	}

	creator, err := user.GetUserByID(input.CreatorID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Créateur introuvable"})
		return
	}
	if creator.MonthlyPrice <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ce créateur n'a pas défini de prix d'abonnement payant"})
		log.Printf("DEBUG: creator.MonthlyPrice = %v", creator.MonthlyPrice)
		return
	}

	log.Printf("DEBUG: creator.MonthlyPrice = %v", creator.MonthlyPrice)

	subscriber, err := user.GetUserByID(uint(subscriberID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Abonné introuvable"})
		return
	}
	customerEmail := subscriber.Email

	successURL := os.Getenv("STRIPE_SUCCESS_URL")
	cancelURL := os.Getenv("STRIPE_CANCEL_URL")

	metadata := map[string]string{
		"creator_id":    strconv.Itoa(int(input.CreatorID)),
		"subscriber_id": strconv.Itoa(subscriberID),
	}

	_, url, err := payment.CreateStripeSubscriptionSession(
		creator.MonthlyPrice,
		"eur",
		successURL,
		cancelURL,
		customerEmail,
		metadata,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur Stripe: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"checkout_url": url})
}
