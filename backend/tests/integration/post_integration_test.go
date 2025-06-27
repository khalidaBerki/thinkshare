package integration

import (
	"backend/internal/post"
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock du repository pour les tests d'intégration
type MockPostRepository struct {
	mock.Mock
}

func (m *MockPostRepository) Create(post *post.Post) error {
	args := m.Called(post)
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
func (m *MockPostRepository) Update(post *post.Post) error {
	args := m.Called(post)
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
func (m *MockPostRepository) GetAllAfter(afterID uint, limit int) ([]*post.Post, error) {
	args := m.Called(afterID, limit)
	return args.Get(0).([]*post.Post), args.Error(1)
}
func (m *MockPostRepository) GetByCreatorAfter(creatorID, afterID uint, limit int) ([]*post.Post, error) {
	args := m.Called(creatorID, afterID, limit)
	return args.Get(0).([]*post.Post), args.Error(1)
}

// --- Setup router ---
func setupPostRouter() (*gin.Engine, *post.Handler, *MockPostRepository) {
	gin.SetMode(gin.TestMode)
	mockRepo := new(MockPostRepository)
	postService := post.NewService(mockRepo)
	postHandler := post.NewHandler(postService)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", 1)
		c.Next()
	})
	r.RedirectTrailingSlash = false
	r.RedirectFixedPath = false
	api := r.Group("/api")
	postHandler.RegisterRoutes(api)
	return r, postHandler, mockRepo
}

// Constantes et variables utiles pour les tests
const (
	testUserID = uint(1)
)

// Fonction d'aide pour créer une réponse de post standard
func createTestPost(id uint) *post.Post {
	return &post.Post{
		ID:         id,
		CreatorID:  testUserID,
		Content:    fmt.Sprintf("Contenu du post de test %d", id),
		Visibility: post.Public,
	}
}

// Test de création d'un post (texte uniquement)
func TestIntegration_CreateTextPost(t *testing.T) {
	// Configurer le router et les mocks
	r, _, mockRepo := setupPostRouter()

	// Configurer le mock pour simuler la création du post
	mockRepo.On("Create", mock.MatchedBy(func(p *post.Post) bool {
		return p.Content == "Contenu du post de test" && p.Visibility == post.Public
	})).Return(nil)

	mockRepo.On("GetPostsWithStats", mock.AnythingOfType("[]*post.Post"), uint(1)).Return([]*post.PostDTO{
		{
			ID:         0,
			CreatorID:  1,
			Content:    "Contenu du post de test",
			Visibility: string(post.Public),
			MediaURLs:  []string{},
		},
	}, nil)

	mockRepo.On("GetByID", uint(0)).Return(&post.Post{
		ID:         0,
		CreatorID:  1,
		Content:    "Contenu du post de test",
		Visibility: post.Public,
	}, nil)

	mockRepo.On("GetCreatorInfo", uint(1)).Return(&post.CreatorInfo{
		ID:       1,
		Username: "testuser",
	}, nil)

	// Créer un buffer pour les données multipart/form-data
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	// Ajouter le contenu du post
	fw, _ := w.CreateFormField("content")
	fw.Write([]byte("Contenu du post de test"))

	// Ajouter la visibilité
	fw, _ = w.CreateFormField("visibility")
	fw.Write([]byte("public"))

	// Fermer le writer
	w.Close()

	// Créer une requête POST avec le type de contenu multipart/form-data
	// Ajouter un slash final à l'URL pour éviter la redirection 307
	req, _ := http.NewRequest("POST", "/api/posts/", &b)
	req.Header.Set("Content-Type", w.FormDataContentType())

	// Journaliser pour débogage
	t.Logf("URL: %s, document_type: %s", req.URL.String(), req.Header.Get("Content-Type"))

	// Exécuter la requête
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	// Journaliser la réponse pour débogage
	t.Logf("Réponse: Code=%d, Body=%s", rec.Code, rec.Body.String())

	// Vérifier le code de statut
	assert.Equal(t, http.StatusCreated, rec.Code)

	// Vérifier que le mock a été appelé comme prévu
	mockRepo.AssertExpectations(t)
}

// Test de récupération d'un post par ID
func TestIntegration_GetPostByID(t *testing.T) {
	// Configurer le router et les mocks
	r, _, mockRepo := setupPostRouter()

	// Créer un post de test à retourner
	testPost := createTestPost(1)

	// Configurer le mock pour simuler la récupération du post
	mockRepo.On("GetByID", uint(1)).Return(testPost, nil)

	// Ajoute ce mock pour GetPostsWithStats
	mockRepo.On("GetPostsWithStats", mock.AnythingOfType("[]*post.Post"), uint(1)).Return(
		[]*post.PostDTO{
			{
				ID:         testPost.ID,
				CreatorID:  testPost.CreatorID,
				Content:    testPost.Content,
				Visibility: string(testPost.Visibility),
				MediaURLs:  []string{},
			},
		}, nil,
	)

	mockRepo.On("GetCreatorInfo", uint(1)).Return(&post.CreatorInfo{
		ID:       1,
		Username: "testuser",
	}, nil)

	// Créer une requête GET pour récupérer le post
	req, _ := http.NewRequest("GET", "/api/posts/1", nil)

	// Exécuter la requête
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Vérifier le code de statut
	assert.Equal(t, http.StatusOK, w.Code)

	// Vérifier la réponse
	var response post.Post
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, testPost.ID, response.ID)
	assert.Equal(t, testPost.Content, response.Content)
	assert.Equal(t, testPost.Visibility, response.Visibility)

	// Vérifier que le mock a été appelé comme prévu
	mockRepo.AssertExpectations(t)
}

// Test de récupération de tous les posts
func TestIntegration_GetAllPosts(t *testing.T) {
	// Configurer le router et les mocks
	r, _, mockRepo := setupPostRouter()

	// Créer plusieurs posts de test à retourner
	testPosts := []*post.Post{}
	for i := 1; i <= 5; i++ {
		testPosts = append(testPosts, createTestPost(uint(i)))
	}

	// Créer les DTO correspondants
	testDTOs := []*post.PostDTO{}
	for _, p := range testPosts {
		testDTOs = append(testDTOs, &post.PostDTO{
			ID:         p.ID,
			CreatorID:  p.CreatorID,
			Content:    p.Content,
			Visibility: string(p.Visibility),
			CreatedAt:  p.CreatedAt,
			UpdatedAt:  p.UpdatedAt,
			MediaURLs:  []string{},
		})
	}

	// Mock scroll infini
	mockRepo.On("GetAllAfter", uint(0), mock.AnythingOfType("int")).Return(testPosts, nil)
	mockRepo.On("GetPostsWithStats", testPosts, uint(1)).Return(testDTOs, nil)

	// Créer une requête GET pour récupérer tous les posts
	req, _ := http.NewRequest("GET", "/api/posts?limit=10", nil)

	// Exécuter la requête
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Vérifier le code de statut et afficher la réponse pour débogage
	t.Logf("Réponse GetAllPosts: %s", w.Body.String())
	assert.Equal(t, http.StatusOK, w.Code)

	// Vérifier la réponse
	var response struct {
		Posts []*post.PostDTO `json:"posts"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response.Posts, 5) // Nous avons créé 5 posts

	// Vérifier que le mock a été appelé comme prévu
	mockRepo.AssertExpectations(t)
}

// Test de mise à jour d'un post
func TestIntegration_UpdatePost(t *testing.T) {
	// Configurer le router et les mocks
	r, _, mockRepo := setupPostRouter()

	// Créer un post de test à retourner lors de la vérification
	testPost := createTestPost(1)
	testPost.Content = "Contenu original"
	testPost.Visibility = post.Private

	// Configurer le mock pour simuler la récupération et la mise à jour du post
	mockRepo.On("GetByID", uint(1)).Return(testPost, nil)
	mockRepo.On("Update", mock.MatchedBy(func(p *post.Post) bool {
		return p.ID == 1 && p.Content == "Contenu mis à jour" && p.Visibility == post.Public
	})).Return(nil)

	mockRepo.On("GetPostsWithStats", []*post.Post{testPost}, uint(1)).Return([]*post.PostDTO{
		{
			ID:         testPost.ID,
			CreatorID:  testPost.CreatorID,
			Content:    "Contenu mis à jour",
			Visibility: string(post.Public),
			CreatedAt:  testPost.CreatedAt,
			UpdatedAt:  testPost.UpdatedAt,
			MediaURLs:  []string{},
		},
	}, nil)

	mockRepo.On("GetCreatorInfo", uint(1)).Return(&post.CreatorInfo{
		ID:       1,
		Username: "testuser",
	}, nil)

	// Préparer les données pour la mise à jour
	updateData := post.UpdatePostInput{
		Content:    "Contenu mis à jour",
		Visibility: post.Public,
	}
	jsonData, _ := json.Marshal(updateData)

	// Créer une requête PUT
	req, _ := http.NewRequest("PUT", "/api/posts/1", bytes.NewBuffer(jsonData))
	req.Header.Set("document_type", "application/json")

	// Exécuter la requête
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Vérifier le code de statut
	assert.Equal(t, http.StatusOK, w.Code)

	// Vérifier que le mock a été appelé comme prévu
	mockRepo.AssertExpectations(t)
}

// Test de suppression d'un post
func TestIntegration_DeletePost(t *testing.T) {
	// Configurer le router et les mocks
	r, _, mockRepo := setupPostRouter()

	// Créer un post de test à retourner lors de la vérification
	testPost := createTestPost(1)

	// Configurer le mock pour simuler la récupération et la suppression du post
	mockRepo.On("GetByID", uint(1)).Return(testPost, nil)
	mockRepo.On("Delete", uint(1)).Return(nil)

	// Créer une requête DELETE
	req, _ := http.NewRequest("DELETE", "/api/posts/1", nil)

	// Exécuter la requête
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Vérifier le code de statut
	assert.Equal(t, http.StatusNoContent, w.Code)

	// Vérifier que le mock a été appelé comme prévu
	mockRepo.AssertExpectations(t)
}

// Test de récupération des statistiques médias
func TestIntegration_GetMediaStats(t *testing.T) {
	// Configurer le router et les mocks
	r, _, mockRepo := setupPostRouter()

	// Configurer le mock pour simuler la récupération des stats
	mockRepo.On("CountMediaByType", "image").Return(int64(10), nil)
	mockRepo.On("CountMediaByType", "video").Return(int64(5), nil)
	mockRepo.On("CountMediaByType", "document").Return(int64(3), nil)

	// Créer une requête GET pour les statistiques
	req, _ := http.NewRequest("GET", "/api/posts/media/stats", nil)

	// Exécuter la requête
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Vérifier le code de statut
	assert.Equal(t, http.StatusOK, w.Code)

	// Journaliser la réponse pour débogage
	t.Logf("Réponse brute: %s", w.Body.String())

	// Vérifier la réponse
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Vérifier la structure de base de la réponse
	assert.Contains(t, response, "statistics", "La réponse devrait contenir un champ 'statistics'")
	assert.Contains(t, response, "recommendations", "La réponse devrait contenir un champ 'recommendations'")

	// Extraire et vérifier les statistiques
	statistics, ok := response["statistics"].(map[string]interface{})
	assert.True(t, ok, "Le champ 'statistics' devrait être un objet")

	// Vérifier le total des médias
	assert.Contains(t, statistics, "total_media", "Les statistiques devraient contenir un champ 'total_media'")
	totalMedia, ok := statistics["total_media"].(float64)
	assert.True(t, ok, "Le champ 'total_media' devrait être un nombre")
	assert.Equal(t, float64(18), totalMedia, "Le nombre total de médias devrait être 18")

	// Vérifier les statistiques par type
	assert.Contains(t, statistics, "media_by_type", "Les statistiques devraient contenir un champ 'media_by_type'")
	mediaByType, ok := statistics["media_by_type"].(map[string]interface{})
	assert.True(t, ok, "Le champ 'media_by_type' devrait être un objet")

	// Vérifier chaque type de média
	assert.Contains(t, mediaByType, "image", "Les stats devraient inclure le type 'image'")
	assert.Contains(t, mediaByType, "video", "Les stats devraient inclure le type 'video'")
	assert.Contains(t, mediaByType, "document", "Les stats devraient inclure le type 'document'")

	// Vérifier les valeurs pour chaque type
	imageCount, ok := mediaByType["image"].(float64)
	assert.True(t, ok, "Le compteur d'images devrait être un nombre")
	assert.Equal(t, float64(10), imageCount, "Le nombre d'images devrait être 10")

	videoCount, ok := mediaByType["video"].(float64)
	assert.True(t, ok, "Le compteur de vidéos devrait être un nombre")
	assert.Equal(t, float64(5), videoCount, "Le nombre de vidéos devrait être 5")

	documentCount, ok := mediaByType["document"].(float64)
	assert.True(t, ok, "Le compteur de documents devrait être un nombre")
	assert.Equal(t, float64(3), documentCount, "Le nombre de documents devrait être 3")

	// Vérifier que le mock a été appelé comme prévu
	mockRepo.AssertExpectations(t)
}

// Les tests d'intégration utilisent déjà le TestMain défini dans user_integration_test.go

// Mocks pour les tests
func checkUploadsDirectorySecurity() error {
	// Simuler une vérification réussie
	return nil
}

func getRecommendedFormats() map[string]interface{} {
	return map[string]interface{}{
		"image": map[string]string{
			"preferred": ".jpg, .png",
		},
		"video": map[string]string{
			"preferred": ".mp4",
		},
		"document": map[string]string{
			"preferred": ".pdf",
		},
	}
}

func TestGetAllPostsAfter_Success(t *testing.T) {
	mockRepo := new(MockPostRepository)
	service := post.NewService(mockRepo)

	posts := []*post.Post{
		{ID: 21, CreatorID: 1, Content: "suite 1", Visibility: post.Public},
		{ID: 22, CreatorID: 1, Content: "suite 2", Visibility: post.Public},
	}
	mockRepo.On("GetAllAfter", uint(20), 2).Return(posts, nil)
	mockRepo.On("GetPostsWithStats", posts, uint(1)).Return([]*post.PostDTO{
		{ID: 21, CreatorID: 1, Content: "suite 1", Visibility: string(post.Public)},
		{ID: 22, CreatorID: 1, Content: "suite 2", Visibility: string(post.Public)},
	}, nil)

	result, err := service.GetAllPostsAfter(20, 2, 1)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	mockRepo.AssertExpectations(t)
}

// Test de récupération de tous les posts après un certain ID
// --- Test scroll infini /api/posts ---
func TestIntegration_GetAllPostsAfter(t *testing.T) {
	r, _, mockRepo := setupPostRouter()

	// Simule des posts à retourner après l'ID 20
	posts := []*post.Post{
		{ID: 21, CreatorID: 1, Content: "suite 1", Visibility: post.Public},
		{ID: 22, CreatorID: 1, Content: "suite 2", Visibility: post.Public},
	}
	mockRepo.On("GetAllAfter", uint(20), 2).Return(posts, nil)
	mockRepo.On("GetPostsWithStats", posts, uint(1)).Return([]*post.PostDTO{
		{ID: 21, CreatorID: 1, Content: "suite 1", Visibility: string(post.Public)},
		{ID: 22, CreatorID: 1, Content: "suite 2", Visibility: string(post.Public)},
	}, nil)

	// Requête GET avec after=20&limit=2
	req, _ := http.NewRequest("GET", "/api/posts?after=20&limit=2", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response struct {
		Posts   []*post.PostDTO `json:"posts"`
		HasMore bool            `json:"has_more"`
		LastID  uint            `json:"last_id"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response.Posts, 2)
	assert.Equal(t, uint(22), response.LastID)
	mockRepo.AssertExpectations(t)
}
