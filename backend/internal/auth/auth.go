package auth

import (
	"backend/internal/db"
	"backend/internal/user"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
	"golang.org/x/crypto/bcrypt"
)

func InitGoth() {
	goth.UseProviders(
		google.New(
			os.Getenv("GOOGLE_KEY"),
			os.Getenv("GOOGLE_SECRET"),
			"http://localhost:8080/auth/google/callback",
		),
	)
}

func ProviderMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		provider := c.Param("provider")
		if provider != "" {
			// Ajoute le provider à la query string pour gothic
			q := c.Request.URL.Query()
			q.Set("provider", provider)
			c.Request.URL.RawQuery = q.Encode()
		}
		c.Next()
	}
}

func RegisterRoutes(r *gin.Engine) {
	r.POST("/register", Register)
	r.POST("/login", Login)
	r.GET("/auth/:provider", ProviderMiddleware(), BeginAuthHandler)
	r.GET("/auth/:provider/callback", ProviderMiddleware(), CallbackHandler)
	r.GET("/logout", LogoutHandler)
}

type RegisterInput struct {
	Name      string `json:"name" binding:"required"`
	FirstName string `json:"firstname" binding:"required"`
	Username  string `json:"username" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
}

// Register godoc
// @Summary Créer un compte avec name, firstname, username, email, password
// @Tags Auth
// @Accept json
// @Produce json
// @Param input body RegisterInput true "Informations d'inscription"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /register [post]
func Register(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	u := user.User{
		Name:         input.Name,
		FirstName:    input.FirstName,
		Username:     input.Username,
		Email:        input.Email,
		PasswordHash: string(hashed),
		Role:         "user",
		CreatedAt:    time.Now(),
	}

	if err := db.GormDB.Create(&u).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email, username, name ou firstname déjà utilisé"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Utilisateur inscrit avec succès"})
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// Login godoc
// @Summary Connexion utilisateur (login email/password)
// @Tags Auth
// @Accept json
// @Produce json
// @Param input body LoginInput true "Identifiants de connexion"
// @Success 200 {object} TokenResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /login [post]
func Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var u user.User
	if err := db.GormDB.Where("email = ?", input.Email).First(&u).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "email ou mot de passe invalide"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "email ou mot de passe invalide"})
		return
	}

	token, err := GenerateJWT(int(u.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erreur lors de la génération du token"})
		return
	}

	c.JSON(http.StatusOK, TokenResponse{
		Token:  token,
		UserID: u.ID,
	})
}

// BeginAuthHandler godoc
// @Summary Début de l'authentification Google OAuth
// @Tags Auth
// @Produce json
// @Param provider path string true "google"
// @Success 302 {string} string "Redirection vers Google"
// @Router /auth/{provider} [get]
func BeginAuthHandler(c *gin.Context) {
	provider := c.Param("provider")
	if provider != "google" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Provider not supported"})
		return
	}
	gothic.BeginAuthHandler(c.Writer, c.Request)
}

// CallbackHandler godoc
// @Summary Callback OAuth Google
// @Tags Auth
// @Produce json
// @Param provider path string true "google"
// @Success 200 {object} user.User
// @Router /auth/{provider}/callback [get]
func CallbackHandler(c *gin.Context) {
	provider := c.Param("provider")
	if provider != "google" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Provider not supported"})
		return
	}
	gUser, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	var u user.User
	result := db.GormDB.Where("email = ?", gUser.Email).First(&u)
	if result.Error != nil {
		// Inscription via Google
		u = user.User{
			Username:     gUser.NickName,
			Email:        gUser.Email,
			PasswordHash: "",
			Role:         "google",
			CreatedAt:    time.Now(),
		}
		db.GormDB.Create(&u)
	}

	// Redirige vers la documentation Swagger après authentification
	c.Redirect(http.StatusTemporaryRedirect, "http://localhost:8080/swagger/index.html")
}

// LogoutHandler godoc
// @Summary Déconnexion utilisateur
// @Tags Auth
// @Produce json
// @Success 302 {string} string "Redirect vers /"
// @Router /logout [get]
func LogoutHandler(c *gin.Context) {
	gothic.Logout(c.Writer, c.Request)
	c.Redirect(http.StatusTemporaryRedirect, "/")
}
