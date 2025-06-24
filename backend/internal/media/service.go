package media

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Interface de service pour les médias
type Service interface {
	GetMediaByID(id uint) (*Media, error)
	GetMediasByPostID(postID uint) ([]Media, error)
	DeleteMedia(id uint) error
	UpdateMediaMetadata(id uint, metadata string) error
	CleanupOrphanedMedia() (int, error)
}

type serviceImpl struct {
	repo Repository
}

// Créer une nouvelle instance du service
func NewService(repo Repository) Service {
	return &serviceImpl{repo: repo}
}

// Récupérer un média par son ID
func (s *serviceImpl) GetMediaByID(id uint) (*Media, error) {
	return s.repo.FindByID(id)
}

// Récupérer tous les médias associés à un post
func (s *serviceImpl) GetMediasByPostID(postID uint) ([]Media, error) {
	return s.repo.FindByPostID(postID)
}

// Supprimer un média
func (s *serviceImpl) DeleteMedia(id uint) error {
	// Récupérer d'abord le média
	media, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	// Supprimer le fichier physique
	if err := os.Remove(media.MediaURL); err != nil && !os.IsNotExist(err) {
		log.Printf("⚠️ Impossible de supprimer le fichier %s: %v", media.MediaURL, err)
		// On continue même si le fichier n'a pas pu être supprimé
	}

	// Supprimer la miniature si elle existe
	if media.ThumbnailURL != "" {
		if err := os.Remove(media.ThumbnailURL); err != nil && !os.IsNotExist(err) {
			log.Printf("⚠️ Impossible de supprimer la miniature %s: %v", media.ThumbnailURL, err)
		}
	}

	// Supprimer l'entrée en base de données
	return s.repo.Delete(id)
}

// Mettre à jour les métadonnées d'un média
func (s *serviceImpl) UpdateMediaMetadata(id uint, metadata string) error {
	// Récupérer d'abord le média
	media, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	// Mettre à jour les métadonnées
	media.Metadata = metadata

	// Enregistrer les modifications
	return s.repo.Update(media)
}

// Nettoyer les fichiers médias orphelins (sans référence en BDD)
func (s *serviceImpl) CleanupOrphanedMedia() (int, error) {
	// Récupérer tous les médias en base de données
	allMedia, err := s.repo.FindAll()
	if err != nil {
		return 0, err
	}

	// Créer une map pour recherche rapide
	mediaMap := make(map[string]bool)
	thumbnailMap := make(map[string]bool)

	for _, m := range allMedia {
		mediaMap[m.MediaURL] = true
		if m.ThumbnailURL != "" {
			thumbnailMap[m.ThumbnailURL] = true
		}
	}

	// Compter les fichiers supprimés
	deleted := 0

	// Parcourir les dossiers de média
	foldersToCheck := []string{"uploads/images", "uploads/videos", "uploads/documents", "uploads/thumbnails"}

	for _, folder := range foldersToCheck {
		// Vérifier si le dossier existe
		if _, err := os.Stat(folder); os.IsNotExist(err) {
			continue
		}

		// Parcourir tous les fichiers du dossier
		err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// Ignorer les dossiers
			if info.IsDir() {
				return nil
			}

			// Vérifier si le fichier est référencé en BDD
			if strings.Contains(path, "/thumbnails/") {
				// Pour les miniatures
				if !thumbnailMap[path] {
					// Le fichier n'est pas référencé, on le supprime
					log.Printf("🗑️ Suppression d'une miniature orpheline: %s", path)
					if err := os.Remove(path); err == nil {
						deleted++
					}
				}
			} else {
				// Pour les autres fichiers média
				if !mediaMap[path] {
					// Le fichier n'est pas référencé, on le supprime
					log.Printf("🗑️ Suppression d'un média orphelin: %s", path)
					if err := os.Remove(path); err == nil {
						deleted++
					}
				}
			}

			return nil
		})

		if err != nil {
			return deleted, fmt.Errorf("erreur lors du parcours du dossier %s: %v", folder, err)
		}
	}

	return deleted, nil
}
