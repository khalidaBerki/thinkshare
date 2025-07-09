package message

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// Service définit la logique métier pour les messages privés.
type Service interface {
	Send(senderID uint, input CreateMessageInput) (*MessageDTO, error)
	GetConversation(user1ID, user2ID uint) ([]*MessageDTO, error)
	GetPreviews(userID uint) ([]*MessagePreviewDTO, error)
	MarkRead(senderID, receiverID uint) error
	UpdateMessage(msgID, userID uint, input UpdateMessageInput) (*MessageDTO, error)
	DeleteMessage(msgID, userID uint) error
}

type service struct {
	repo Repository
	db   *gorm.DB
}

type UpdateMessageInput struct {
	Content string `json:"content" binding:"required"`
}

// NewService initialise un nouveau service de messagerie.
func NewService(repo Repository, db *gorm.DB) Service {
	if repo == nil || db == nil {
		panic("message repository and db cannot be nil")
	}
	return &service{repo: repo, db: db}
}

// Send crée un message et renvoie son DTO enrichi.
func (s *service) Send(senderID uint, input CreateMessageInput) (*MessageDTO, error) {
	if senderID == input.ReceiverID {
		return nil, errors.New("you can't send a message to yourself")
	}

	msg := &Message{
		SenderID:   senderID,
		ReceiverID: input.ReceiverID,
		Content:    input.Content,
		Status:     StatusUnread,
	}

	if err := s.repo.CreateMessage(msg); err != nil {
		return nil, err
	}

	// Enrichir avec les infos utilisateur
	senderInfo, err := s.getUserInfoByID(senderID)
	if err != nil {
		return nil, err
	}
	receiverInfo, err := s.getUserInfoByID(input.ReceiverID)
	if err != nil {
		return nil, err
	}

	return &MessageDTO{
		ID:        msg.ID,
		Content:   msg.Content,
		Status:    msg.Status,
		CreatedAt: msg.CreatedAt,
		Sender:    senderInfo,
		Receiver:  receiverInfo,
	}, nil
}

// GetConversation récupère tous les messages entre deux utilisateurs, enrichis.
func (s *service) GetConversation(user1ID, user2ID uint) ([]*MessageDTO, error) {
	msgs, err := s.repo.GetConversation(user1ID, user2ID)
	if err != nil {
		return nil, err
	}

	user1, err := s.getUserInfoByID(user1ID)
	if err != nil {
		return nil, err
	}
	user2, err := s.getUserInfoByID(user2ID)
	if err != nil {
		return nil, err
	}

	var dtos []*MessageDTO
	for _, m := range msgs {
		dto := &MessageDTO{
			ID:        m.ID,
			Content:   m.Content,
			Status:    m.Status,
			CreatedAt: m.CreatedAt,
		}

		if m.SenderID == user1ID {
			dto.Sender = user1
			dto.Receiver = user2
		} else {
			dto.Sender = user2
			dto.Receiver = user1
		}

		dtos = append(dtos, dto)
	}
	return dtos, nil
}

// GetPreviews retourne un aperçu des dernières conversations avec chaque utilisateur.
func (s *service) GetPreviews(userID uint) ([]*MessagePreviewDTO, error) {
	rawPreviews, err := s.repo.GetConversationPreviews(userID)
	if err != nil {
		return nil, err
	}

	var previews []*MessagePreviewDTO
	for _, raw := range rawPreviews {
		dto := &MessagePreviewDTO{
			ConversationID: generateConversationKey(userID, raw.OtherUserID),
			LastMessage:    raw.LastMessage,
			Timestamp:      raw.CreatedAt,
			UnreadCount:    raw.UnreadCount,
			OtherUser: &UserInfo{
				ID:        raw.OtherUserID,
				Username:  raw.OtherUsername,
				AvatarURL: raw.OtherAvatarURL,
			},
		}
		previews = append(previews, dto)
	}
	return previews, nil
}

// MarkRead marque tous les messages de sender vers receiver comme lus.
func (s *service) MarkRead(senderID, receiverID uint) error {
	return s.repo.MarkMessagesAsRead(senderID, receiverID)
}

// getUserInfoByID récupère les infos publiques (username, avatar) d'un utilisateur.
func (s *service) getUserInfoByID(userID uint) (*UserInfo, error) {
	var user struct {
		ID        uint
		Username  string
		AvatarURL string
	}
	err := s.db.
		Table("users").
		Select("id, username, avatar_url").
		Where("id = ?", userID).
		First(&user).Error

	if err != nil {
		return nil, err
	}

	return &UserInfo{
		ID:        user.ID,
		Username:  user.Username,
		AvatarURL: user.AvatarURL,
	}, nil
}

// Utilitaire : génère une clé unique pour une conversation.
func generateConversationKey(user1ID, user2ID uint) string {
	if user1ID < user2ID {
		return fmt.Sprintf("%d-%d", user1ID, user2ID)
	}
	return fmt.Sprintf("%d-%d", user2ID, user1ID)
}

// UpdateMessage met à jour le contenu d'un message si l'utilisateur est l'expéditeur.
func (s *service) UpdateMessage(msgID, userID uint, input UpdateMessageInput) (*MessageDTO, error) {
	// Met à jour le message
	err := s.repo.UpdateMessage(msgID, userID, input.Content)
	if err != nil {
		return nil, err
	}

	// Récupère le message mis à jour
	updated, err := s.repo.GetMessageByID(msgID)
	if err != nil {
		return nil, err
	}

	senderInfo, _ := s.getUserInfoByID(updated.SenderID)
	receiverInfo, _ := s.getUserInfoByID(updated.ReceiverID)

	return &MessageDTO{
		ID:        updated.ID,
		Content:   updated.Content,
		Status:    updated.Status,
		CreatedAt: updated.CreatedAt,
		Sender:    senderInfo,
		Receiver:  receiverInfo,
	}, nil
}

// DeleteMessage supprime un message si l'utilisateur est l'expéditeur.
func (s *service) DeleteMessage(msgID, userID uint) error {
	return s.repo.DeleteMessage(msgID, userID)
}
