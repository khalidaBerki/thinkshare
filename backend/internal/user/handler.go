package user

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetProfileHandler godoc
// @Summary      Get current user profile
// @Description  Returns the profile information of the authenticated user
// @Tags         user
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object} ProfileDTO
// @Failure      401  {object} map[string]string "Unauthorized"
// @Failure      404  {object} map[string]string "User not found"
// @Router       /api/profile [get]
func GetProfileHandler(c *gin.Context) {
	// Récupère l'ID depuis le contexte (JWT)
	userID := c.GetInt("user_id")

	user, err := GetProfile(uint(userID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Conversion en DTO (sans données sensibles)
	profile := ProfileDTO{
		ID:        user.ID,
		FullName:  user.FullName,
		Bio:       user.Bio,
		AvatarURL: user.AvatarURL,
	}

	c.JSON(http.StatusOK, profile)
}

// UpdateProfileHandler godoc
// @Summary      Update current user profile
// @Description  Update profile fields (full name, bio, avatar)
// @Tags         user
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body  body  user.UpdateUserInput  true  "Updatable profile fields"
// @Success      200  {object} map[string]string "Profile updated successfully"
// @Failure      400  {object} map[string]string "Invalid input"
// @Failure      401  {object} map[string]string "Unauthorized"
// @Failure      500  {object} map[string]string "Internal server error"
// @Router       /api/profile [put]
func UpdateProfileHandler(c *gin.Context) {
	var input UpdateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid entry"})
		return
	}

	// Récupère l'ID depuis le contexte (JWT)
	userID := c.GetInt("user_id")

	if err := UpdateProfile(uint(userID), input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
}

// GetUserProfileHandler godoc
// @Summary      Get public user profile
// @Description  Returns the public profile of a user by their ID
// @Tags         user
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200  {object} ProfileDTO
// @Failure      400  {object} map[string]string "Invalid user ID"
// @Failure      404  {object} map[string]string "User not found"
// @Router       /api/users/{id}/profile [get]
func GetUserProfileHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := GetProfile(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	profile := ProfileDTO{
		ID:        user.ID,
		FullName:  user.FullName,
		Bio:       user.Bio,
		AvatarURL: user.AvatarURL,
	}

	c.JSON(http.StatusOK, profile)
}
