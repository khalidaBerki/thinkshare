package post

import (
	"backend/internal/db"
	"errors"
	"fmt"
	"os"
)

type Repository interface {
	Create(post *Post) error
	GetAll(page, pageSize int, visibility Visibility) ([]Post, error)
	GetByID(id uint) (*Post, error)
	Update(post *Post) error
	Delete(id uint, userID uint) error
}

type repository struct{}

func NewRepository() Repository {
	return &repository{}
}

func (r *repository) Create(post *Post) error {
	return db.GormDB.Create(post).Error
}

func (r *repository) GetAll(page, pageSize int, visibility Visibility) ([]Post, error) {
	var posts []Post
	offset := (page - 1) * pageSize

	query := db.GormDB.Preload("Media")

	// Filtrer par visibilité si spécifiée
	if visibility != "" {
		query = query.Where("visibility = ?", visibility)
	}

	err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&posts).Error
	return posts, err
}

func (r *repository) GetByID(id uint) (*Post, error) {
	var post Post
	err := db.GormDB.Preload("Media").First(&post, id).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *repository) Update(post *Post) error {
	return db.GormDB.Save(post).Error
}

func (r *repository) Delete(id uint, userID uint) error {
	tx := db.GormDB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var post Post
	if err := tx.Preload("Media").First(&post, id).Error; err != nil {
		tx.Rollback()
		return err
	}

	if post.CreatorID != userID {
		tx.Rollback()
		return errors.New("unauthorized")
	}

	// 1. D'abord supprimer les entrées de médias dans la base de données
	if len(post.Media) > 0 {
		fmt.Printf("Suppression de %d médias pour le post ID %d\n", len(post.Media), post.ID)
		if err := tx.Delete(&post.Media).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("échec de la suppression des médias en base: %v", err)
		}
	}

	// 2. Ensuite supprimer les fichiers physiques
	for _, m := range post.Media {
		// Ignorer l'erreur si le fichier n'existe pas déjà
		err := os.Remove(m.MediaURL)
		if err != nil && !os.IsNotExist(err) {
			// Journaliser l'erreur mais continuer
			fmt.Printf("Avertissement: Impossible de supprimer le fichier %s: %v\n", m.MediaURL, err)
		}
	}

	// 3. Enfin, supprimer le post
	if err := tx.Delete(&post).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
