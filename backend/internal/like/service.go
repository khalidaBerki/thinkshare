package like

import (
	"backend/internal/post"
	"errors"
)

// PostRepository interface pour vérifier l'existence des posts
type PostRepository interface {
	GetByID(id uint) (*post.Post, error)
}

// Service interface pour la logique métier des likes
type Service interface {
	ToggleLike(userID, postID uint) (*PostLikeStats, error)
	GetPostLikeStats(postID, userID uint) (*PostLikeStats, error)
}

// service implémentation de Service
type service struct {
	repo     Repository
	postRepo PostRepository
}

// NewService crée une nouvelle instance du service
func NewService(repo Repository, postRepo PostRepository) Service {
	return &service{
		repo:     repo,
		postRepo: postRepo,
	}
}

// ToggleLike ajoute ou retire un like sur un post
func (s *service) ToggleLike(userID, postID uint) (*PostLikeStats, error) {
	// Vérifier que l'utilisateur est authentifié
	if userID == 0 {
		return nil, errors.New("utilisateur non authentifié")
	}

	// Vérifier que le post existe
	_, err := s.postRepo.GetByID(postID)
	if err != nil {
		return nil, errors.New("post non trouvé")
	}

	// Vérifier si l'utilisateur a déjà liké ce post
	existingLike, err := s.repo.GetByUserAndPost(userID, postID)

	if existingLike != nil {
		// L'utilisateur a déjà liké -> on retire le like (unlike)
		if err := s.repo.Delete(userID, postID); err != nil {
			return nil, errors.New("erreur lors de la suppression du like")
		}
	} else {
		// L'utilisateur n'a pas encore liké -> on ajoute le like
		like := &Like{
			PostID: postID,
			UserID: userID,
		}

		if err := s.repo.Create(like); err != nil {
			return nil, errors.New("erreur lors de la création du like")
		}
	}

	// Retourner les nouvelles statistiques
	stats, err := s.repo.GetPostLikeStats(postID, userID)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des statistiques")
	}

	return stats, nil
}

// GetPostLikeStats récupère les statistiques de likes d'un post
func (s *service) GetPostLikeStats(postID, userID uint) (*PostLikeStats, error) {
	// Vérifier que le post existe
	_, err := s.postRepo.GetByID(postID)
	if err != nil {
		return nil, errors.New("post non trouvé")
	}

	stats, err := s.repo.GetPostLikeStats(postID, userID)
	if err != nil {
		return nil, errors.New("erreur lors de la récupération des statistiques")
	}

	return stats, nil
}
