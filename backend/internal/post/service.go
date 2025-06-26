package post

import (
	"errors"
	"strings"
)

type Service interface {
	CreatePost(creatorID uint, input CreatePostInput) (*PostDTO, error)
	GetPostByID(id, userID uint) (*PostDTO, error)
	GetAllPosts(page, limit int, userID uint) ([]*PostDTO, int64, error)
	GetPostsByCreator(creatorID uint, page, limit int, userID uint) ([]*PostDTO, int64, error)
	UpdatePost(postID, creatorID uint, input UpdatePostInput) (*PostDTO, error)
	DeletePost(postID, creatorID uint) error
	GetMediaStatistics() (interface{}, interface{})
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	if repo == nil {
		panic("repository cannot be nil")
	}
	return &service{repo: repo}
}

func (s *service) GetMediaStatistics() (interface{}, interface{}) {
	types := []string{"image", "video", "document"}
	mediaByType := make(map[string]int64)
	var total int64

	for _, t := range types {
		count, err := s.repo.CountMediaByType(t)
		if err != nil {
			mediaByType[t] = 0
		} else {
			mediaByType[t] = count
			total += count
		}
	}

	statistics := map[string]interface{}{
		"total_media":   total,
		"media_by_type": mediaByType,
	}

	recommendations := []string{} //getRecommendedFormats() si on a cette fonction

	return statistics, recommendations
}

// CreatePost crée un nouveau post
func (s *service) CreatePost(creatorID uint, input CreatePostInput) (*PostDTO, error) {
	if creatorID == 0 {
		return nil, errors.New("utilisateur non authentifié")
	}
	if strings.TrimSpace(input.Content) == "" {
		return nil, errors.New("le contenu ne peut pas être vide")
	}
	if input.Visibility != Public && input.Visibility != Private {
		return nil, errors.New("invalid visibility")
	}
	post := &Post{
		CreatorID:    creatorID,
		Content:      strings.TrimSpace(input.Content),
		Visibility:   input.Visibility,
		DocumentType: input.DocumentType,
		Media:        input.Media,
	}
	if err := s.repo.Create(post); err != nil {
		return nil, errors.New("erreur lors de la création du post")
	}
	return s.GetPostByID(post.ID, creatorID)
}

// GetPostByID récupère un post + statistiques + créateur
func (s *service) GetPostByID(id, userID uint) (*PostDTO, error) {
	post, err := s.repo.GetByID(id)
	if err != nil {
		return nil, errors.New("post non trouvé")
	}

	postsDTO, err := s.repo.GetPostsWithStats([]*Post{post}, userID)
	if err != nil || len(postsDTO) == 0 {
		return nil, errors.New("erreur lors de la récupération des statistiques")
	}

	dto := postsDTO[0]
	removeDuplicateMediaURLs(dto) // ✅

	creator, err := s.repo.GetCreatorInfo(post.CreatorID)
	if err == nil {
		dto.Creator = creator
	}
	return dto, nil
}

func removeDuplicateMediaURLs(dto *PostDTO) {
	seen := map[string]bool{}
	unique := []string{}
	for _, url := range dto.MediaURLs {
		if !seen[url] {
			seen[url] = true
			unique = append(unique, url)
		}
	}
	dto.MediaURLs = unique
}

func (s *service) GetAllPosts(page, limit int, userID uint) ([]*PostDTO, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	posts, total, err := s.repo.GetAll(page, limit)
	if err != nil {
		return nil, 0, errors.New("erreur lors de la récupération des posts")
	}

	postsDTO, err := s.repo.GetPostsWithStats(posts, userID)
	if err != nil {
		return nil, 0, errors.New("erreur lors de la récupération des statistiques")
	}

	for _, dto := range postsDTO {
		removeDuplicateMediaURLs(dto) // ✅
	}

	return postsDTO, total, nil
}

// GetPostsByCreator récupère les posts d'un utilisateur donné
func (s *service) GetPostsByCreator(creatorID uint, page, limit int, userID uint) ([]*PostDTO, int64, error) {
	if creatorID == 0 {
		return nil, 0, errors.New("ID créateur invalide")
	}

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	posts, total, err := s.repo.GetByCreatorID(creatorID, page, limit)
	if err != nil {
		return nil, 0, errors.New("erreur lors de la récupération des posts")
	}

	postsDTO, err := s.repo.GetPostsWithStats(posts, userID)
	if err != nil {
		return nil, 0, errors.New("erreur lors de la récupération des statistiques")
	}

	for _, dto := range postsDTO {
		removeDuplicateMediaURLs(dto) // ✅
	}

	return postsDTO, total, nil
}

// UpdatePost met à jour un post
func (s *service) UpdatePost(postID, creatorID uint, input UpdatePostInput) (*PostDTO, error) {
	post, err := s.repo.GetByID(postID)
	if err != nil {
		return nil, errors.New("post non trouvé")
	}

	if post.CreatorID != creatorID {
		return nil, errors.New("non autorisé")
	}

	if input.Content != "" {
		post.Content = strings.TrimSpace(input.Content)
	}
	if input.Visibility != "" {
		post.Visibility = input.Visibility
	}
	if input.DocumentType != "" {
		post.DocumentType = input.DocumentType
	}

	if err := s.repo.Update(post); err != nil {
		return nil, errors.New("erreur lors de la mise à jour")
	}

	return s.GetPostByID(postID, creatorID)
}

// DeletePost supprime un post
func (s *service) DeletePost(postID, creatorID uint) error {
	post, err := s.repo.GetByID(postID)
	if err != nil {
		return errors.New("post non trouvé")
	}
	if post.CreatorID != creatorID {
		return errors.New("non autorisé")
	}
	return s.repo.Delete(postID)
}
