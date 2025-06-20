package post

import (
	"errors"
	"strings"
)

type Service interface {
	CreatePost(p *Post) error
	GetAllPosts(page, pageSize int, visibility Visibility) ([]Post, error)
	UpdatePost(postID uint, userID uint, input UpdatePostInput) error
	DeletePost(postID uint, userID uint) error
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	if r == nil {
		panic("repository cannot be nil")
	}
	return &service{repo: r}
}

func (s *service) CreatePost(p *Post) error {
	if p.Visibility != Public && p.Visibility != Private {
		return errors.New("invalid visibility")
	}
	if strings.TrimSpace(p.Content) == "" && len(p.Media) == 0 {
		return errors.New("post must have either content or media")
	}
	return s.repo.Create(p)
}

func (s *service) GetAllPosts(page, pageSize int, visibility Visibility) ([]Post, error) {
	return s.repo.GetAll(page, pageSize, visibility)
}

func (s *service) UpdatePost(postID uint, userID uint, input UpdatePostInput) error {
	post, err := s.repo.GetByID(postID)
	if err != nil {
		return err
	}
	if post.CreatorID != userID {
		return errors.New("unauthorized")
	}

	post.Content = input.Content
	post.Visibility = Visibility(input.Visibility)

	// Ici tu peux envisager de mettre à jour les médias aussi,
	// mais ce sera plus complexe, pour l'instant on gère juste le contenu

	return s.repo.Update(post)
}

func (s *service) DeletePost(postID uint, userID uint) error {
	post, err := s.repo.GetByID(postID)
	if err != nil {
		return err
	}
	if post.CreatorID != userID {
		return errors.New("unauthorized")
	}
	return s.repo.Delete(postID, userID)
}
