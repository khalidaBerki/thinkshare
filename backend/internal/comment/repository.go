package comment

import (
	"errors"
	"gorm.io/gorm"
)

// Repository interface pour l'accès aux données des commentaires
type Repository interface {
	Create(comment *Comment) error
	GetByID(id uint) (*Comment, error)
	GetByPostID(postID uint, limit, offset int) ([]Comment, error)
	Update(comment *Comment) error
	Delete(id uint) error
	CountByPostID(postID uint) (int64, error)
}

// repository implémentation de Repository
type repository struct {
	db *gorm.DB
}

// NewRepository crée une nouvelle instance du repository
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// Create crée un nouveau commentaire
func (r *repository) Create(comment *Comment) error {
	return r.db.Create(comment).Error
}

// GetByID récupère un commentaire par son ID
func (r *repository) GetByID(id uint) (*Comment, error) {
	var comment Comment
	if err := r.db.First(&comment, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("commentaire non trouvé")
		}
		return nil, err
	}
	return &comment, nil
}

// GetByPostID récupère tous les commentaires d'un post avec pagination
func (r *repository) GetByPostID(postID uint, limit, offset int) ([]Comment, error) {
	var comments []Comment
	query := r.db.Where("post_id = ?", postID).Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&comments).Error; err != nil {
		return nil, err
	}
	return comments, nil
}

// Update met à jour un commentaire
func (r *repository) Update(comment *Comment) error {
	return r.db.Save(comment).Error
}

// Delete supprime un commentaire
func (r *repository) Delete(id uint) error {
	return r.db.Delete(&Comment{}, id).Error
}

// CountByPostID compte le nombre de commentaires d'un post
func (r *repository) CountByPostID(postID uint) (int64, error) {
	var count int64
	if err := r.db.Model(&Comment{}).Where("post_id = ?", postID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
