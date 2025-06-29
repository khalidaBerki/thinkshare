package user

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetProfileHandler godoc
// @Summary      Récupérer le profil utilisateur
// @Description  Retourne les informations du profil de l'utilisateur connecté
// @Tags         user
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object} ProfileDTO
// @Failure      404  {object} map[string]string
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
// @Summary      Modifier le profil utilisateur
// @Description  Met à jour les champs du profil (nom, bio, avatar)
// @Tags         user
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body  body  user.UpdateUserInput  true  "Champs modifiables du profil"
// @Success      200  {object} map[string]string
// @Failure      400  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /api/profile [put]
func UpdateProfileHandler(c *gin.Context) {
	var input UpdateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Entrée invalide"})
		return
	}

	// Récupère l'ID depuis le contexte (JWT)
	userID := c.GetInt("user_id")

	if err := UpdateProfile(uint(userID), input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profil mis à jour avec succès"})
}

// GetUserProfileHandler retourne le profil public d'un utilisateur par son ID
func GetUserProfileHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID utilisateur invalide"})
		return
	}

	user, err := GetProfile(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Utilisateur non trouvé"})
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
