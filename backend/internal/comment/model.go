package comment

import (
	_ "backend/internal/postaccess"
	_ "os/user"
	"time"
)

// Comment représente un commentaire sur un post
type Comment struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	PostID    uint      `json:"post_id" gorm:"not null;index"`
	UserID    uint      `json:"user_id" gorm:"not null;index"`
	Text      string    `json:"text" gorm:"type:text;not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateCommentRequest DTO pour créer un commentaire
type CreateCommentRequest struct {
	PostID uint   `json:"post_id" binding:"required"`
	Text   string `json:"text" binding:"required,min=1,max=1000"`
}

// UpdateCommentRequest DTO pour modifier un commentaire
type UpdateCommentRequest struct {
	Text string `json:"text" binding:"required,min=1,max=1000"`
}

// CommentResponse DTO pour la réponse enrichie
type CommentResponse struct {
	ID        uint      `json:"id"`
	PostID    uint      `json:"post_id"`
	UserID    uint      `json:"user_id"`
	Username  string    `json:"username"`
	AvatarURL string    `json:"avatar_url"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
