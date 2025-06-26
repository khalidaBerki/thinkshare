package comment

import (
	"backend/internal/post"
	"errors"
	"time"
)

// PostRepository interface pour vérifier l'existence des posts
type PostRepository interface {
	GetByID(id uint) (*post.Post, error)
}

// Service interface pour la logique métier des commentaires
type Service interface {
	CreateComment(userID uint, req CreateCommentRequest) (*CommentResponse, error)
	GetCommentsByPostID(postID uint, page, limit int) ([]CommentResponse, int64, error)
	UpdateComment(userID, commentID uint, req UpdateCommentRequest) (*CommentResponse, error)
	DeleteComment(userID, commentID uint) error
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

// CreateComment crée un nouveau commentaire
func (s *service) CreateComment(userID uint, req CreateCommentRequest) (*CommentResponse, error) {
	// Vérifier que l'utilisateur est authentifié
	if userID == 0 {
		return nil, errors.New("utilisateur non authentifié")
	}

	// Vérifier que le post existe
	_, err := s.postRepo.GetByID(req.PostID)
	if err != nil {
		return nil, errors.New("post non trouvé")
	}

	// Créer le commentaire
	comment := &Comment{
		PostID:    req.PostID,
		UserID:    userID,
		Text:      req.Text,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.Create(comment); err != nil {
		return nil, errors.New("erreur lors de la création du commentaire")
	}

	response := comment.ToResponse()
	return &response, nil
}

// GetCommentsByPostID récupère les commentaires d'un post avec pagination
func (s *service) GetCommentsByPostID(postID uint, page, limit int) ([]CommentResponse, int64, error) {
	// Valeurs par défaut pour la pagination
	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit

	// Récupérer les commentaires
	comments, err := s.repo.GetByPostID(postID, limit, offset)
	if err != nil {
		return nil, 0, errors.New("erreur lors de la récupération des commentaires")
	}

	// Compter le total
	total, err := s.repo.CountByPostID(postID)
	if err != nil {
		return nil, 0, errors.New("erreur lors du comptage des commentaires")
	}

	// Convertir en réponses
	responses := make([]CommentResponse, len(comments))
	for i, comment := range comments {
		responses[i] = comment.ToResponse()
	}

	return responses, total, nil
}

// UpdateComment met à jour un commentaire
func (s *service) UpdateComment(userID, commentID uint, req UpdateCommentRequest) (*CommentResponse, error) {
	// Vérifier que l'utilisateur est authentifié
	if userID == 0 {
		return nil, errors.New("utilisateur non authentifié")
	}

	// Récupérer le commentaire
	comment, err := s.repo.GetByID(commentID)
	if err != nil {
		return nil, err
	}

	// Vérifier que l'utilisateur est le propriétaire
	if comment.UserID != userID {
		return nil, errors.New("vous n'êtes pas autorisé à modifier ce commentaire")
	}

	// Mettre à jour le commentaire
	comment.Text = req.Text
	comment.UpdatedAt = time.Now()

	if err := s.repo.Update(comment); err != nil {
		return nil, errors.New("erreur lors de la mise à jour du commentaire")
	}

	response := comment.ToResponse()
	return &response, nil
}

// DeleteComment supprime un commentaire
func (s *service) DeleteComment(userID, commentID uint) error {
	// Vérifier que l'utilisateur est authentifié
	if userID == 0 {
		return errors.New("utilisateur non authentifié")
	}

	// Récupérer le commentaire
	comment, err := s.repo.GetByID(commentID)
	if err != nil {
		return err
	}

	// Vérifier que l'utilisateur est le propriétaire
	if comment.UserID != userID {
		return errors.New("vous n'êtes pas autorisé à supprimer ce commentaire")
	}

	// Supprimer le commentaire
	if err := s.repo.Delete(commentID); err != nil {
		return errors.New("erreur lors de la suppression du commentaire")
	}

	return nil
}
