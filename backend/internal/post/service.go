package post

import (
	"errors"
	"log"
	"strings"
)

// Interface du service post
type Service interface {
	CreatePost(p *Post) error
	GetAllPosts(page, pageSize int, visibility Visibility) ([]Post, error)
	GetPostByID(id uint) (*Post, error)
	UpdatePost(postID uint, userID uint, input UpdatePostInput) error
	DeletePost(postID uint, userID uint) error
	GetMediaStatistics() (map[string]interface{}, error)
}

// Implémentation du service
type service struct {
	repo Repository
}

// Créer un nouveau service
func NewService(r Repository) Service {
	if r == nil {
		panic("repository cannot be nil")
	}
	return &service{repo: r}
}

// Créer un nouveau post
func (s *service) CreatePost(p *Post) error {
	// Valider les champs obligatoires
	if p.Visibility != Public && p.Visibility != Private {
		return errors.New("invalid visibility")
	}
	if strings.TrimSpace(p.Content) == "" && len(p.Media) == 0 {
		return errors.New("post must have either content or media")
	}

	// Vérification des limites de taille pour l'ensemble des médias
	var totalSize int64
	for _, m := range p.Media {
		totalSize += m.FileSize
	}

	// Limite totale pour tous les médias combinés (200MB)
	const maxTotalSize = 200 * 1024 * 1024 // 200 MB
	if totalSize > maxTotalSize {
		log.Printf("❌ Taille totale des médias trop grande: %.2f MB (max 200MB)",
			float64(totalSize)/(1024*1024))
		return errors.New("total media size exceeds the maximum allowed (200MB)")
	}

	return s.repo.Create(p)
}

// Récupérer tous les posts avec pagination
func (s *service) GetAllPosts(page, pageSize int, visibility Visibility) ([]Post, error) {
	return s.repo.GetAll(page, pageSize, visibility)
}

// Récupérer un post par son ID
func (s *service) GetPostByID(id uint) (*Post, error) {
	return s.repo.GetByID(id)
}

// Mettre à jour un post existant
func (s *service) UpdatePost(postID uint, userID uint, input UpdatePostInput) error {
	// Vérifier que le post existe et appartient à l'utilisateur
	post, err := s.repo.GetByID(postID)
	if err != nil {
		return err
	}
	if post.CreatorID != userID {
		return errors.New("unauthorized")
	}

	// Mettre à jour les champs
	post.Content = input.Content
	post.Visibility = Visibility(input.Visibility)

	return s.repo.Update(post)
}

// Supprimer un post et ses médias associés
func (s *service) DeletePost(postID uint, userID uint) error {
	// Vérifier que le post existe et appartient à l'utilisateur
	post, err := s.repo.GetByID(postID)
	if err != nil {
		return err
	}
	if post.CreatorID != userID {
		return errors.New("unauthorized")
	}
	return s.repo.Delete(postID, userID)
}

// Récupérer des statistiques sur les médias
func (s *service) GetMediaStatistics() (map[string]interface{}, error) {
	// Obtenir les statistiques de base sur les médias
	result := make(map[string]interface{})

	// 1. Nombre total de médias par type
	imageCount, err := s.repo.CountMediaByType("image")
	if err != nil {
		return nil, err
	}

	videoCount, err := s.repo.CountMediaByType("video")
	if err != nil {
		return nil, err
	}

	documentCount, err := s.repo.CountMediaByType("document")
	if err != nil {
		return nil, err
	}

	// 2. Statistiques sur les types de médias
	result["total_media"] = imageCount + videoCount + documentCount
	result["media_by_type"] = map[string]int64{
		"image":    imageCount,
		"video":    videoCount,
		"document": documentCount,
	}

	// 3. Formats supportés
	result["supported_formats"] = map[string]interface{}{
		"images":    getImageFormatList(),
		"videos":    getVideoFormatList(),
		"documents": getDocumentFormatList(),
	}

	return result, nil
}
