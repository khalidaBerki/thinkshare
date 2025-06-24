package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"backend/internal/db"
	"backend/internal/user"
)

var testUser user.User

func TestMain(m *testing.M) {
	// Init DB
	db.InitDB()

	// Nettoyer et remigrer la table
	db.GormDB.Migrator().DropTable(&user.User{})
	db.GormDB.AutoMigrate(&user.User{})

	// Créer l'utilisateur de test
	testUser = createTestUser()

	// Lancer les tests
	code := m.Run()
	os.Exit(code)
}

func setupRouter() *gin.Engine {
	r := gin.Default()

	// Middleware simulé pour les tests
	r.Use(func(c *gin.Context) {
		// Utiliser la clé exacte "user_id" comme dans le handler
		c.Set("user_id", 1)
		c.Next()
	})

	r.GET("/api/profile", user.GetProfileHandler)
	r.PUT("/api/profile", user.UpdateProfileHandler)
	return r
}

func createTestUser() user.User {
	// Supprimer l'utilisateur existant s'il existe
	db.GormDB.Unscoped().Where("id = ?", 1).Delete(&user.User{})

	u := user.User{
		ID:           1,
		FullName:     "Integration Test",
		Bio:          "Test bio",
		AvatarURL:    "https://avatar.test/avatar.png",
		Username:     "integration_test",
		Email:        "integration@test.com",
		PasswordHash: "hashed_password",
		Role:         "user",
		CreatedAt:    time.Now(),
	}

	// Créer l'utilisateur avec un ID spécifique
	result := db.GormDB.Create(&u)
	if result.Error != nil {
		panic("Failed to create test user: " + result.Error.Error())
	}

	// Vérifier que l'utilisateur a été créé avec le bon ID
	var createdUser user.User
	if err := db.GormDB.First(&createdUser, 1).Error; err != nil {
		panic("Failed to verify test user creation: " + err.Error())
	}

	return u
}

func TestGetProfileHandler(t *testing.T) {
	r := setupRouter()

	req, _ := http.NewRequest("GET", "/api/profile", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var profile user.ProfileDTO
	err := json.Unmarshal(w.Body.Bytes(), &profile)
	assert.NoError(t, err)

	assert.Equal(t, testUser.FullName, profile.FullName)
	assert.Equal(t, testUser.Bio, profile.Bio)
	assert.Equal(t, testUser.AvatarURL, profile.AvatarURL)
}

func TestUpdateProfileHandler(t *testing.T) {
	r := setupRouter()

	updateInput := user.UpdateUserInput{
		FullName:  "Updated Integration",
		Bio:       "Updated Bio",
		AvatarURL: "https://avatar.test/updated.png",
	}
	body, _ := json.Marshal(updateInput)

	req, _ := http.NewRequest("PUT", "/api/profile", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var updatedUser user.User
	err := db.GormDB.First(&updatedUser, testUser.ID).Error
	assert.NoError(t, err)

	assert.Equal(t, updateInput.FullName, updatedUser.FullName)
	assert.Equal(t, updateInput.Bio, updatedUser.Bio)
	assert.Equal(t, updateInput.AvatarURL, updatedUser.AvatarURL)
}
