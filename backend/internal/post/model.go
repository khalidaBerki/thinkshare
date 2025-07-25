package post

import (
	"time"

	"backend/internal/media"
	"backend/internal/postaccess"
)

type Visibility string

const (
	Public  Visibility = "public"
	Private Visibility = "private"
)

// DTO pour créer un post
type CreatePostInput struct {
	Content      string        `json:"content" binding:"required"`
	Visibility   Visibility    `json:"visibility" binding:"required,oneof=public private"`
	IsPaidOnly   bool          `json:"is_paid_only"` // Nouveau champ pour création
	DocumentType string        `json:"document_type,omitempty"`
	Media        []media.Media `json:"media"`
}

type UpdatePostInput struct {
	Content      string     `json:"content"`
	Visibility   Visibility `json:"visibility" binding:"omitempty,oneof=public private"`
	IsPaidOnly   *bool      `json:"is_paid_only,omitempty"` // Pointeur pour permettre la mise à jour
	DocumentType string     `json:"document_type,omitempty"`
}

type Post struct {
	ID           uint       `gorm:"primaryKey"`
	CreatorID    uint       `gorm:"not null;index"`
	Content      string     `gorm:"type:text"`
	Visibility   Visibility `gorm:"type:varchar(10);default:'public'"`
	IsPaidOnly   bool       `gorm:"default:false"` // Nouveau champ pour les posts payants
	DocumentType string     `gorm:"type:varchar(50)"`
	CreatedAt    time.Time
	UpdatedAt    time.Time

	Media      []media.Media           `gorm:"foreignKey:PostID"`
	PostAccess []postaccess.PostAccess `gorm:"foreignKey:PostID"`
}

// PostDTO pour les réponses API
type PostDTO struct {
	ID           uint      `json:"id"`
	CreatorID    uint      `json:"creator_id"`
	Content      string    `json:"content"`
	Visibility   string    `json:"visibility"`
	IsPaidOnly   bool      `json:"is_paid_only"`
	DocumentType string    `json:"document_type,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	MediaURLs    []string  `json:"media_urls"`

	// Statistiques
	LikeCount    int  `json:"like_count"`
	CommentCount int  `json:"comment_count"`
	UserHasLiked bool `json:"user_has_liked"`

	// Informations du créateur
	Creator *CreatorInfo `json:"creator,omitempty"`

	// Contrôle d'accès
	HasAccess bool `json:"has_access"` // Indique si l'utilisateur peut voir le contenu
}

type CreatorInfo struct {
	ID           uint    `json:"id"`
	Username     string  `json:"username"`
	FullName     string  `json:"full_name"`
	AvatarURL    string  `json:"avatar_url,omitempty"`
	MonthlyPrice float64 `json:"monthly_price"` // Ajouté pour le feed Flutter
}

// Vérifie que le modèle User existe bien dans ton projet avec les champs suivants :
// ID uint, Username string, FullName string, AvatarURL string

type PostStats struct {
	PostID       uint `json:"post_id"`
	LikeCount    int  `json:"like_count"`
	CommentCount int  `json:"comment_count"`
	UserHasLiked bool `json:"user_has_liked"`
}
