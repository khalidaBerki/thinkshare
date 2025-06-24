package comment

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Handler gestionnaire HTTP pour les commentaires
type Handler struct {
	service Service
}

// NewHandler crée une nouvelle instance du handler
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// RegisterRoutes enregistre les routes du handler
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	comments := rg.Group("/comments")

	// Routes pour les commentaires
	comments.POST("", h.CreateComment)              // POST /api/comments
	comments.GET("/:postID", h.GetCommentsByPostID) // GET /api/comments/:postID
	comments.PUT("/:id", h.UpdateComment)           // PUT /api/comments/:id
	comments.DELETE("/:id", h.DeleteComment)        // DELETE /api/comments/:id
}

// CreateComment crée un nouveau commentaire
func (h *Handler) CreateComment(c *gin.Context) {
	// 🔧 FIX : Utiliser GetInt au lieu de GetUint
	userID := c.GetInt("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentification requise"})
		return
	}

	// Lire les données de la requête
	var req CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données de requête invalides", "details": err.Error()})
		return
	}

	// Créer le commentaire - conversion uint nécessaire pour le service
	comment, err := h.service.CreateComment(uint(userID), req)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "post non trouvé" {
			status = http.StatusNotFound
		} else if err.Error() == "utilisateur non authentifié" {
			status = http.StatusUnauthorized
		}

		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Commentaire créé avec succès",
		"comment": comment,
	})
}

// GetCommentsByPostID récupère les commentaires d'un post
func (h *Handler) GetCommentsByPostID(c *gin.Context) {
	// Récupérer l'ID du post
	postID, err := strconv.ParseUint(c.Param("postID"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de post invalide"})
		return
	}

	// Récupérer les paramètres de pagination
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	// Récupérer les commentaires
	comments, total, err := h.service.GetCommentsByPostID(uint(postID), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Calculer les informations de pagination
	totalPages := (int(total) + limit - 1) / limit
	hasNext := page < totalPages
	hasPrev := page > 1

	c.JSON(http.StatusOK, gin.H{
		"comments": comments,
		"pagination": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": totalPages,
			"has_next":    hasNext,
			"has_prev":    hasPrev,
		},
	})
}

// UpdateComment met à jour un commentaire
func (h *Handler) UpdateComment(c *gin.Context) {
	// 🔧 FIX : Utiliser GetInt au lieu de GetUint
	userID := c.GetInt("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentification requise"})
		return
	}

	// Récupérer l'ID du commentaire
	commentID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de commentaire invalide"})
		return
	}

	// Lire les données de la requête
	var req UpdateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données de requête invalides", "details": err.Error()})
		return
	}

	// Mettre à jour le commentaire - conversion uint nécessaire pour le service
	comment, err := h.service.UpdateComment(uint(userID), uint(commentID), req)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "commentaire non trouvé" {
			status = http.StatusNotFound
		} else if err.Error() == "vous n'êtes pas autorisé à modifier ce commentaire" {
			status = http.StatusForbidden
		} else if err.Error() == "utilisateur non authentifié" {
			status = http.StatusUnauthorized
		}

		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Commentaire mis à jour avec succès",
		"comment": comment,
	})
}

// DeleteComment supprime un commentaire
func (h *Handler) DeleteComment(c *gin.Context) {
	// 🔧 FIX : Utiliser GetInt au lieu de GetUint
	userID := c.GetInt("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentification requise"})
		return
	}

	// Récupérer l'ID du commentaire
	commentID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de commentaire invalide"})
		return
	}

	// Supprimer le commentaire - conversion uint nécessaire pour le service
	if err := h.service.DeleteComment(uint(userID), uint(commentID)); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "commentaire non trouvé" {
			status = http.StatusNotFound
		} else if err.Error() == "vous n'êtes pas autorisé à supprimer ce commentaire" {
			status = http.StatusForbidden
		} else if err.Error() == "utilisateur non authentifié" {
			status = http.StatusUnauthorized
		}

		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Commentaire supprimé avec succès",
		"id":      commentID,
	})
}
