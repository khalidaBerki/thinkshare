package post

import (
	"time"

	"backend/internal/comment"
	"backend/internal/like"
	"backend/internal/media"
	"backend/internal/postaccess"
)

type Post struct {
	ID         uint `gorm:"primaryKey"`
	CreatorID  uint
	Content    string
	Visibility string
	CreatedAt  time.Time

	Media      []media.Media
	Comments   []comment.Comment
	Likes      []like.Like
	PostAccess []postaccess.PostAccess
}
