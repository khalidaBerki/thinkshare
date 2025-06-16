package main

import (
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

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title ThinkShare Auth API
// @version 1.0
// @description API d'authentification pour ThinkShare : Google + formulaire
// @host localhost:8080
// @BasePath /

func main() {
	// ✅ Initialise la connexion DB une fois pour toute
	db.InitDB()

	logs := []struct {
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

	for _, l := range logs {
		err := db.GormDB.AutoMigrate(l.model)
		if err != nil {
			println("❌ Erreur migration table", l.name, ":", err.Error())
		} else {
			println("✅ Table", l.name, "OK ou déjà existante")
		}
	}

	r := gin.Default()

	auth.InitGoth()        // Google OAuth
	auth.RegisterRoutes(r) // tes routes d'auth

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Run(":8080")
}
