package main

import (
	"log"
	"os"
	"path/filepath"
	"time"

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

	// ✅ S'assurer que le dossier uploads existe avec les bonnes permissions
	uploadsDir := "uploads"
	// Vérifier si le dossier existe
	if _, err := os.Stat(uploadsDir); os.IsNotExist(err) {
		log.Printf("📁 Création du dossier uploads...")
		if err := os.MkdirAll(uploadsDir, 0755); err != nil {
			log.Printf("❌ Erreur lors de la création du dossier uploads: %v", err)
		} else {
			log.Printf("✅ Dossier uploads créé avec succès")
		}
	} else {
		// Si le dossier existe, on vérifie les permissions
		info, _ := os.Stat(uploadsDir)
		mode := info.Mode()
		log.Printf("📁 Dossier uploads existe déjà avec les permissions: %v", mode)

		// Vérifier si on peut écrire dans le dossier
		testFile := filepath.Join(uploadsDir, "test_write_permission.tmp")
		if f, err := os.Create(testFile); err != nil {
			log.Printf("❌ ATTENTION: Impossible d'écrire dans le dossier uploads: %v", err)
		} else {
			f.Close()
			os.Remove(testFile)
			log.Printf("✅ Les permissions d'écriture dans uploads sont correctes")
		}
	}

	// ✅ Démarrage serveur
	r := gin.Default()

	// Middleware CORS pour permettre les requêtes cross-origin
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Page d'accueil simple
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Bienvenue sur ThinkShare API",
			"version": "1.0.0",
			"endpoints": gin.H{
				"auth": []string{"/register", "/login", "/auth/google"},
				"api":  []string{"/api/posts", "/api/comments", "/api/profile"},
				"docs": "/swagger/index.html",
			},
		})
	})

	// Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 🔐 Auth routes (registration, login, OAuth)
	auth.InitGoth()
	auth.RegisterRoutes(r)

	// 🔧 Routes de debug (seulement en mode développement)
	if gin.Mode() != gin.ReleaseMode {
		r.GET("/debug/mode", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"mode":         gin.Mode(),
				"env_gin_mode": os.Getenv("GIN_MODE"),
				"is_debug":     gin.Mode() != gin.ReleaseMode,
				"server_time":  time.Now().Format(time.RFC3339),
			})
		})

		// Route pour tester l'authentification
		r.GET("/api/test-auth", auth.AuthMiddleware(), func(c *gin.Context) {
			userID := c.GetInt("user_id")
			c.JSON(200, gin.H{
				"message": "Authentification réussie",
				"user_id": userID,
				"time":    time.Now().Format(time.RFC3339),
			})
		})

		log.Printf("🔧 Routes de debug activées (mode développement)")
	}

	// 🔐 Routes API protégées
	api := r.Group("/api", auth.AuthMiddleware())
	{
		// 👤 Routes utilisateur
		api.GET("/profile", user.GetProfileHandler)
		api.PUT("/profile", user.UpdateProfileHandler)

		// 📝 Routes posts
		postRepo := post.NewRepository()
		postService := post.NewService(postRepo)
		postHandler := post.NewHandler(postService)
		postHandler.RegisterRoutes(api)

		// 💬 Routes commentaires
		commentRepo := comment.NewRepository(db.GormDB)
		commentService := comment.NewService(commentRepo, postRepo)
		commentHandler := comment.NewHandler(commentService)
		commentHandler.RegisterRoutes(api)

		// 💖 Routes likes
		likeRepo := like.NewRepository(db.GormDB)
		likeService := like.NewService(likeRepo, postRepo)
		likeHandler := like.NewHandler(likeService)
		likeHandler.RegisterRoutes(api)

		// 📩 Routes messagerie privée
		messageRepo := message.NewRepository(db.GormDB)
		messageService := message.NewService(messageRepo, db.GormDB)
		messageHandler := message.NewHandler(messageService)
		messageHandler.RegisterRoutes(api)

		log.Printf("✅ Routes API protégées configurées")
	}

	// Port dynamique
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Serveur ThinkShare lancé sur le port : %s", port)
	log.Printf("📚 Documentation Swagger : http://localhost:%s/swagger/index.html", port)

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("❌ Erreur de lancement : %v", err)
	}
}
