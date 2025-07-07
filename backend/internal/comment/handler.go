package comment

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Handler HTTP for comments
type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	comments := rg.Group("/comments")

	// Routes pour les commentaires
	comments.POST("", h.CreateComment)              // POST /api/comments
	comments.GET("/:postID", h.GetCommentsByPostID) // GET /api/comments/:postID
	comments.PUT("/:id", h.UpdateComment)           // PUT /api/comments/:id
	comments.DELETE("/:id", h.DeleteComment)        // DELETE /api/comments/:id
}

// CreateComment creates a new comment
func (h *Handler) CreateComment(c *gin.Context) {
	userID := c.GetInt("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}
	var req CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data", "details": err.Error()})
		return
	}
	comment, err := h.service.CreateComment(uint(userID), req)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "post non trouvé" {
			status = http.StatusNotFound
			c.JSON(status, gin.H{"error": "Post not found"})
			return
		} else if err.Error() == "utilisateur non authentifié" {
			status = http.StatusUnauthorized
			c.JSON(status, gin.H{"error": "Authentication required"})
			return
		}
		c.JSON(status, gin.H{"error": "Failed to create comment"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "Comment created successfully",
		"comment": comment,
	})
}

// GetCommentsByPostID retrieves comments for a post
func (h *Handler) GetCommentsByPostID(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("postID"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	comments, total, err := h.service.GetCommentsByPostID(uint(postID), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve comments"})
		return
	}
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

// UpdateComment updates a comment
func (h *Handler) UpdateComment(c *gin.Context) {
	userID := c.GetInt("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}
	commentID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}
	var req UpdateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data", "details": err.Error()})
		return
	}
	comment, err := h.service.UpdateComment(uint(userID), uint(commentID), req)
	if err != nil {
		status := http.StatusInternalServerError
		switch err.Error() {
		case "commentaire non trouvé":
			status = http.StatusNotFound
			c.JSON(status, gin.H{"error": "Comment not found"})
			return
		case "vous n'êtes pas autorisé à modifier ce commentaire":
			status = http.StatusForbidden
			c.JSON(status, gin.H{"error": "You are not allowed to edit this comment"})
			return
		case "utilisateur non authentifié":
			status = http.StatusUnauthorized
			c.JSON(status, gin.H{"error": "Authentication required"})
			return
		}
		c.JSON(status, gin.H{"error": "Failed to update comment"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Comment updated successfully",
		"comment": comment,
	})
}

// DeleteComment deletes a comment
func (h *Handler) DeleteComment(c *gin.Context) {
	userID := c.GetInt("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}
	commentID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}
	if err := h.service.DeleteComment(uint(userID), uint(commentID)); err != nil {
		status := http.StatusInternalServerError
		switch err.Error() {
		case "commentaire non trouvé":
			status = http.StatusNotFound
			c.JSON(status, gin.H{"error": "Comment not found"})
			return
		case "vous n'êtes pas autorisé à supprimer ce commentaire":
			status = http.StatusForbidden
			c.JSON(status, gin.H{"error": "You are not allowed to delete this comment"})
			return
		case "utilisateur non authentifié":
			status = http.StatusUnauthorized
			c.JSON(status, gin.H{"error": "Authentication required"})
			return
		}
		c.JSON(status, gin.H{"error": "Failed to delete comment"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Comment deleted successfully",
		"id":      commentID,
	})
}
