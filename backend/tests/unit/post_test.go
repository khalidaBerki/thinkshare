package unit

import (
	"errors"
	"testing"

	"backend/internal/media"
	"backend/internal/post"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- Mock Repository ---

type MockPostRepository struct {
	mock.Mock
}

func (m *MockPostRepository) Create(p *post.Post) error {
	args := m.Called(p)
	return args.Error(0)
}

func (m *MockPostRepository) GetByID(id uint) (*post.Post, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*post.Post), args.Error(1)
}

func (m *MockPostRepository) GetAll(page, limit int) ([]*post.Post, int64, error) {
	args := m.Called(page, limit)
	return args.Get(0).([]*post.Post), args.Get(1).(int64), args.Error(2)
}

func (m *MockPostRepository) GetByCreatorID(creatorID uint, page, limit int) ([]*post.Post, int64, error) {
	args := m.Called(creatorID, page, limit)
	return args.Get(0).([]*post.Post), args.Get(1).(int64), args.Error(2)
}

func (m *MockPostRepository) Update(p *post.Post) error {
	args := m.Called(p)
	return args.Error(0)
}

func (m *MockPostRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockPostRepository) GetPostStats(postID, userID uint) (*post.PostStats, error) {
	args := m.Called(postID, userID)
	return args.Get(0).(*post.PostStats), args.Error(1)
}

func (m *MockPostRepository) GetPostsWithStats(posts []*post.Post, userID uint) ([]*post.PostDTO, error) {
	args := m.Called(posts, userID)
	return args.Get(0).([]*post.PostDTO), args.Error(1)
}

func (m *MockPostRepository) GetCreatorInfo(userID uint) (*post.CreatorInfo, error) {
	args := m.Called(userID)
	return args.Get(0).(*post.CreatorInfo), args.Error(1)
}

func (m *MockPostRepository) CountMediaByType(mediaType string) (int64, error) {
	args := m.Called(mediaType)
	return args.Get(0).(int64), args.Error(1)
}

// --- Tests ---

func TestCreatePost_Success(t *testing.T) {
	mockRepo := new(MockPostRepository)
	service := post.NewService(mockRepo)

	input := post.CreatePostInput{
		Content:    "Test post",
		Visibility: post.Public,
	}

	createdPost := &post.Post{
		ID:         1,
		CreatorID:  1,
		Content:    input.Content,
		Visibility: input.Visibility,
	}

	mockRepo.On("Create", mock.AnythingOfType("*post.Post")).Return(nil)
	mockRepo.On("GetByID", mock.AnythingOfType("uint")).Return(createdPost, nil)
	mockRepo.On("GetPostsWithStats", []*post.Post{createdPost}, uint(1)).
		Return([]*post.PostDTO{
			{ID: 1, CreatorID: 1, Content: input.Content, Visibility: string(input.Visibility)},
		}, nil)
	mockRepo.On("GetCreatorInfo", uint(1)).
		Return(&post.CreatorInfo{ID: 1, Username: "testuser"}, nil)

	dto, err := service.CreatePost(1, input)

	assert.NoError(t, err)
	assert.NotNil(t, dto)
	assert.Equal(t, input.Content, dto.Content)
	mockRepo.AssertExpectations(t)
}

func TestCreatePost_EmptyContent(t *testing.T) {
	mockRepo := new(MockPostRepository)
	service := post.NewService(mockRepo)

	input := post.CreatePostInput{
		Content:    "",
		Visibility: post.Public,
		Media:      []media.Media{},
	}

	_, err := service.CreatePost(1, input)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "vide")
	mockRepo.AssertNotCalled(t, "Create")
}

func TestCreatePost_InvalidVisibility(t *testing.T) {
	mockRepo := new(MockPostRepository)
	service := post.NewService(mockRepo)

	input := post.CreatePostInput{
		Content:    "Contenu",
		Visibility: "secret",
	}

	_, err := service.CreatePost(1, input)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid visibility")
	mockRepo.AssertNotCalled(t, "Create")
}

func TestGetPostByID_Success(t *testing.T) {
	mockRepo := new(MockPostRepository)
	service := post.NewService(mockRepo)

	existing := &post.Post{
		ID:         1,
		CreatorID:  2,
		Content:    "Un post",
		Visibility: post.Public,
	}

	expectedDTO := &post.PostDTO{
		ID:         1,
		CreatorID:  2,
		Content:    "Un post",
		Visibility: string(post.Public),
	}

	mockRepo.On("GetByID", uint(1)).Return(existing, nil)
	mockRepo.On("GetPostsWithStats", []*post.Post{existing}, uint(0)).Return([]*post.PostDTO{expectedDTO}, nil)
	mockRepo.On("GetCreatorInfo", uint(2)).Return(&post.CreatorInfo{ID: 2, Username: "auteur"}, nil)

	dto, err := service.GetPostByID(1, 0)

	assert.NoError(t, err)
	assert.Equal(t, expectedDTO.Content, dto.Content)
	mockRepo.AssertExpectations(t)
}

func TestGetPostByID_NotFound(t *testing.T) {
	mockRepo := new(MockPostRepository)
	service := post.NewService(mockRepo)

	mockRepo.On("GetByID", uint(404)).Return(nil, errors.New("post non trouvé"))

	// Correction ici : ajoute un second argument (userID, par exemple 0)
	dto, err := service.GetPostByID(404, 0)

	assert.Error(t, err)
	assert.Nil(t, dto)
	assert.Contains(t, err.Error(), "post non trouvé")
}

func TestUpdatePost_Success(t *testing.T) {
	mockRepo := new(MockPostRepository)
	service := post.NewService(mockRepo)

	existing := &post.Post{
		ID:         1,
		CreatorID:  2,
		Content:    "Old",
		Visibility: post.Private,
	}

	mockRepo.On("GetByID", uint(1)).Return(existing, nil)
	mockRepo.On("Update", mock.MatchedBy(func(p *post.Post) bool {
		return p.Content == "New content" && p.Visibility == post.Public
	})).Return(nil)
	mockRepo.On("GetPostsWithStats", mock.Anything, mock.Anything).Return([]*post.PostDTO{
		{
			ID:         1,
			Content:    "New content",
			Visibility: string(post.Public),
		},
	}, nil)
	mockRepo.On("GetCreatorInfo", uint(2)).Return(&post.CreatorInfo{ID: 2, Username: "auteur"}, nil) // <-- AJOUTE CETTE LIGNE

	_, err := service.UpdatePost(1, 2, post.UpdatePostInput{
		Content:    "New content",
		Visibility: post.Public,
	})

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUpdatePost_Unauthorized(t *testing.T) {
	mockRepo := new(MockPostRepository)
	service := post.NewService(mockRepo)

	postToEdit := &post.Post{
		ID:        1,
		CreatorID: 2,
		Content:   "Texte",
	}

	mockRepo.On("GetByID", uint(1)).Return(postToEdit, nil)

	// Correction ici : capture les deux valeurs de retour
	_, err := service.UpdatePost(1, 99, post.UpdatePostInput{
		Content:    "Trop tard",
		Visibility: post.Private,
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "non autorisé")
	mockRepo.AssertNotCalled(t, "Update")
}

func TestDeletePost_Success(t *testing.T) {
	mockRepo := new(MockPostRepository)
	service := post.NewService(mockRepo)

	postToDelete := &post.Post{ID: 1, CreatorID: 2}

	mockRepo.On("GetByID", uint(1)).Return(postToDelete, nil)
	mockRepo.On("Delete", uint(1)).Return(nil)

	err := service.DeletePost(1, 2)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDeletePost_Unauthorized(t *testing.T) {
	mockRepo := new(MockPostRepository)
	service := post.NewService(mockRepo)

	postToDelete := &post.Post{ID: 1, CreatorID: 2}

	mockRepo.On("GetByID", uint(1)).Return(postToDelete, nil)

	err := service.DeletePost(1, 99)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "non autorisé")
	mockRepo.AssertNotCalled(t, "Delete")
}
