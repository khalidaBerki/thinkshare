package post

import (
	"backend/internal/media"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	posts := rg.Group("/posts")
	posts.POST("/", h.CreatePost)
	posts.GET("/", h.GetAllPosts)
	posts.PUT("/:id", h.UpdatePost)
	posts.DELETE("/:id", h.DeletePost)
}

func (h *Handler) CreatePost(c *gin.Context) {
	userID := c.GetInt("user_id")
	log.Printf("📝 Création post par utilisateur ID: %d", userID)

	content := c.PostForm("content")
	visibility := c.PostForm("visibility")
	log.Printf("📝 Contenu: '%s', Visibilité: %s", content, visibility)

	if visibility != string(Public) && visibility != string(Private) {
		log.Printf("❌ Visibilité invalide: %s", visibility)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid visibility value"})
		return
	}

	// Journaliser les types MIME acceptés
	log.Printf("📋 Content-Type: %s", c.GetHeader("Content-Type"))

	form, err := c.MultipartForm()
	if err != nil {
		log.Printf("❌ Erreur lors de la lecture du formulaire multipart: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form data"})
		return
	}

	// Journaliser les clés présentes dans le formulaire
	log.Printf("🔑 Clés dans le formulaire: %v", getMapKeys(form.File))

	images := form.File["images"]
	videos := form.File["video"]

	log.Printf("🖼️ Nombre d'images: %d", len(images))
	log.Printf("🎬 Nombre de vidéos: %d", len(videos))

	if len(images) > 0 && len(videos) > 0 {
		log.Printf("❌ Tentative d'upload d'images ET de vidéos simultanément")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot upload both images and video in one post"})
		return
	}
	if len(images) > 10 {
		log.Printf("❌ Trop d'images: %d (max 10)", len(images))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Maximum 10 images allowed"})
		return
	}
	if len(videos) > 1 {
		log.Printf("❌ Trop de vidéos: %d (max 1)", len(videos))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only one video allowed"})
		return
	}

	var medias []media.Media
	for _, img := range images {
		log.Printf("🖼️ Traitement image: %s, taille: %d bytes, type MIME: %s",
			img.Filename, img.Size, img.Header.Get("Content-Type"))

		if !isValidImage(img.Filename) || !isUnderSize(img, 10*1024*1024) {
			log.Printf("❌ Format ou taille d'image invalide: %s (%d bytes)", img.Filename, img.Size)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid image format or size (max 10MB)"})
			return
		}
		path, err := saveFile(uint(userID), img)
		if err != nil {
			log.Printf("❌ Échec de la sauvegarde de l'image: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
			return
		}
		log.Printf("✅ Image enregistrée à: %s", path)
		medias = append(medias, media.Media{MediaURL: path, MediaType: "image"})
	}

	if len(videos) == 1 {
		video := videos[0]
		log.Printf("🎬 DÉBUT TRAITEMENT VIDÉO =================================================")
		log.Printf("🎬 Nom du fichier: %s", video.Filename)
		log.Printf("🎬 Taille: %d bytes (%.2f MB)", video.Size, float64(video.Size)/(1024*1024))
		log.Printf("🎬 Type MIME: %s", video.Header.Get("Content-Type"))

		// Vérification de l'extension
		ext := strings.ToLower(filepath.Ext(video.Filename))
		log.Printf("🎬 Extension détectée: %s", ext)
		if !isValidVideo(video.Filename) {
			log.Printf("❌ Format vidéo invalide: %s - Les formats acceptés sont .mp4, .mov, .webm", ext)
			c.JSON(http.StatusBadRequest, gin.H{
				"error":           "Format vidéo invalide - Les formats acceptés sont .mp4, .mov, .webm",
				"detected_format": ext,
				"filename":        video.Filename,
			})
			return
		}

		// Vérification de la taille
		if video.Size > 100*1024*1024 {
			log.Printf("❌ Taille vidéo trop grande: %.2f MB (max 100MB)", float64(video.Size)/(1024*1024))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Video size exceeds maximum allowed (100MB)"})
			return
		}

		// Débuter l'enregistrement
		log.Printf("⏳ Début de l'enregistrement de la vidéo...")

		// Créer le dossier uploads s'il n'existe pas
		dir := "uploads"
		if err := os.MkdirAll(dir, 0750); err != nil {
			log.Printf("❌ Erreur création dossier uploads: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la création du dossier uploads"})
			return
		}

		// Générer un nom de fichier unique
		unique := uuid.New().String()
		filepath := filepath.Join(dir, fmt.Sprintf("user_%d_video_%s%s", userID, unique, ext))
		log.Printf("📁 Chemin de destination: %s", filepath)

		// Essayer d'ouvrir le fichier source
		src, err := video.Open()
		if err != nil {
			log.Printf("❌ Erreur ouverture fichier source: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de l'ouverture du fichier vidéo"})
			return
		}
		defer src.Close()

		// Créer le fichier destination
		dst, err := os.Create(filepath)
		if err != nil {
			log.Printf("❌ Erreur création fichier destination: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la création du fichier destination"})
			return
		}
		defer dst.Close()

		// Copier les données
		log.Printf("⏳ Copie des données en cours...")
		bytes, err := io.Copy(dst, src)
		if err != nil {
			log.Printf("❌ Erreur copie de données (%d bytes copiés): %v", bytes, err)
			os.Remove(filepath)
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Erreur lors de la copie des données: %v", err)})
			return
		}

		log.Printf("✅ Vidéo enregistrée avec succès! %d bytes écrits", bytes)
		log.Printf("🎬 FIN TRAITEMENT VIDÉO ====================================================")

		// Ajouter au média
		medias = append(medias, media.Media{MediaURL: filepath, MediaType: "video"})
	}

	post := Post{
		CreatorID:  uint(userID),
		Content:    content,
		Visibility: Visibility(visibility),
		CreatedAt:  time.Now(),
		Media:      medias,
	}

	log.Printf("📝 Tentative de création du post: %d médias attachés", len(medias))
	if err := h.service.CreatePost(&post); err != nil {
		status := http.StatusInternalServerError
		if strings.Contains(err.Error(), "invalid") {
			status = http.StatusBadRequest
		}
		log.Printf("❌ Échec de la création du post: %v", err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	log.Printf("✅ Post créé avec succès! ID: %d", post.ID)
	c.JSON(http.StatusCreated, post)
}

func (h *Handler) GetAllPosts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	visibility := c.Query("visibility")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 50 {
		pageSize = 10
	}

	posts, err := h.service.GetAllPosts(page, pageSize, Visibility(visibility))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, posts)
}

func (h *Handler) UpdatePost(c *gin.Context) {
	postID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}
	userID := c.GetInt("user_id")

	var input UpdatePostInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdatePost(uint(postID), uint(userID), input); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

func (h *Handler) DeletePost(c *gin.Context) {
	postID, _ := strconv.Atoi(c.Param("id"))
	userID := c.GetInt("user_id")

	log.Printf("🗑️ Tentative de suppression du post ID %d par utilisateur ID %d", postID, userID)

	if err := h.service.DeletePost(uint(postID), uint(userID)); err != nil {
		log.Printf("❌ Échec de la suppression du post ID %d: %v", postID, err)

		// Déterminer le code d'état approprié
		statusCode := http.StatusForbidden
		if strings.Contains(err.Error(), "record not found") {
			statusCode = http.StatusNotFound
		} else if strings.Contains(err.Error(), "failed to delete media") {
			statusCode = http.StatusInternalServerError
		}

		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	log.Printf("✅ Post ID %d supprimé avec succès", postID)
	c.Status(http.StatusNoContent)
}

func isValidImage(name string) bool {
	ext := strings.ToLower(filepath.Ext(name))
	isValid := ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif" || ext == ".webp"
	log.Printf("🔍 Vérification format image: %s -> %v", ext, isValid)
	return isValid
}

func isValidVideo(name string) bool {
	ext := strings.ToLower(filepath.Ext(name))
	// Liste étendue des formats vidéo acceptés
	validFormats := map[string]bool{
		".mp4":  true,
		".mov":  true,
		".webm": true,
		".avi":  true,
		".mkv":  true,
		".flv":  true,
		".wmv":  true,
		".3gp":  true,
		".m4v":  true,
	}
	isValid := validFormats[ext]
	log.Printf("🔍 Vérification format vidéo: %s -> %v", ext, isValid)
	return isValid
}

func isUnderSize(f *multipart.FileHeader, max int64) bool {
	isValid := f.Size <= max
	log.Printf("🔍 Vérification taille fichier: %d bytes (max %d) -> %v", f.Size, max, isValid)
	return isValid
}

// Fonction utilitaire pour afficher les clés d'une map dans les logs
func getMapKeys(m map[string][]*multipart.FileHeader) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func saveFile(userID uint, f *multipart.FileHeader) (string, error) {
	log.Printf("💾 Début sauvegarde fichier: %s (taille: %d bytes)", f.Filename, f.Size)

	dir := "uploads"
	if err := os.MkdirAll(dir, 0750); err != nil {
		log.Printf("❌ Erreur création dossier uploads: %v", err)
		return "", err
	}

	// Vérifier simplement que le fichier ne dépasse pas une taille maximale
	// La vérification d'espace disque avancée est désactivée pour plus de compatibilité
	if f.Size > 100*1024*1024 { // 100 MB max pour une taille absolue
		log.Printf("❌ Fichier trop volumineux: %d bytes", f.Size)
		return "", errors.New("file too large")
	}

	ext := filepath.Ext(f.Filename)
	unique := uuid.New().String()
	path := filepath.Join(dir, fmt.Sprintf("user_%d_%s%s", userID, unique, ext))
	log.Printf("📁 Chemin de destination: %s", path)

	src, err := f.Open()
	if err != nil {
		log.Printf("❌ Erreur ouverture fichier source: %v", err)
		return "", err
	}
	defer src.Close()

	dst, err := os.Create(path)
	if err != nil {
		log.Printf("❌ Erreur création fichier destination: %v", err)
		return "", err
	}
	defer dst.Close()

	bytes, err := io.Copy(dst, src)
	if err != nil {
		log.Printf("❌ Erreur copie de données (%d bytes copiés): %v", bytes, err)
		os.Remove(path)
		return "", err
	}
	log.Printf("✅ Fichier enregistré avec succès: %d bytes écrits", bytes)

	return path, nil
}
