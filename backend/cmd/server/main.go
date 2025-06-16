package main

import (
	"github.com/gin-gonic/gin"
	"github.com/khalidaBerki/thinkshare/backend/internal/auth"
	"github.com/khalidaBerki/thinkshare/backend/internal/user"
)

func main() {
	r := gin.Default()

	api := r.Group("/api")
	{
		// ✅ Route temporaire pour générer un token JWT simulé (utile en dev avant d'avoir la vraie connexion)
		api.POST("/fake-login", func(c *gin.Context) {
			// Simule un utilisateur avec ID = 1
			token, err := auth.GenerateJWT(1)
			if err != nil {
				c.JSON(500, gin.H{"error": "could not generate token"})
				return
			}
			c.JSON(200, gin.H{"token": token})
		})

		// 🔐 Routes protégées
		api.GET("/profile", auth.AuthMiddleware(), user.GetProfileHandler)
		api.PUT("/profile", auth.AuthMiddleware(), user.UpdateProfileHandler)

		// 📌 En prod, on remplacera "/fake-login" par une vraie route /login et /register
	}

	r.Run(":8080") // lance le serveur en local
}
