package postaccess

type PostAccess struct {
	ID        uint `gorm:"primaryKey"`
	PostID    uint
	UserID    uint
	CommentID uint 
}
