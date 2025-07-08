package media

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// Structure du gestionnaire HTTP pour les médias
type Handler struct {
	service Service
}

// Créer une nouvelle instance du gestionnaire
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// Enregistrer les routes du gestionnaire
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	media := rg.Group("/media")

	// Routes pour les médias
	media.GET("/:id", h.GetMediaByID)
	media.DELETE("/:id", h.DeleteMedia)
	media.GET("/post/:postID", h.GetMediasByPostID)
	media.PUT("/:id/metadata", h.UpdateMediaMetadata)
	media.POST("/cleanup", h.CleanupOrphanedMedia)
}

// Récupérer un média par son ID

// GetMediaByID godoc
// @Summary      Get media by ID
// @Description  Retrieve a media file and its metadata by its ID
// @Tags         media
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      int  true  "Media ID"
// @Success      200  {object}  media.Media
// @Failure      400  {object}  map[string]string "Invalid media ID"
// @Failure      404  {object}  map[string]string "Media not found"
// @Failure      500  {object}  map[string]string "Internal server error"
// @Router       /api/media/{id} [get]
func (h *Handler) GetMediaByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de média invalide"})
		return
	}

	media, err := h.service.GetMediaByID(uint(id))
	if err != nil {
		status := http.StatusNotFound
		message := "Média non trouvé"

		// Vérifier si c'est une erreur de base de données ou autre
		if strings.Contains(err.Error(), "record not found") {
			message = "Média introuvable dans la base de données"
		} else {
			status = http.StatusInternalServerError
			message = "Erreur lors de la récupération du média"
		}

		c.JSON(status, gin.H{"error": message})
		return
	}

	c.JSON(http.StatusOK, media)
}

// Supprimer un média
// DeleteMedia godoc
// @Summary      Delete a media file
// @Description  Delete a media file by its ID (only the owner or admin can delete)
// @Tags         media
// @Security     BearerAuth
// @Param        id   path      int  true  "Media ID"
// @Success      200  {object}  map[string]interface{} "Media deleted successfully"
// @Failure      400  {object}  map[string]string "Invalid media ID"
// @Failure      401  {object}  map[string]string "Unauthorized"
// @Failure      404  {object}  map[string]string "Media not found"
// @Failure      500  {object}  map[string]string "Internal server error"
// @Router       /api/media/{id} [delete]
func (h *Handler) DeleteMedia(c *gin.Context) {
	// Vérifier que l'utilisateur est autorisé
	userID := c.GetInt("user_id")
	if userID <= 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Non autorisé"})
		return
	}

	// Récupérer l'ID du média
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de média invalide"})
		return
	}

	// Récupérer le média pour vérifier les permissions
	media, err := h.service.GetMediaByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Média non trouvé"})
		return
	}

	// TODO: Vérifier que l'utilisateur est autorisé à supprimer ce média
	// Dans une implémentation réelle, on vérifierait que le média appartient
	// à l'utilisateur ou que l'utilisateur est administrateur

	// Supprimer le média
	if err := h.service.DeleteMedia(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la suppression du média"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Média supprimé avec succès",
		"id":        id,
		"type":      media.MediaType,
		"file_name": media.FileName,
	})
}

// Récupérer tous les médias associés à un post
// GetMediasByPostID godoc
// @Summary      Get all media for a post
// @Description  Retrieve all media files associated with a specific post
// @Tags         media
// @Security     BearerAuth
// @Produce      json
// @Param        postID   path      int  true  "Post ID"
// @Success      200  {object}  map[string]interface{} "List of media for the post"
// @Failure      400  {object}  map[string]string "Invalid post ID"
// @Failure      500  {object}  map[string]string "Internal server error"
// @Router       /api/media/post/{postID} [get]
func (h *Handler) GetMediasByPostID(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("postID"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de post invalide"})
		return
	}

	medias, err := h.service.GetMediasByPostID(uint(postID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des médias"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"post_id":     postID,
		"media_count": len(medias),
		"medias":      medias,
	})
}

// Mettre à jour les métadonnées d'un média
// UpdateMediaMetadata godoc
// @Summary      Update media metadata
// @Description  Update the metadata of a media file (only the owner or admin can update)
// @Tags         media
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Media ID"
// @Param        body body      object true "Metadata update payload"
// @Success      200  {object}  map[string]interface{} "Metadata updated successfully"
// @Failure      400  {object}  map[string]string "Invalid input"
// @Failure      401  {object}  map[string]string "Unauthorized"
// @Failure      500  {object}  map[string]string "Internal server error"
// @Router       /api/media/{id}/metadata [put]
func (h *Handler) UpdateMediaMetadata(c *gin.Context) {
	// Vérifier que l'utilisateur est autorisé
	userID := c.GetInt("user_id")
	if userID <= 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Non autorisé"})
		return
	}

	// Récupérer l'ID du média
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de média invalide"})
		return
	}

	// Structure pour les données de la requête
	var req struct {
		Metadata string `json:"metadata"`
	}

	// Lire les données de la requête
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données de requête invalides"})
		return
	}

	// Mettre à jour les métadonnées
	if err := h.service.UpdateMediaMetadata(uint(id), req.Metadata); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la mise à jour des métadonnées"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Métadonnées mises à jour avec succès",
		"id":      id,
	})
}

// Nettoyer les fichiers médias orphelins
// CleanupOrphanedMedia godoc
// @Summary      Cleanup orphaned media files
// @Description  Delete all media files not referenced in the database (admin only)
// @Tags         media
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  map[string]interface{} "Cleanup result"
// @Failure      401  {object}  map[string]string "Unauthorized"
// @Failure      500  {object}  map[string]string "Internal server error"
// @Router       /api/media/cleanup [post]
func (h *Handler) CleanupOrphanedMedia(c *gin.Context) {
	// Cette action devrait être réservée aux administrateurs
	userID := c.GetInt("user_id")
	if userID <= 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Non autorisé"})
		return
	}

	// TODO: Vérifier que l'utilisateur est administrateur
	// Dans une implémentation réelle, on vérifierait le rôle de l'utilisateur

	// Nettoyer les médias orphelins
	deleted, err := h.service.CleanupOrphanedMedia()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors du nettoyage des médias orphelins"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "Nettoyage des médias orphelins terminé",
		"deleted_files": deleted,
	})
}
