package post

import (
	"backend/internal/media"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// Structure principale du gestionnaire HTTP
type Handler struct {
	service Service
}

// Créer une nouvelle instance du gestionnaire
func NewHandler(s Service) *Handler {
	// Vérifier la sécurité du dossier uploads au démarrage
	if err := checkUploadsDirectorySecurity(); err != nil {
		log.Printf("⚠️ AVERTISSEMENT: Problème avec le dossier uploads: %v", err)
		log.Printf("⚠️ Les téléchargements de fichiers pourraient ne pas fonctionner correctement")
	}

	return &Handler{service: s}
}

// Enregistrer les routes du gestionnaire
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	posts := rg.Group("/posts")

	// Routes principales
	posts.POST("/", h.CreatePost)
	posts.GET("/", h.GetAllPosts)
	posts.GET("/:id", h.GetPostByID)
	posts.PUT("/:id", h.UpdatePost)
	posts.DELETE("/:id", h.DeletePost)

	// Statistiques
	posts.GET("/media/stats", h.GetMediaStats)
}

// Retourne les clés d'une map comme slice de string
func getMapKeys(m map[string][]*multipart.FileHeader) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// --- HANDLERS DES ROUTES ---

// Créer un nouveau post
func (h *Handler) CreatePost(c *gin.Context) {
	userID := c.GetInt("user_id")
	log.Printf("📝 Création post par utilisateur ID: %d", userID)

	content := c.PostForm("content")
	visibility := c.PostForm("visibility")
	documentType := c.PostForm("document_type") // Type de document (cours, devoir, support, etc.)
	log.Printf("📝 Contenu: '%s', Visibilité: %s, Type document: %s", content, visibility, documentType)

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
	documents := form.File["documents"]

	log.Printf("🖼️ Nombre d'images: %d", len(images))
	log.Printf("🎬 Nombre de vidéos: %d", len(videos))
	log.Printf("📄 Nombre de documents: %d", len(documents))

	// Vérifier les combinaisons de médias non autorisées
	mediaTypes := 0
	if len(images) > 0 {
		mediaTypes++
	}
	if len(videos) > 0 {
		mediaTypes++
	}
	if len(documents) > 0 {
		mediaTypes++
	}

	if mediaTypes > 1 {
		log.Printf("❌ Tentative d'upload de plusieurs types de médias simultanément")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot upload different media types in one post (images, videos, or documents)"})
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
	if len(documents) > 5 {
		log.Printf("❌ Trop de documents: %d (max 5)", len(documents))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Maximum 5 documents allowed"})
		return
	}

	var medias []media.Media

	// Traiter les images
	if len(images) > 0 {
		for _, img := range images {
			log.Printf("🖼️ Traitement image: %s, taille: %d bytes, type MIME: %s",
				img.Filename, img.Size, img.Header.Get("Content-Type"))

			isValid, _ := isUnderSize(img, 10*1024*1024)
			if !isValidImage(img.Filename) || !isValid {
				log.Printf("❌ Format ou taille d'image invalide: %s (%d bytes)", img.Filename, img.Size)
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid image format or size (max 10MB)"})
				return
			}
			path, _, fileSize, err := saveFile(uint(userID), img)
			if err != nil {
				log.Printf("❌ Échec de la sauvegarde de l'image: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
				return
			}
			log.Printf("✅ Image enregistrée à: %s", path)
			medias = append(medias, media.Media{MediaURL: path, MediaType: "image", FileSize: fileSize})
		}
	}

	// Traiter les documents
	if len(documents) > 0 {
		for _, doc := range documents {
			log.Printf("📄 DÉBUT TRAITEMENT DOCUMENT ==================================================")
			log.Printf("📄 Nom du fichier: %s", doc.Filename)
			log.Printf("📄 Taille: %d bytes (%.2f MB)", doc.Size, float64(doc.Size)/(1024*1024))
			log.Printf("📄 Type MIME: %s", doc.Header.Get("Content-Type"))

			// Vérification de l'extension
			if !isValidDocument(doc.Filename) {
				formats := strings.Join(getDocumentFormatList(), ", ")
				log.Printf("❌ Format document invalide: %s - Les formats acceptés sont %s",
					strings.ToLower(filepath.Ext(doc.Filename)), formats)
				c.JSON(http.StatusBadRequest, gin.H{
					"error":           fmt.Sprintf("Format document invalide - Les formats acceptés sont %s", formats),
					"detected_format": strings.ToLower(filepath.Ext(doc.Filename)),
					"filename":        doc.Filename,
				})
				return
			}

			// Vérification de la taille
			if doc.Size > 20*1024*1024 {
				log.Printf("❌ Taille document trop grande: %.2f MB (max 20MB)", float64(doc.Size)/(1024*1024))
				c.JSON(http.StatusBadRequest, gin.H{"error": "Document size exceeds maximum allowed (20MB)"})
				return
			}

			// Vérification de la qualité éducative du document
			valid, message := validateEducationalDocumentQuality(doc)
			if !valid {
				log.Printf("❌ Document ne respecte pas les critères de qualité: %s", message)
				c.JSON(http.StatusBadRequest, gin.H{"error": message})
				return
			}

			// Analyser les informations du document
			docInfo := getDocumentInfo(doc)
			log.Printf("📄 Information document: Format=%s, Taille=%.2fMB, Type=%s, PDF=%v",
				docInfo.Format, float64(docInfo.FileSize)/(1024*1024), docInfo.Category, docInfo.IsPDF)

			// Enregistrer le document
			path, _, fileSize, err := saveFile(uint(userID), doc)
			if err != nil {
				log.Printf("❌ Échec de la sauvegarde du document: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save document"})
				return
			}
			log.Printf("✅ Document enregistré à: %s", path)
			medias = append(medias, media.Media{MediaURL: path, MediaType: "document", FileSize: fileSize})
			log.Printf("📄 FIN TRAITEMENT DOCUMENT ==================================================")
		}
	}

	// Traiter les vidéos
	if len(videos) == 1 {
		video := videos[0]
		log.Printf("🎬 DÉBUT TRAITEMENT VIDÉO ==================================================")
		log.Printf("🎬 Nom du fichier: %s", video.Filename)
		log.Printf("🎬 Taille: %d bytes (%.2f MB)", video.Size, float64(video.Size)/(1024*1024))
		log.Printf("🎬 Type MIME: %s", video.Header.Get("Content-Type"))

		// Vérification de l'extension
		if !isValidVideo(video.Filename) {
			log.Printf("❌ Format vidéo invalide: %s - Les formats acceptés sont .mp4, .mov, .webm",
				strings.ToLower(filepath.Ext(video.Filename)))
			c.JSON(http.StatusBadRequest, gin.H{
				"error":           "Format vidéo invalide - Les formats acceptés sont .mp4, .mov, .webm",
				"detected_format": strings.ToLower(filepath.Ext(video.Filename)),
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

		// Enregistrer la vidéo
		path, _, fileSize, err := saveFile(uint(userID), video)
		if err != nil {
			log.Printf("❌ Échec de la sauvegarde de la vidéo: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save video"})
			return
		}

		log.Printf("✅ Vidéo enregistrée avec succès à %s", path)
		log.Printf("🎬 FIN TRAITEMENT VIDÉO =====================================================")

		// Ajouter au média
		medias = append(medias, media.Media{MediaURL: path, MediaType: "video", FileSize: fileSize})
	}

	// Créer le post
	post := Post{
		CreatorID:    uint(userID),
		Content:      content,
		Visibility:   Visibility(visibility),
		DocumentType: documentType,
		CreatedAt:    time.Now(),
		Media:        medias,
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

// Récupérer tous les posts
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

// Récupérer un post par son ID
func (h *Handler) GetPostByID(c *gin.Context) {
	// Convertir l'ID de la route en nombre
	postID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de post invalide"})
		return
	}

	// Récupérer le post depuis le service
	post, err := h.service.GetPostByID(uint(postID))
	if err != nil {
		statusCode := http.StatusInternalServerError

		// Si le post n'existe pas
		if strings.Contains(err.Error(), "record not found") {
			statusCode = http.StatusNotFound
		}

		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	// Retourner le post
	c.JSON(http.StatusOK, post)
}

// Mettre à jour un post existant
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

// Supprimer un post
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

// Obtenir des statistiques sur les médias
func (h *Handler) GetMediaStats(c *gin.Context) {
	// Récupérer les statistiques sur les médias par type
	stats, err := h.service.GetMediaStatistics()
	if err != nil {
		log.Printf("❌ Erreur lors de la récupération des statistiques des médias: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve media statistics"})
		return
	}

	// Ajouter des recommandations pour les formats
	recommendations := getRecommendedFormats()

	// Préparer la réponse
	response := gin.H{
		"statistics":      stats,
		"recommendations": recommendations,
	}

	c.JSON(http.StatusOK, response)
}
