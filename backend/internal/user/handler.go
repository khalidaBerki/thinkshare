package user

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetProfileHandler(c *gin.Context) {
	// Pour le test, on fixe un ID = 1
	user, err := GetProfile(1)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func UpdateProfileHandler(c *gin.Context) {
	var input UpdateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Entrée invalide"})
		return
	}

	// ID = 1 (mocké)
	if err := UpdateProfile(1, input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profil mis à jour avec succès"})
}
