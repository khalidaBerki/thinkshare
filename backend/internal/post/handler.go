package post

import (
	"backend/internal/media"
	"log"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// Handler structure
type Handler struct {
	service Service
}

// NewHandler instancie un gestionnaire de route
func NewHandler(s Service) *Handler {
	if err := checkUploadsDirectorySecurity(); err != nil {
		log.Printf("⚠️ Problème dossier uploads: %v", err)
	}
	return &Handler{service: s}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	posts := rg.Group("/posts")

	posts.POST("/", h.CreatePost)
	posts.GET("/", h.GetAllPosts)
	posts.GET("/:id", h.GetPostByID)
	posts.PUT("/:id", h.UpdatePost)
	posts.DELETE("/:id", h.DeletePost)

	posts.GET("/media/stats", h.GetMediaStats)
}

// Utilitaire: extraire les clés du form
func getMapKeys(m map[string][]*multipart.FileHeader) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// POST /posts : Créer un post
func (h *Handler) CreatePost(c *gin.Context) {
	userID := c.GetInt("user_id")
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 2<<30) // 2GB max

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Formulaire invalide (taille max dépassée ?)"})
		return
	}

	getFirst := func(m map[string][]string, key string) string {
		if v, ok := m[key]; ok && len(v) > 0 {
			return v[0]
		}
		return ""
	}
	content := getFirst(form.Value, "content")
	visibility := getFirst(form.Value, "visibility")
	documentType := getFirst(form.Value, "document_type")

	images := form.File["images"]
	videos := form.File["video"]
	documents := form.File["documents"]

	if visibility != string(Public) && visibility != string(Private) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Visibility invalide"})
		return
	}

	// Restrictions combinées
	if (len(images) > 0 && len(videos) > 0) || (len(images) > 0 && len(documents) > 0) || (len(videos) > 0 && len(documents) > 0) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Types médias multiples non autorisés (image OU vidéo OU document)"})
		return
	}
	if len(images) > 10 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Maximum 10 images autorisées"})
		return
	}
	if len(videos) > 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Une seule vidéo autorisée"})
		return
	}
	if len(documents) > 5 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Maximum 5 documents autorisés"})
		return
	}

	var medias []media.Media

	// Images
	for _, img := range images {
		if !IsValidImage(img.Filename) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Format image invalide"})
			return
		}
		if !IsUnderSize(img, 100*1024*1024) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Image trop lourde (max 100MB)"})
			return
		}
		path, _, fileSize, err := saveFile(uint(userID), img)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur sauvegarde image"})
			return
		}
		medias = append(medias, media.Media{MediaURL: path, MediaType: "image", FileSize: fileSize})
	}

	// Documents
	for _, doc := range documents {
		if !IsValidDocument(doc.Filename) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Format document invalide"})
			return
		}
		if !IsUnderSize(doc, 200*1024*1024) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Document trop lourd (max 200MB)"})
			return
		}
		path, _, fileSize, err := saveFile(uint(userID), doc)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur sauvegarde document"})
			return
		}
		medias = append(medias, media.Media{MediaURL: path, MediaType: "document", FileSize: fileSize})
	}

	// Vidéo
	if len(videos) == 1 {
		video := videos[0]
		if !IsValidVideo(video.Filename) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Format vidéo invalide"})
			return
		}
		if !IsUnderSize(video, 2*1024*1024*1024) { // 2GB max pour la vidéo
			c.JSON(http.StatusBadRequest, gin.H{"error": "Vidéo trop lourde (max 2GB)"})
			return
		}
		path, _, fileSize, err := saveFile(uint(userID), video)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur sauvegarde vidéo"})
			return
		}
		medias = append(medias, media.Media{MediaURL: path, MediaType: "video", FileSize: fileSize})
	}

	input := CreatePostInput{
		Content:      content,
		Visibility:   Visibility(visibility),
		DocumentType: documentType,
		Media:        medias,
	}

	postDTO, err := h.service.CreatePost(uint(userID), input)
	if err != nil {
		status := http.StatusInternalServerError
		if strings.Contains(err.Error(), "invalid") || strings.Contains(err.Error(), "invalide") {
			status = http.StatusBadRequest
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, postDTO)
}

// GET /posts/:id
func (h *Handler) GetPostByID(c *gin.Context) {
	userID := c.GetInt("user_id")
	postID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de post invalide"})
		return
	}
	post, err := h.service.GetPostByID(uint(postID), uint(userID))
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "post non trouvé" {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, post)
}

// GET /posts
func (h *Handler) GetAllPosts(c *gin.Context) {
	userID := c.GetInt("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	posts, total, err := h.service.GetAllPosts(page, limit, uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	totalPages := (int(total) + limit - 1) / limit
	c.JSON(http.StatusOK, gin.H{
		"posts": posts,
		"pagination": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": totalPages,
			"has_next":    page < totalPages,
			"has_prev":    page > 1,
		},
	})
}

// PUT /posts/:id
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

	postDTO, err := h.service.UpdatePost(uint(postID), uint(userID), input)
	if err != nil {
		status := http.StatusForbidden
		if strings.Contains(err.Error(), "post non trouvé") {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, postDTO)
}

// DELETE /posts/:id
func (h *Handler) DeletePost(c *gin.Context) {
	postID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}
	userID := c.GetInt("user_id")

	if err := h.service.DeletePost(uint(postID), uint(userID)); err != nil {
		status := http.StatusForbidden
		if strings.Contains(err.Error(), "post non trouvé") {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// GET /posts/media/stats
func (h *Handler) GetMediaStats(c *gin.Context) {
	statistics, recommendations := h.service.GetMediaStatistics()
	if statistics == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve media stats"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"statistics":      statistics,
		"recommendations": recommendations,
	})
}
