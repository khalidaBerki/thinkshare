package access

import "time"

type PostAccess struct {
	ID         uint `gorm:"primaryKey"`
	UserID     uint
	PostID     uint
	AccessType string
	AccessDate time.Time
	Permanent  bool
}
