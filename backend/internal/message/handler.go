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
// SendMessage godoc
// @Summary      Send a private message
// @Description  Send a private message to another user
// @Tags         messages
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body  body  message.CreateMessageInput  true  "Message content and receiver ID"
// @Success      201   {object}  message.MessageDTO
// @Failure      400   {object}  map[string]string "Invalid input"
// @Failure      401   {object}  map[string]string "Unauthorized"
// @Router       /api/messages [post]
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
// GetPreviews godoc
// @Summary      Get all conversations
// @Description  Get a preview of all conversations (last message, user info, unread count)
// @Tags         messages
// @Security     BearerAuth
// @Produce      json
// @Success      200   {array}   message.MessagePreviewDTO
// @Failure      401   {object}  map[string]string "Unauthorized"
// @Failure      500   {object}  map[string]string "Internal server error"
// @Router       /api/messages/conversations [get]
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
// GetConversation godoc
// @Summary      Get conversation with a user
// @Description  Get all messages exchanged with a specific user
// @Tags         messages
// @Security     BearerAuth
// @Produce      json
// @Param        otherUserID  path  int  true  "Other user ID"
// @Success      200   {array}   message.MessageDTO
// @Failure      400   {object}  map[string]string "Invalid user ID"
// @Failure      401   {object}  map[string]string "Unauthorized"
// @Failure      500   {object}  map[string]string "Internal server error"
// @Router       /api/messages/{otherUserID} [get]
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
// MarkAsRead godoc
// @Summary      Mark messages as read
// @Description  Mark all messages from a sender as read for the authenticated user
// @Tags         messages
// @Security     BearerAuth
// @Produce      json
// @Param        senderID  path  int  true  "Sender user ID"
// @Success      200   {object}  map[string]string "Messages marked as read"
// @Failure      400   {object}  map[string]string "Invalid sender ID"
// @Failure      401   {object}  map[string]string "Unauthorized"
// @Failure      500   {object}  map[string]string "Internal server error"
// @Router       /api/messages/{senderID}/read [patch]
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
// UpdateMessage godoc
// @Summary      Update a message
// @Description  Update the content of a message (only the sender can update)
// @Tags         messages
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id    path  int  true  "Message ID"
// @Param        body  body  message.UpdateMessageInput  true  "Updated content"
// @Success      200   {object}  message.MessageDTO
// @Failure      400   {object}  map[string]string "Invalid input"
// @Failure      401   {object}  map[string]string "Unauthorized"
// @Failure      500   {object}  map[string]string "Internal server error"
// @Router       /api/messages/{id} [put]
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
// DeleteMessage godoc
// @Summary      Delete a message
// @Description  Delete a message (only the sender can delete)
// @Tags         messages
// @Security     BearerAuth
// @Param        id    path  int  true  "Message ID"
// @Success      200   {object}  map[string]string "Message deleted"
// @Failure      400   {object}  map[string]string "Invalid message ID"
// @Failure      401   {object}  map[string]string "Unauthorized"
// @Failure      500   {object}  map[string]string "Internal server error"
// @Router       /api/messages/{id} [delete]
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
