package like

import (
	"time"
)

// Like représente un like sur un post
type Like struct {
	ID        uint      `gorm:"primaryKey"`
	PostID    uint      `gorm:"not null;index"`
	UserID    uint      `gorm:"not null;index"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

// Index unique pour empêcher qu'un utilisateur like plusieurs fois le même post
func (Like) TableName() string {
	return "likes"
}

// LikeResponse DTO pour les réponses API
type LikeResponse struct {
	ID        uint      `json:"id"`
	PostID    uint      `json:"post_id"`
	UserID    uint      `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

// ToResponse convertit un Like en LikeResponse
func (l *Like) ToResponse() LikeResponse {
	return LikeResponse{
		ID:        l.ID,
		PostID:    l.PostID,
		UserID:    l.UserID,
		CreatedAt: l.CreatedAt,
	}
}

// PostLikeStats statistiques de likes d'un post
type PostLikeStats struct {
	PostID       uint `json:"post_id"`
	TotalLikes   int  `json:"total_likes"`
	UserHasLiked bool `json:"user_has_liked"`
}
