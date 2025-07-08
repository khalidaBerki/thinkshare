package message

import (
	"errors"
	_ "fmt"

	"gorm.io/gorm"
)

type Repository interface {
	CreateMessage(msg *Message) error
	GetConversation(user1ID, user2ID uint) ([]*Message, error)
	GetConversationPreviews(userID uint) ([]*MessagePreviewRaw, error)
	MarkMessagesAsRead(senderID, receiverID uint) error
	UpdateMessage(msgID, userID uint, content string) error
	DeleteMessage(msgID, userID uint) error
	GetMessageByID(msgID uint) (*Message, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

// Create a new message
func (r *repository) CreateMessage(msg *Message) error {
	return r.db.Create(msg).Error
}

// Get full conversation between two users (ordered by CreatedAt)
func (r *repository) GetConversation(user1ID, user2ID uint) ([]*Message, error) {
	var messages []*Message
	err := r.db.
		Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)",
			user1ID, user2ID, user2ID, user1ID).
		Order("created_at ASC").
		Find(&messages).Error

	return messages, err
}

// Get preview of all conversations with last message, user info and unread count
func (r *repository) GetConversationPreviews(userID uint) ([]*MessagePreviewRaw, error) {
	var previews []*MessagePreviewRaw

	// Subquery to find the latest message per conversation
	subquery := r.db.
		Table("messages").
		Select("CASE WHEN sender_id = ? THEN receiver_id ELSE sender_id END AS other_user_id, MAX(created_at) AS last_time", userID).
		Where("sender_id = ? OR receiver_id = ?", userID, userID).
		Group("other_user_id")

	// Join to fetch message + user info + unread count
	tx := r.db.Table("messages AS m").
		Select(`
			m.content AS last_message,
			u.id AS other_user_id,
			u.username AS other_username,
			u.avatar_url AS other_avatar_url,
			m.created_at,
			(
				SELECT COUNT(*) FROM messages AS unread
				WHERE unread.sender_id = m.sender_id
				AND unread.receiver_id = m.receiver_id
				AND unread.status = 'UNREAD'
				AND unread.receiver_id = ?
			) AS unread_count
		`, userID).
		Joins("JOIN users u ON u.id = CASE WHEN m.sender_id = ? THEN m.receiver_id ELSE m.sender_id END", userID).
		Joins("JOIN (?) AS conv ON ((m.sender_id = ? AND m.receiver_id = conv.other_user_id) OR (m.receiver_id = ? AND m.sender_id = conv.other_user_id)) AND m.created_at = conv.last_time", subquery, userID, userID).
		Order("m.created_at DESC")

	if err := tx.Scan(&previews).Error; err != nil {
		return nil, err
	}

	return previews, nil
}

// Mark all unread messages from sender to receiver as read
func (r *repository) MarkMessagesAsRead(senderID, receiverID uint) error {
	res := r.db.Model(&Message{}).
		Where("sender_id = ? AND receiver_id = ? AND status = ?", senderID, receiverID, "UNREAD").
		Update("status", "READ")

	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("no messages to mark as read")
	}
	return nil
}

// Utilitaire pour générer une clé de conversation unique entre 2 utilisateurs (ex: "2-5")
/*func generateConversationID(userID1, userID2 uint) string {
	if userID1 < userID2 {
		return fmt.Sprintf("%d-%d", userID1, userID2)
	}
	return fmt.Sprintf("%d-%d", userID2, userID1)
}*/

// Met à jour le contenu d'un message (seul l'auteur peut modifier)
func (r *repository) UpdateMessage(msgID, userID uint, content string) error {
	res := r.db.Model(&Message{}).
		Where("id = ? AND sender_id = ?", msgID, userID).
		Update("content", content)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("message not found or not owned by user")
	}
	return nil
}

// Supprime un message (seul l'auteur peut supprimer)
func (r *repository) DeleteMessage(msgID, userID uint) error {
	res := r.db.Where("id = ? AND sender_id = ?", msgID, userID).Delete(&Message{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("message not found or not owned by user")
	}
	return nil
}

// GetMessageByID récupère un message par son ID
func (r *repository) GetMessageByID(msgID uint) (*Message, error) {
	var msg Message
	err := r.db.First(&msg, msgID).Error
	if err != nil {
		return nil, err
	}
	return &msg, nil
}
