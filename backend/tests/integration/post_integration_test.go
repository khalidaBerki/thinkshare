package integration

import (
	"backend/internal/post"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Mock du repository pour les tests d'intégration
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

// Configuration du router pour les tests d'intégration
func setupPostRouter() (*gin.Engine, *post.Handler, *MockPostRepository) {
	// Configuration du mode de test pour Gin
	gin.SetMode(gin.TestMode)

	// Initialisation des mocks
	mockRepo := new(MockPostRepository)
	postService := post.NewService(mockRepo)
	postHandler := post.NewHandler(postService)

	// Création du router
	r := gin.New() // Utiliser gin.New() au lieu de gin.Default() pour éviter les logs

	// Middleware pour simuler un utilisateur authentifié
	r.Use(func(c *gin.Context) {
		c.Set("user_id", 1) // Utilisateur de test avec ID 1
		c.Next()
	})

	// Désactiver les redirections pour les tests
	r.RedirectTrailingSlash = false
	// Désactiver la correction automatique des URLs
	r.RedirectFixedPath = false

	// Enregistrement des routes
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
	t.Logf("URL: %s, Content-Type: %s", req.URL.String(), req.Header.Get("Content-Type"))

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
	testPosts := []post.Post{}
	for i := 1; i <= 5; i++ {
		testPosts = append(testPosts, *createTestPost(uint(i)))
	}

	// Configurer le mock pour simuler la récupération des posts
	mockRepo.On("GetAll", 1, 10, post.Visibility("")).Return(testPosts, nil)

	// Créer une requête GET pour récupérer tous les posts
	req, _ := http.NewRequest("GET", "/api/posts/?page=1&pageSize=10", nil)

	// Exécuter la requête
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Vérifier le code de statut et afficher la réponse pour débogage
	t.Logf("Réponse GetAllPosts: %s", w.Body.String())
	assert.Equal(t, http.StatusOK, w.Code)

	// Vérifier la réponse
	var response []post.Post
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response, 5) // Nous avons créé 5 posts

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

	// Préparer les données pour la mise à jour
	updateData := post.UpdatePostInput{
		Content:    "Contenu mis à jour",
		Visibility: post.Public,
	}
	jsonData, _ := json.Marshal(updateData)

	// Créer une requête PUT
	req, _ := http.NewRequest("PUT", "/api/posts/1", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

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
	mockRepo.On("Delete", uint(1), uint(1)).Return(nil)

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
