package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// ✅ Middleware à appliquer sur les routes protégées
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Récupère le header "Authorization: Bearer <token>"
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header required"})
			c.Abort()
			return
		}

		// Sépare "Bearer" du token réel
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header"})
			c.Abort()
			return
		}

		// Vérifie et décode le token
		userID, err := ParseJWT(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		// ✅ Injecte user_id dans le contexte pour le handler
		c.Set("user_id", userID)
		c.Next()
	}
}
