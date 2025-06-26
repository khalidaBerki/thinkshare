package like

import (
	"errors"
	"gorm.io/gorm"
)

// Repository interface pour les likes
type Repository interface {
	Create(like *Like) error
	Delete(userID, postID uint) error
	GetByUserAndPost(userID, postID uint) (*Like, error)
	CountByPostID(postID uint) (int64, error)
	GetPostLikeStats(postID, userID uint) (*PostLikeStats, error)
	IsLikedByUser(userID, postID uint) (bool, error)
}

// repository implémentation de Repository
type repository struct {
	db *gorm.DB
}

// NewRepository crée une nouvelle instance du repository
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// Create crée un nouveau like
func (r *repository) Create(like *Like) error {
	return r.db.Create(like).Error
}

// Delete supprime un like
func (r *repository) Delete(userID, postID uint) error {
	result := r.db.Where("user_id = ? AND post_id = ?", userID, postID).Delete(&Like{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("like non trouvé")
	}
	return nil
}

// GetByUserAndPost récupère un like spécifique
func (r *repository) GetByUserAndPost(userID, postID uint) (*Like, error) {
	var like Like
	err := r.db.Where("user_id = ? AND post_id = ?", userID, postID).First(&like).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("like non trouvé")
		}
		return nil, err
	}
	return &like, nil
}

// CountByPostID compte le nombre de likes d'un post
func (r *repository) CountByPostID(postID uint) (int64, error) {
	var count int64
	err := r.db.Model(&Like{}).Where("post_id = ?", postID).Count(&count).Error
	return count, err
}

// GetPostLikeStats récupère les statistiques de likes d'un post
func (r *repository) GetPostLikeStats(postID, userID uint) (*PostLikeStats, error) {
	// Compter le total de likes
	totalLikes, err := r.CountByPostID(postID)
	if err != nil {
		return nil, err
	}

	// Vérifier si l'utilisateur a liké
	userHasLiked, err := r.IsLikedByUser(userID, postID)
	if err != nil {
		return nil, err
	}

	return &PostLikeStats{
		PostID:       postID,
		TotalLikes:   int(totalLikes),
		UserHasLiked: userHasLiked,
	}, nil
}

// IsLikedByUser vérifie si un utilisateur a liké un post
func (r *repository) IsLikedByUser(userID, postID uint) (bool, error) {
	var count int64
	err := r.db.Model(&Like{}).Where("user_id = ? AND post_id = ?", userID, postID).Count(&count).Error
	return count > 0, err
}
