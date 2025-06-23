package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"

	_ "backend/docs"

	"backend/internal/auth"
	"backend/internal/comment"
	"backend/internal/db"
	"backend/internal/like"
	"backend/internal/media"
	"backend/internal/message"
	"backend/internal/post"
	"backend/internal/postaccess"
	"backend/internal/subscription"
	"backend/internal/user"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	// ✅ Configurer le mode Gin
	mode := os.Getenv("GIN_MODE")
	if mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// ✅ Initialiser la DB
	db.InitDB()

	// ✅ Migrer les tables (si nécessaires)
	migrations := []struct {
		name  string
		model interface{}
	}{
		{"users", &user.User{}},
		{"auth_tokens", &auth.AuthToken{}},
		{"posts", &post.Post{}},
		{"comments", &comment.Comment{}},
		{"likes", &like.Like{}},
		{"media", &media.Media{}},
		{"subscriptions", &subscription.Subscription{}},
		{"messages", &message.Message{}},
		{"postaccess", &postaccess.PostAccess{}},
	}

	for _, m := range migrations {
		if err := db.GormDB.AutoMigrate(m.model); err != nil {
			log.Printf("❌ Erreur migration %s : %v", m.name, err)
		} else {
			log.Printf("✅ Table %s migrée ou déjà existante", m.name)
		}
	}

	// ✅ Démarrage serveur
	r := gin.Default()

	// Ajoute cette route :
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Bienvenue sur ThinkShare"})
	})

	// Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Auth (Google, login, etc.)
	auth.InitGoth()
	auth.RegisterRoutes(r)

	// ✅ ✅ ✅ FAKE LOGIN - UNIQUEMENT POUR LE DEV ✅ ✅ ✅
	if os.Getenv("GIN_MODE") != "release" {
		r.POST("/api/fake-login", func(c *gin.Context) {
			// Simule un utilisateur avec ID = 1
			token, err := auth.GenerateJWT(1)
			if err != nil {
				c.JSON(500, gin.H{"error": "Impossible de générer le token"})
				return
			}
			c.JSON(200, gin.H{"token": token})
		})
		// ❗ À supprimer avant le déploiement prod
	}
	// 🔐 Routes API protégées
	api := r.Group("/api", auth.AuthMiddleware())
	{
		api.GET("/profile", user.GetProfileHandler)
		api.PUT("/profile", user.UpdateProfileHandler)
		api.POST("/subscribe", auth.AuthMiddleware(), subscription.SubscribeHandler)
		api.POST("/unsubscribe", auth.AuthMiddleware(), subscription.UnsubscribeHandler)
		r.GET("/api/followers/:id", auth.AuthMiddleware(), subscription.GetFollowersByUserHandler)
		r.GET("/api/subscriptions", auth.AuthMiddleware(), subscription.GetMySubscriptionsHandler)
	}

	// Port dynamique
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Serveur lancé sur le port : %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("❌ Erreur de lancement : %v", err)
	}
}
