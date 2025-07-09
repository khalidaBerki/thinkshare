package unit

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"backend/internal/user"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) GetProfile(id uint) (*user.ProfileDTO, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.ProfileDTO), args.Error(1)
}

func (m *MockUserService) UpdateProfile(id uint, input user.UpdateUserInput) error {
	args := m.Called(id, input)
	return args.Error(0)
}

func setupRouter(mockSvc *MockUserService) *gin.Engine {
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Set("userID", uint(1))
		c.Next()
	})

	r.GET("/api/profile", func(c *gin.Context) {
		userID := c.GetUint("userID")
		profile, err := mockSvc.GetProfile(userID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, profile)
	})

	r.PUT("/api/profile", func(c *gin.Context) {
		userID := c.GetUint("userID")
		var input user.UpdateUserInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Entrée invalide"})
			return
		}
		if err := mockSvc.UpdateProfile(userID, input); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Profil mis à jour avec succès"})
	})

	return r
}

func TestGetProfileHandler_Success(t *testing.T) {
	mockSvc := new(MockUserService)
	expectedProfile := &user.ProfileDTO{
		ID:        1,
		FullName:  "Test User",
		Bio:       "Bio test",
		AvatarURL: "https://avatar.test/image.png",
	}
	mockSvc.On("GetProfile", uint(1)).Return(expectedProfile, nil)

	r := setupRouter(mockSvc)
	req, _ := http.NewRequest("GET", "/api/profile", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var profile user.ProfileDTO
	json.Unmarshal(w.Body.Bytes(), &profile)
	assert.Equal(t, expectedProfile.FullName, profile.FullName)
	mockSvc.AssertExpectations(t)
}

func TestGetProfileHandler_NotFound(t *testing.T) {
	mockSvc := new(MockUserService)
	mockSvc.On("GetProfile", uint(1)).Return(nil, errors.New("user not found"))

	r := setupRouter(mockSvc)
	req, _ := http.NewRequest("GET", "/api/profile", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestUpdateProfileHandler_Success(t *testing.T) {
	mockSvc := new(MockUserService)
	mockSvc.On("UpdateProfile", uint(1), mock.AnythingOfType("user.UpdateUserInput")).Return(nil)

	r := setupRouter(mockSvc)
	updateInput := user.UpdateUserInput{
		FullName:  "Updated Name",
		Bio:       "Updated Bio",
		AvatarURL: "https://avatar.test/updated.png",
	}
	jsonValue, _ := json.Marshal(updateInput)
	req, _ := http.NewRequest("PUT", "/api/profile", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestUpdateProfileHandler_BadRequest(t *testing.T) {
	mockSvc := new(MockUserService)

	r := setupRouter(mockSvc)
	req, _ := http.NewRequest("PUT", "/api/profile", bytes.NewBuffer([]byte(`invalid-json`)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateProfileHandler_Error(t *testing.T) {
	mockSvc := new(MockUserService)
	mockSvc.On("UpdateProfile", uint(1), mock.AnythingOfType("user.UpdateUserInput")).Return(errors.New("update failed"))

	r := setupRouter(mockSvc)
	updateInput := user.UpdateUserInput{
		FullName:  "Updated Name",
		Bio:       "Updated Bio",
		AvatarURL: "https://avatar.test/updated.png",
	}
	jsonValue, _ := json.Marshal(updateInput)
	req, _ := http.NewRequest("PUT", "/api/profile", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockSvc.AssertExpectations(t)
}
