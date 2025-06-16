package media

type Media struct {
	ID        uint `gorm:"primaryKey"`
	PostID    uint
	MediaURL  string
	MediaType string
}
