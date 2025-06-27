package message

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// RegisterRoutes ajoute les routes liées à la messagerie
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	msg := rg.Group("/messages")
	msg.Use() // Auth middleware déjà appliqué au niveau de rg

	msg.POST("", h.SendMessage)
	msg.GET("/conversations", h.GetPreviews)
	msg.GET("/:otherUserID", h.GetConversation)
	msg.PATCH("/:senderID/read", h.MarkAsRead)

	msg.PUT("/:id", h.UpdateMessage)
	msg.DELETE("/:id", h.DeleteMessage)
}

// POST /messages
func (h *Handler) SendMessage(c *gin.Context) {
	var input CreateMessageInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	senderID := c.GetInt("user_id")
	dto, err := h.service.Send(uint(senderID), input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, dto)
}

// GET /messages/conversations
func (h *Handler) GetPreviews(c *gin.Context) {
	userID := c.GetInt("user_id")
	previews, err := h.service.GetPreviews(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load conversations"})
		return
	}
	c.JSON(http.StatusOK, previews)
}

// GET /messages/:otherUserID
func (h *Handler) GetConversation(c *gin.Context) {
	userID := c.GetInt("user_id")
	otherID, err := strconv.Atoi(c.Param("otherUserID"))
	if err != nil || otherID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	convo, err := h.service.GetConversation(uint(userID), uint(otherID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load messages"})
		return
	}
	c.JSON(http.StatusOK, convo)
}

// PATCH /messages/:senderID/read
func (h *Handler) MarkAsRead(c *gin.Context) {
	receiverID := c.GetInt("user_id")
	senderID, err := strconv.Atoi(c.Param("senderID"))
	if err != nil || senderID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sender ID"})
		return
	}

	err = h.service.MarkRead(uint(senderID), uint(receiverID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark messages as read"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Messages marked as read"})
}

// PUT /messages/:id
func (h *Handler) UpdateMessage(c *gin.Context) {
	userID := c.GetInt("user_id")
	msgID, err := strconv.Atoi(c.Param("id"))
	if err != nil || msgID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message ID"})
		return
	}

	var input UpdateMessageInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	updated, err := h.service.UpdateMessage(uint(msgID), uint(userID), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update message"})
		return
	}
	c.JSON(http.StatusOK, updated)
}

// DELETE /messages/:id
func (h *Handler) DeleteMessage(c *gin.Context) {
	userID := c.GetInt("user_id")
	msgID, err := strconv.Atoi(c.Param("id"))
	if err != nil || msgID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message ID"})
		return
	}

	err = h.service.DeleteMessage(uint(msgID), uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete message"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Message deleted"})
}
