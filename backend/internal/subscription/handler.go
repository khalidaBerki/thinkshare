package subscription

import (
	"net/http"
	"strconv"
	"time"

	"backend/internal/db"

	"github.com/gin-gonic/gin"
)

// SubscriptionInput pour la requête
type SubscriptionInput struct {
	CreatorID uint   `json:"creator_id" binding:"required"`
	Type      string `json:"type" binding:"required,oneof=paid free"`
}

// SubscribeHandler godoc
// @Summary S’abonner à un créateur (payant ou gratuit)
// @Tags Subscription
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body SubscriptionInput true "Données d’abonnement"
// @Success 200 {object} Subscription
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/subscribe [post]
func SubscribeHandler(c *gin.Context) {
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

	var existing Subscription
	err := db.GormDB.Where("subscriber_id = ? AND creator_id = ?", subscriberID, input.CreatorID).First(&existing).Error

	// Si déjà abonné
	if err == nil {
		// Si déjà abonné "paid" et on redemande "paid", on bloque
		if existing.Type == "paid" && input.Type == "paid" && existing.IsActive && existing.EndDate.After(time.Now()) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Vous êtes déjà abonné payant, renouvellement impossible avant expiration"})
			return
		}

		// Si le type d'abonnement change (free <-> paid)
		if existing.Type != input.Type {
			existing.Type = input.Type
			existing.IsActive = true
			if input.Type == "paid" {
				existing.StartDate = time.Now()
				existing.EndDate = time.Now().AddDate(0, 1, 0)
			} else {
				existing.EndDate = time.Time{}
			}
			if err := db.GormDB.Save(&existing).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la mise à jour de l'abonnement"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"message": "Abonnement mis à jour", "subscription": existing})
			return
		}

		// Si même type "free", ne rien faire
		if input.Type == "free" {
			c.JSON(http.StatusOK, gin.H{"message": "Déjà abonné gratuitement", "subscription": existing})
			return
		}

		// Si même type "paid" mais abonnement expiré, on autorise le renouvellement
		if input.Type == "paid" && (!existing.IsActive || existing.EndDate.Before(time.Now())) {
			now := time.Now()
			existing.StartDate = now
			existing.EndDate = now.AddDate(0, 1, 0)
			existing.IsActive = true
			if err := db.GormDB.Save(&existing).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors du renouvellement"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"message": "Abonnement renouvelé", "subscription": existing})
			return
		}
	}

	// Pas d'abonnement existant, on crée
	sub := Subscription{
		SubscriberID: uint(subscriberID),
		CreatorID:    input.CreatorID,
		StartDate:    time.Now(),
		IsActive:     true,
		Type:         input.Type,
	}
	if input.Type == "paid" {
		sub.EndDate = sub.StartDate.AddDate(0, 1, 0)
	}
	if err := db.GormDB.Create(&sub).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de l'abonnement"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Abonnement réussi", "subscription": sub})
}

// UnsubscribeHandler godoc
// @Summary Se désabonner d’un créateur
// @Tags Subscription
// @Security BearerAuth
// @Param creator_id query int true "ID du créateur"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/unsubscribe [post]
func UnsubscribeHandler(c *gin.Context) {
	creatorID, err := strconv.Atoi(c.Query("creator_id"))
	if err != nil || creatorID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "creator_id invalide"})
		return
	}

	subscriberID := c.GetInt("user_id")

	var sub Subscription
	if err := db.GormDB.Where("subscriber_id = ? AND creator_id = ?", subscriberID, creatorID).First(&sub).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Abonnement non trouvé"})
		return
	}

	// Désactive l'abonnement (soft delete)
	sub.IsActive = false
	if err := db.GormDB.Save(&sub).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors du désabonnement"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Désabonnement réussi"})
}
