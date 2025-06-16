package comment

import (
	"backend/internal/postaccess"
	"time"
)

type Comment struct {
	ID         uint `gorm:"primaryKey"`
	PostID     uint
	UserID     uint
	Content    string
	CreatedAt  time.Time
	PostAccess []postaccess.PostAccess `gorm:"foreignKey:CommentID"`
}
