package unit

import (
	"errors"
	"testing"
	"time"

	"backend/internal/media"
	"backend/internal/post"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock du repository pour les tests unitaires
type MockPostRepository struct {
	mock.Mock
}

func (m *MockPostRepository) Create(post *post.Post) error {
	args := m.Called(post)
	return args.Error(0)
}

func (m *MockPostRepository) GetAll(page, pageSize int, visibility post.Visibility) ([]post.Post, error) {
	args := m.Called(page, pageSize, visibility)
	return args.Get(0).([]post.Post), args.Error(1)
}

func (m *MockPostRepository) GetByID(id uint) (*post.Post, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*post.Post), args.Error(1)
}

func (m *MockPostRepository) Update(post *post.Post) error {
	args := m.Called(post)
	return args.Error(0)
}

func (m *MockPostRepository) Delete(id uint, userID uint) error {
	args := m.Called(id, userID)
	return args.Error(0)
}

func (m *MockPostRepository) CountMediaByType(mediaType string) (int64, error) {
	args := m.Called(mediaType)
	return args.Get(0).(int64), args.Error(1)
}

// Tests unitaires pour le service post
func TestCreatePost_Success(t *testing.T) {
	mockRepo := new(MockPostRepository)
	service := post.NewService(mockRepo)

	// Créer un post valide avec contenu textuel et sans média
	testPost := &post.Post{
		CreatorID:  1,
		Content:    "Test post content",
		Visibility: post.Public,
		CreatedAt:  time.Now(),
	}

	// Configurer le mock pour retourner nil (succès)
	mockRepo.On("Create", testPost).Return(nil)

	// Appeler la méthode et vérifier le résultat
	err := service.CreatePost(testPost)

	// Assertions
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestCreatePost_EmptyContentNoMedia(t *testing.T) {
	mockRepo := new(MockPostRepository)
	service := post.NewService(mockRepo)

	// Créer un post sans contenu ni média (invalide)
	testPost := &post.Post{
		CreatorID:  1,
		Content:    "", // Contenu vide
		Visibility: post.Public,
		CreatedAt:  time.Now(),
		Media:      []media.Media{}, // Pas de média
	}

	// Appeler la méthode
	err := service.CreatePost(testPost)

	// Assertions
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "must have either content or media")
	// Le mock ne devrait pas être appelé car la validation échoue avant
	mockRepo.AssertNotCalled(t, "Create")
}

func TestCreatePost_InvalidVisibility(t *testing.T) {
	mockRepo := new(MockPostRepository)
	service := post.NewService(mockRepo)

	// Créer un post avec une visibilité invalide
	testPost := &post.Post{
		CreatorID:  1,
		Content:    "Test post content",
		Visibility: "invalid_visibility", // Visibilité invalide
		CreatedAt:  time.Now(),
	}

	// Appeler la méthode
	err := service.CreatePost(testPost)

	// Assertions
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid visibility")
	// Le mock ne devrait pas être appelé car la validation échoue avant
	mockRepo.AssertNotCalled(t, "Create")
}

func TestGetPostByID_Success(t *testing.T) {
	mockRepo := new(MockPostRepository)
	service := post.NewService(mockRepo)

	// Post attendu à retourner
	expectedPost := &post.Post{
		ID:         1,
		CreatorID:  2,
		Content:    "Test post content",
		Visibility: post.Public,
		CreatedAt:  time.Now(),
	}

	// Configurer le mock pour retourner le post
	mockRepo.On("GetByID", uint(1)).Return(expectedPost, nil)

	// Appeler la méthode
	result, err := service.GetPostByID(1)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, expectedPost.ID, result.ID)
	assert.Equal(t, expectedPost.Content, result.Content)
	mockRepo.AssertExpectations(t)
}

func TestGetPostByID_NotFound(t *testing.T) {
	mockRepo := new(MockPostRepository)
	service := post.NewService(mockRepo)

	// Configurer le mock pour simuler un post non trouvé
	mockRepo.On("GetByID", uint(999)).Return(nil, errors.New("record not found"))

	// Appeler la méthode
	result, err := service.GetPostByID(999)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "not found")
	mockRepo.AssertExpectations(t)
}

func TestUpdatePost_Success(t *testing.T) {
	mockRepo := new(MockPostRepository)
	service := post.NewService(mockRepo)

	// Post existant en base de données
	existingPost := &post.Post{
		ID:         1,
		CreatorID:  2,
		Content:    "Original content",
		Visibility: post.Private,
	}

	// Configurer le mock pour retourner le post existant
	mockRepo.On("GetByID", uint(1)).Return(existingPost, nil)

	// Préparation des données pour la mise à jour
	// (pas besoin de créer un nouvel objet post car on utilise l'input directement)

	// Configurer le mock pour accepter la mise à jour
	mockRepo.On("Update", mock.MatchedBy(func(p *post.Post) bool {
		return p.ID == 1 && p.Content == "Updated content" && p.Visibility == post.Public
	})).Return(nil)

	// Préparer les données d'entrée pour la mise à jour
	updateInput := post.UpdatePostInput{
		Content:    "Updated content",
		Visibility: post.Public,
	}

	// Appeler la méthode de mise à jour
	err := service.UpdatePost(1, 2, updateInput)

	// Assertions
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUpdatePost_Unauthorized(t *testing.T) {
	mockRepo := new(MockPostRepository)
	service := post.NewService(mockRepo)

	// Post existant en base de données
	existingPost := &post.Post{
		ID:         1,
		CreatorID:  2, // Créateur ID 2
		Content:    "Original content",
		Visibility: post.Private,
	}

	// Configurer le mock pour retourner le post existant
	mockRepo.On("GetByID", uint(1)).Return(existingPost, nil)

	// Préparer les données d'entrée pour la mise à jour
	updateInput := post.UpdatePostInput{
		Content:    "Updated content",
		Visibility: post.Public,
	}

	// Utilisateur 3 (différent du créateur) tente de mettre à jour
	err := service.UpdatePost(1, 3, updateInput)

	// Assertions
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unauthorized")
	// Vérifier que Update n'a pas été appelé
	mockRepo.AssertNotCalled(t, "Update")
}

func TestDeletePost_Success(t *testing.T) {
	mockRepo := new(MockPostRepository)
	service := post.NewService(mockRepo)

	// Post existant en base de données
	existingPost := &post.Post{
		ID:         1,
		CreatorID:  2,
		Content:    "Post to delete",
		Visibility: post.Public,
	}

	// Configurer le mock pour retourner le post existant
	mockRepo.On("GetByID", uint(1)).Return(existingPost, nil)

	// Configurer le mock pour accepter la suppression
	mockRepo.On("Delete", uint(1), uint(2)).Return(nil)

	// Appeler la méthode de suppression
	err := service.DeletePost(1, 2)

	// Assertions
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDeletePost_Unauthorized(t *testing.T) {
	mockRepo := new(MockPostRepository)
	service := post.NewService(mockRepo)

	// Post existant en base de données
	existingPost := &post.Post{
		ID:         1,
		CreatorID:  2, // Créateur ID 2
		Content:    "Post to delete",
		Visibility: post.Public,
	}

	// Configurer le mock pour retourner le post existant
	mockRepo.On("GetByID", uint(1)).Return(existingPost, nil)

	// Utilisateur 3 (différent du créateur) tente de supprimer
	err := service.DeletePost(1, 3)

	// Assertions
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unauthorized")
	// Vérifier que Delete n'a pas été appelé
	mockRepo.AssertNotCalled(t, "Delete")
}
