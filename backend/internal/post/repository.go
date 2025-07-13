package post

import (
	"backend/internal/db"
	"backend/internal/media"

	userModel "backend/internal/user"
	"errors"

	"gorm.io/gorm"
)

type Repository interface {
	Create(post *Post) error
	GetByID(id uint) (*Post, error)
	GetAll(page, limit int) ([]*Post, int64, error)
	GetByCreatorID(creatorID uint, page, limit int) ([]*Post, int64, error)
	Update(post *Post) error
	Delete(id uint) error

	// Méthodes pour les statistiques
	GetPostStats(postID, userID uint) (*PostStats, error)
	GetPostsWithStats(posts []*Post, userID uint) ([]*PostDTO, error)

	// Ajout de la méthode manquante
	GetCreatorInfo(userID uint) (*CreatorInfo, error)

	CountMediaByType(mediaType string) (int64, error)

	// Méthodes pour le scroll infini
	GetAllAfter(afterID uint, limit int) ([]*Post, error)
	GetByCreatorAfter(creatorID, afterID uint, limit int) ([]*Post, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository() Repository {
	return &repository{db: db.GormDB}
}

func (r *repository) Create(post *Post) error {
	if post == nil {
		return errors.New("post cannot be nil")
	}
	// Sauvegarde le post
	if err := r.db.Create(post).Error; err != nil {
		return err
	}
	// Sauvegarde les médias associés (si présents)
	if len(post.Media) > 0 {
		for i := range post.Media {
			post.Media[i].PostID = post.ID
			post.Media[i].ID = 0 // Laisse GORM gérer l'auto-incrément
		}
		// Crée tous les médias en une seule requête
		if err := r.db.Create(&post.Media).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *repository) GetByID(id uint) (*Post, error) {
	var post Post
	err := r.db.Preload("Media").First(&post, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("post not found")
		}
		return nil, err
	}
	return &post, nil
}

func (r *repository) GetAll(page, limit int) ([]*Post, int64, error) {
	var posts []*Post
	var total int64

	// Compter le total
	if err := r.db.Model(&Post{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Récupérer les posts avec pagination
	offset := (page - 1) * limit
	err := r.db.Preload("Media").Order("id DESC").Offset(offset).Limit(limit).Find(&posts).Error
	if err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}

func (r *repository) GetByCreatorID(creatorID uint, page, limit int) ([]*Post, int64, error) {
	var posts []*Post
	var total int64

	// Compter le total
	if err := r.db.Model(&Post{}).Where("creator_id = ?", creatorID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Récupérer les posts avec pagination
	offset := (page - 1) * limit
	err := r.db.Preload("Media").Where("creator_id = ?", creatorID).Order("created_at DESC").Offset(offset).Limit(limit).Find(&posts).Error
	if err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}

func (r *repository) Update(post *Post) error {
	if post == nil {
		return errors.New("post cannot be nil")
	}
	return r.db.Save(post).Error
}

func (r *repository) Delete(id uint) error {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Supprimer les médias associés
	if err := tx.Where("post_id = ?", id).Delete(&media.Media{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Supprimer le post
	if err := tx.Delete(&Post{}, id).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// GetPostStats récupère les statistiques d'un post
func (r *repository) GetPostStats(postID, userID uint) (*PostStats, error) {
	stats := &PostStats{PostID: postID}

	// Compter les likes
	var likeCount int64
	if err := r.db.Table("likes").Where("post_id = ?", postID).Count(&likeCount).Error; err != nil {
		return nil, err
	}
	stats.LikeCount = int(likeCount)

	// Compter les commentaires
	var commentCount int64
	if err := r.db.Table("comments").Where("post_id = ?", postID).Count(&commentCount).Error; err != nil {
		return nil, err
	}
	stats.CommentCount = int(commentCount)

	// Vérifier si l'utilisateur a liké
	if userID > 0 {
		var exists bool
		err := r.db.Table("likes").Select("count(*) > 0").Where("post_id = ? AND user_id = ?", postID, userID).Find(&exists).Error
		if err != nil {
			return nil, err
		}
		stats.UserHasLiked = exists
	}

	return stats, nil
}

// GetPostsWithStats convertit les posts en PostDTO avec statistiques
func (r *repository) GetPostsWithStats(posts []*Post, userID uint) ([]*PostDTO, error) {
	if len(posts) == 0 {
		return []*PostDTO{}, nil
	}

	result := make([]*PostDTO, 0, len(posts))

	for _, post := range posts {

		// Récupérer les statistiques
		stats, err := r.GetPostStats(post.ID, userID)
		if err != nil {
			return nil, err
		}

		// Récupérer les URLs des médias
		mediaURLs := make([]string, len(post.Media))
		for i, media := range post.Media {
			mediaURLs[i] = media.MediaURL
		}

		// Récupérer les infos du créateur
		var creator *CreatorInfo
		creatorInfo, err := r.GetCreatorInfo(post.CreatorID)
		if err == nil {
			creator = creatorInfo
		}

		postDTO := &PostDTO{
			ID:           post.ID,
			CreatorID:    post.CreatorID,
			Content:      post.Content,
			Visibility:   string(post.Visibility),
			IsPaidOnly:   post.IsPaidOnly, // <-- Ajouté pour le mapping correct
			DocumentType: post.DocumentType,
			CreatedAt:    post.CreatedAt,
			UpdatedAt:    post.UpdatedAt,
			MediaURLs:    mediaURLs,
			LikeCount:    stats.LikeCount,
			CommentCount: stats.CommentCount,
			UserHasLiked: stats.UserHasLiked,
			Creator:      creator,
		}

		result = append(result, postDTO)
	}

	return result, nil
}

// GetCreatorInfo récupère les informations du créateur d'un post
func (r *repository) GetCreatorInfo(userID uint) (*CreatorInfo, error) {
	var user userModel.User // adapte selon ton modèle utilisateur
	if err := r.db.First(&user, userID).Error; err != nil {
		return nil, err
	}
	return &CreatorInfo{
		ID:           user.ID,
		Username:     user.Username,
		FullName:     user.FullName,
		AvatarURL:    user.AvatarURL,
		MonthlyPrice: user.MonthlyPrice, // Ajout pour le feed Flutter
	}, nil
}

func (r *repository) CountMediaByType(mediaType string) (int64, error) {
	var count int64
	err := r.db.Model(&media.Media{}).Where("media_type = ?", mediaType).Count(&count).Error
	return count, err
}

// Récupère tous les posts après un certain ID (scroll infini)
func (r *repository) GetAllAfter(afterID uint, limit int) ([]*Post, error) {
	var posts []*Post
	query := r.db.Preload("Media").Order("id DESC").Limit(limit)
	if afterID > 0 {
		query = query.Where("id < ?", afterID)
	}
	err := query.Find(&posts).Error
	return posts, err
}

// Récupère les posts d'un créateur après un certain ID (scroll infini)
func (r *repository) GetByCreatorAfter(creatorID, afterID uint, limit int) ([]*Post, error) {
	var posts []*Post
	query := r.db.Preload("Media").Where("creator_id = ?", creatorID).Order("id ASC").Limit(limit)
	if afterID > 0 {
		query = query.Where("id > ?", afterID).Where("creator_id = ?", creatorID)
	}
	err := query.Find(&posts).Error
	return posts, err
}
