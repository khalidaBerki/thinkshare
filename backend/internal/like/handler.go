package like

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Handler gestionnaire HTTP pour les likes
type Handler struct {
	service Service
}

// NewHandler crée une nouvelle instance du handler
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// RegisterRoutes enregistre les routes du handler
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	likes := rg.Group("/likes")

	// Routes pour les likes
	likes.POST("/posts/:postID", h.ToggleLike)      // POST /api/likes/posts/:postID
	likes.GET("/posts/:postID", h.GetPostLikeStats) // GET /api/likes/posts/:postID
}

// ToggleLike godoc
// @Summary      Toggle like on a post
// @Description  Add or remove a like on a post by the authenticated user
// @Tags         likes
// @Security     BearerAuth
// @Param        postID  path  int  true  "Post ID"
// @Success      200  {object}  map[string]interface{} "Like toggled, returns stats"
// @Failure      400  {object}  map[string]string "Invalid post ID"
// @Failure      401  {object}  map[string]string "Authentication required"
// @Failure      404  {object}  map[string]string "Post not found"
// @Failure      500  {object}  map[string]string "Internal server error"
// @Router       /api/likes/posts/{postID} [post]
func (h *Handler) ToggleLike(c *gin.Context) {
	// Récupérer l'ID de l'utilisateur depuis le middleware d'authentification
	userID := c.GetInt("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentification requise"})
		return
	}

	// Récupérer l'ID du post
	postID, err := strconv.ParseUint(c.Param("postID"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de post invalide"})
		return
	}

	// Toggle le like
	stats, err := h.service.ToggleLike(uint(userID), uint(postID))
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

	// Déterminer le message selon l'action
	message := "Like ajouté"
	if !stats.UserHasLiked {
		message = "Like retiré"
	}

	c.JSON(http.StatusOK, gin.H{
		"message": message,
		"stats":   stats,
	})
}

// GetPostLikeStats godoc
// @Summary      Get like stats for a post
// @Description  Get the total number of likes and whether the authenticated user has liked the post
// @Tags         likes
// @Security     BearerAuth
// @Param        postID  path  int  true  "Post ID"
// @Success      200  {object}  map[string]interface{} "Like stats"
// @Failure      400  {object}  map[string]string "Invalid post ID"
// @Failure      404  {object}  map[string]string "Post not found"
// @Failure      500  {object}  map[string]string "Internal server error"
// @Router       /api/likes/posts/{postID} [get]
func (h *Handler) GetPostLikeStats(c *gin.Context) {
	// Récupérer l'ID de l'utilisateur (optionnel pour cette route)
	userID := c.GetInt("user_id")

	// Récupérer l'ID du post
	postID, err := strconv.ParseUint(c.Param("postID"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de post invalide"})
		return
	}

	// Récupérer les statistiques
	stats, err := h.service.GetPostLikeStats(uint(postID), uint(userID))
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "post non trouvé" {
			status = http.StatusNotFound
		}

		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"stats": stats,
	})
}
