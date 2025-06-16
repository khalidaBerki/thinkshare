package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/khalidaBerki/thinkshare/backend/internal/auth"
	"github.com/khalidaBerki/thinkshare/backend/internal/user"
)

func main() {
	// ✅ En dev, charger le fichier .env localement
	if os.Getenv("GIN_MODE") != "release" {
		err := godotenv.Load("../../.env") // Chemin relatif au fichier .env
		if err != nil {
			log.Println("⚠️  .env non trouvé, on suppose que les variables sont déjà dans l'environnement")
		}
	}

	// ✅ Configurer le mode Gin (release = prod)
	mode := os.Getenv("GIN_MODE")
	if mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	r := gin.Default()

	// ✅ Routes publiques (à activer quand tu auras login/register)
	// public := r.Group("/api")
	// {
	//     public.POST("/login", auth.LoginHandler)
	//     public.POST("/register", auth.RegisterHandler)
	// }

	// 🔐 Routes privées (protégées par vérification de session JWT)
	private := r.Group("/api", auth.AuthMiddleware())
	{
		private.GET("/profile", user.GetProfileHandler)
		private.PUT("/profile", user.UpdateProfileHandler)
	}

	// ✅ Déterminer le port dynamiquement
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // par défaut en local
	}

	log.Printf("✅ Serveur lancé sur le port : %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("❌ Erreur au démarrage du serveur : %v", err)
	}
}
