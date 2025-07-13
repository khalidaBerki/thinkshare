package post

import (
	"backend/internal/db"
	"backend/internal/models"
	"log"
)

// CheckPostAccess vérifie si un utilisateur a accès à un post payant
func CheckPostAccess(userID uint, creatorID uint, isPaidOnly bool) bool {
	// Si le post n'est pas payant, accès libre
	if !isPaidOnly {
		log.Printf("[ACCESS] userID=%d, creatorID=%d, isPaidOnly=%v => accès libre", userID, creatorID, isPaidOnly)
		return true
	}

	// Si c'est le créateur lui-même
	if userID == creatorID {
		log.Printf("[ACCESS] userID=%d est le créateur (creatorID=%d) => accès libre", userID, creatorID)
		return true
	}

	// Nouvelle logique : si l'utilisateur a au moins une subscription active (peu importe le type/status)
	var count int64
	err := db.GormDB.Model(&models.Subscription{}).
		Where("subscriber_id = ? AND creator_id = ? AND is_active = ?", userID, creatorID, true).
		Count(&count).Error
	if err != nil {
		log.Printf("[ACCESS][ERROR] Erreur DB lors du comptage des subscriptions: %v", err)
	}
	log.Printf("[ACCESS] userID=%d, creatorID=%d, isPaidOnly=%v, nb_subscriptions_actives=%d", userID, creatorID, isPaidOnly, count)

	return count > 0
}

// FilterPostsWithAccess filtre une liste de posts selon l'accès de l'utilisateur
func FilterPostsWithAccess(posts []*PostDTO, userID uint) []*PostDTO {
	var result []*PostDTO

	for _, post := range posts {
		hasAccess := CheckPostAccess(userID, post.CreatorID, post.IsPaidOnly)

		postDTO := &PostDTO{
			ID:           post.ID,
			CreatorID:    post.CreatorID,
			Visibility:   post.Visibility,
			IsPaidOnly:   post.IsPaidOnly,
			DocumentType: post.DocumentType,
			CreatedAt:    post.CreatedAt,
			UpdatedAt:    post.UpdatedAt,
			HasAccess:    hasAccess,
			MediaURLs:    post.MediaURLs,
		}

		if hasAccess {
			// L'utilisateur a accès, on montre le contenu complet
			postDTO.Content = post.Content
		} else {
			// L'utilisateur n'a pas accès, on montre un message
			postDTO.Content = "🔒 Ce contenu est réservé aux abonnés payants. Abonnez-vous pour y accéder !"
		}

		result = append(result, postDTO)
	}

	return result
}

// FilterPostsFromModelsWithAccess filtre une liste de posts modèles selon l'accès de l'utilisateur
func FilterPostsFromModelsWithAccess(posts []*Post, userID uint) []*PostDTO {
	var result []*PostDTO

	for _, post := range posts {
		hasAccess := CheckPostAccess(userID, post.CreatorID, post.IsPaidOnly)

		postDTO := &PostDTO{
			ID:           post.ID,
			CreatorID:    post.CreatorID,
			Visibility:   string(post.Visibility),
			IsPaidOnly:   post.IsPaidOnly,
			DocumentType: post.DocumentType,
			CreatedAt:    post.CreatedAt,
			UpdatedAt:    post.UpdatedAt,
			HasAccess:    hasAccess,
		}

		if hasAccess {
			// L'utilisateur a accès, on montre le contenu complet
			postDTO.Content = post.Content
		} else {
			// L'utilisateur n'a pas accès, on montre un message
			postDTO.Content = "🔒 Ce contenu est réservé aux abonnés payants. Abonnez-vous pour y accéder !"
		}

		result = append(result, postDTO)
	}

	return result
}
