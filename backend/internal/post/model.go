package post

import (
	"time"

	"backend/internal/comment"
	"backend/internal/like"
	"backend/internal/media"
	"backend/internal/postaccess"
)

type Visibility string

const (
	Public  Visibility = "public"
	Private Visibility = "private"
)

// DTO pour créer un post (multipart)
type CreatePostInput struct {
	Content    string     `json:"content" binding:"required"`
	Visibility Visibility `json:"visibility" binding:"required,oneof=public private"`
	// Pour simplifier, on peut avoir un champ "MediaIDs" ou "Media" plus tard (fichiers uploadés)
	// On gérera le multipart/form-data dans handler.go, ce n’est pas dans l’input JSON classique
}

type UpdatePostInput struct {
	Content    string     `json:"content"`
	Visibility Visibility `json:"visibility" binding:"omitempty,oneof=public private"`
}

type Post struct {
	ID         uint       `gorm:"primaryKey"`
	CreatorID  uint       `gorm:"not null"`
	Content    string     `gorm:"type:text"`
	Visibility Visibility `gorm:"type:varchar(10);default:'public'"`
	CreatedAt  time.Time
	UpdatedAt  time.Time

	Media      []media.Media           `gorm:"foreignKey:PostID"`
	Comments   []comment.Comment       `gorm:"foreignKey:PostID"`
	Likes      []like.Like             `gorm:"foreignKey:PostID"`
	PostAccess []postaccess.PostAccess `gorm:"foreignKey:PostID"`
}

// Pour afficher un post
type PostDTO struct {
	ID         uint      `json:"id"`
	Content    string    `json:"content"`
	Visibility string    `json:"visibility"`
	CreatedAt  time.Time `json:"created_at"`
	MediaURLs  []string  `json:"media_urls"`
}
