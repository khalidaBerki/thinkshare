package media

import (
	"errors"
	"gorm.io/gorm"
)

// Interface du repository pour les médias
type Repository interface {
	Create(media *Media) error
	FindByID(id uint) (*Media, error)
	FindByPostID(postID uint) ([]Media, error)
	FindAll() ([]Media, error)
	Update(media *Media) error
	Delete(id uint) error
}

type repositoryImpl struct {
	db *gorm.DB
}

// Créer une nouvelle instance du repository
func NewRepository(db *gorm.DB) Repository {
	return &repositoryImpl{db: db}
}

// Créer un nouveau média
func (r *repositoryImpl) Create(media *Media) error {
	return r.db.Create(media).Error
}

// Trouver un média par son ID
func (r *repositoryImpl) FindByID(id uint) (*Media, error) {
	var media Media
	result := r.db.First(&media, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("média non trouvé")
		}
		return nil, result.Error
	}
	return &media, nil
}

// Trouver tous les médias associés à un post
func (r *repositoryImpl) FindByPostID(postID uint) ([]Media, error) {
	var medias []Media
	result := r.db.Where("post_id = ?", postID).Find(&medias)
	return medias, result.Error
}

// Trouver tous les médias
func (r *repositoryImpl) FindAll() ([]Media, error) {
	var medias []Media
	result := r.db.Find(&medias)
	return medias, result.Error
}

// Mettre à jour un média
func (r *repositoryImpl) Update(media *Media) error {
	return r.db.Save(media).Error
}

// Supprimer un média
func (r *repositoryImpl) Delete(id uint) error {
	return r.db.Delete(&Media{}, id).Error
}
