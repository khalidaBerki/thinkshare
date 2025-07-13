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
	posts.GET("", h.GetAllPosts)
	posts.GET("/user/:id", h.GetPostsByUser) // Posts d'un utilisateur spécifique
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
// CreatePost godoc
// @Summary      Create a new post
// @Description  Create a new post with text and optional media (images, video, or documents)
// @Tags         posts
// @Security     BearerAuth
// @Accept       multipart/form-data
// @Produce      json
// @Param        content        formData  string  true   "Post content"
// @Param        visibility     formData  string  true   "Post visibility (public or private)"
// @Param        document_type  formData  string  false  "Document type (optional)"
// @Param        images         formData  file    false  "Images (max 10, only if no video/documents)"
// @Param        video          formData  file    false  "Video (only if no images/documents)"
// @Param        documents      formData  file    false  "Documents (max 5, only if no images/video)"
// @Success      201  {object}  post.PostDTO
// @Failure      400  {object}  map[string]string "Invalid input"
// @Failure      401  {object}  map[string]string "Unauthorized"
// @Failure      500  {object}  map[string]string "Internal server error"
// @Router       /api/posts [post]
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
	isPaidOnlyStr := getFirst(form.Value, "is_paid_only")

	// Convertir is_paid_only en booléen
	isPaidOnly := isPaidOnlyStr == "true"

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
		IsPaidOnly:   isPaidOnly,
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
// GetPostByID godoc
// @Summary      Get a post by ID
// @Description  Retrieve a post and its details by its ID
// @Tags         posts
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      int  true  "Post ID"
// @Success      200  {object}  post.PostDTO
// @Failure      400  {object}  map[string]string "Invalid post ID"
// @Failure      404  {object}  map[string]string "Post not found"
// @Router       /api/posts/{id} [get]
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
// GetAllPosts godoc
// @Summary      Get all posts (infinite scroll)
// @Description  Retrieve all posts with optional infinite scroll (after/limit)
// @Tags         posts
// @Security     BearerAuth
// @Produce      json
// @Param        after  query     int  false  "Last post ID already loaded (for infinite scroll)"
// @Param        limit  query     int  false  "Number of posts to return (default 20)"
// @Success      200  {object}   map[string]interface{} "List of posts and pagination info"
// @Failure      401  {object}   map[string]string "Unauthorized"
// @Failure      500  {object}   map[string]string "Internal server error"
// @Router       /api/posts [get]
func (h *Handler) GetAllPosts(c *gin.Context) {
	userID := c.GetInt("user_id")
	afterID, _ := strconv.Atoi(c.DefaultQuery("after", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20")) // on garder un limit pour éviter de tout charger

	posts, err := h.service.GetAllPostsAfter(uint(afterID), limit, uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	hasMore := len(posts) == limit
	c.JSON(http.StatusOK, gin.H{
		"posts":    posts,
		"has_more": hasMore,
		"last_id": func() uint {
			if len(posts) > 0 {
				return posts[len(posts)-1].ID
			}
			return 0
		}(),
	})
}

// PUT /posts/:id
// UpdatePost godoc
// @Summary      Update a post
// @Description  Update the content, visibility, or document type of a post
// @Tags         posts
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id    path      int                true  "Post ID"
// @Param        body  body      post.UpdatePostInput  true  "Fields to update"
// @Success      200   {object}  post.PostDTO
// @Failure      400   {object}  map[string]string "Invalid input"
// @Failure      401   {object}  map[string]string "Unauthorized"
// @Failure      403   {object}  map[string]string "Forbidden"
// @Failure      404   {object}  map[string]string "Post not found"
// @Router       /api/posts/{id} [put]
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
// DeletePost godoc
// @Summary      Delete a post
// @Description  Delete a post by its ID
// @Tags         posts
// @Security     BearerAuth
// @Param        id   path      int  true  "Post ID"
// @Success      204  "No Content"
// @Failure      401  {object}  map[string]string "Unauthorized"
// @Failure      403  {object}  map[string]string "Forbidden"
// @Failure      404  {object}  map[string]string "Post not found"
// @Router       /api/posts/{id} [delete]
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
// GetMediaStats godoc
// @Summary      Get media statistics
// @Description  Retrieve statistics and recommendations about uploaded media
// @Tags         posts
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  map[string]interface{} "Media statistics and recommendations"
// @Failure      401  {object}  map[string]string "Unauthorized"
// @Failure      500  {object}  map[string]string "Internal server error"
// @Router       /api/posts/media/stats [get]
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

// GET /posts/user/:id
// GetPostsByUser godoc
// @Summary      Get posts by user
// @Description  Retrieve all posts created by a specific user (with infinite scroll)
// @Tags         posts
// @Security     BearerAuth
// @Produce      json
// @Param        id     path      int  true  "User ID"
// @Param        after  query     int  false "Last post ID already loaded (for infinite scroll)"
// @Param        limit  query     int  false "Number of posts to return (default 20)"
// @Success      200  {object}   map[string]interface{} "List of posts and pagination info"
// @Failure      400  {object}   map[string]string "Invalid user ID"
// @Failure      401  {object}   map[string]string "Unauthorized"
// @Failure      500  {object}   map[string]string "Internal server error"
// @Router       /api/posts/user/{id} [get]
func (h *Handler) GetPostsByUser(c *gin.Context) {
	userID := c.GetInt("user_id")
	creatorID, err := strconv.Atoi(c.Param("id"))
	if err != nil || creatorID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	afterID, _ := strconv.Atoi(c.DefaultQuery("after", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	posts, err := h.service.GetPostsByCreatorAfter(uint(creatorID), uint(afterID), limit, uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	hasMore := len(posts) == limit
	c.JSON(http.StatusOK, gin.H{
		"posts":    posts,
		"has_more": hasMore,
		"last_id": func() uint {
			if len(posts) > 0 {
				return posts[len(posts)-1].ID
			}
			return 0
		}(),
	})
}
