package media

type Media struct {
	ID           uint `gorm:"primaryKey"`
	PostID       uint
	MediaURL     string
	MediaType    string
	ThumbnailURL string // URL de la miniature pour les images et vidéos
	Metadata     string `gorm:"type:text"` // Métadonnées au format JSON
	FileSize     int64  `gorm:"default:0"` // Taille du fichier en octets
	FileName     string // Nom original du fichier
}
