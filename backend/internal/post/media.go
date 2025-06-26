package post

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Types de médias supportés
const (
	ImageType    = "image"
	VideoType    = "video"
	DocumentType = "document"
)

// Structure pour stocker les informations sur une image
type ImageInfo struct {
	Width       int     `json:"width"`
	Height      int     `json:"height"`
	Format      string  `json:"format"`
	FileSize    int64   `json:"filesize"`
	AspectRatio float64 `json:"aspect_ratio"`
	IsSquare    bool    `json:"is_square"`
	IsLandscape bool    `json:"is_landscape"`
	IsPortrait  bool    `json:"is_portrait"`
}

// Structure pour stocker les informations sur une vidéo
type VideoInfo struct {
	Width        int     `json:"width"`
	Height       int     `json:"height"`
	Duration     string  `json:"duration"`
	DurationSecs float64 `json:"duration_seconds"`
	Codec        string  `json:"codec"`
	Bitrate      string  `json:"bitrate"`
	Framerate    string  `json:"framerate"`
	FileSize     int64   `json:"filesize"`
	HasAudio     bool    `json:"has_audio"`
}

// Structure pour stocker les informations sur un document
type DocumentInfo struct {
	Format       string `json:"format"`
	FileSize     int64  `json:"filesize"`
	Category     string `json:"category"`
	DocumentType string `json:"document_type"`
	IsPDF        bool   `json:"is_pdf"`
	IsBinary     bool   `json:"is_binary"`
}

// --- FORMATS DE FICHIERS ---

// Récupérer la liste des formats document supportés
func getDocumentFormatList() []string {
	return []string{
		".pdf",          // Format document portable
		".doc", ".docx", // Formats Microsoft Word
		".ppt", ".pptx", // Formats Microsoft PowerPoint
		".xls", ".xlsx", // Formats Microsoft Excel
		".txt",                 // Texte brut
		".csv",                 // Valeurs séparées par des virgules
		".md",                  // Markdown
		".odt", ".ods", ".odp", // OpenDocument
	}
}

// Récupérer la liste des formats vidéo supportés
func getVideoFormatList() []string {
	return []string{
		".mp4",  // Format le plus courant et compatible
		".webm", // Format web ouvert
		".mov",  // Format Apple QuickTime
		".avi",  // Format Microsoft
		".mkv",  // Format conteneur Matroska
		".m4v",  // Format Apple iTunes
	}
}

// Récupérer la liste des formats image supportés
func getImageFormatList() []string {
	return []string{
		".jpg", ".jpeg", // Format JPEG
		".png",  // Format PNG pour les images avec transparence
		".gif",  // Format GIF pour les animations simples
		".webp", // Format Web moderne avec compression améliorée
		".svg",  // Format vectoriel
	}
}

// Récupérer la liste des formats recommandés avec explications
func getRecommendedFormats() map[string]map[string]string {
	return map[string]map[string]string{
		"image": {
			".jpg":  "Pour les photos et images complexes (sans transparence)",
			".png":  "Pour les images avec transparence ou haute qualité",
			".webp": "Format moderne avec meilleure compression et transparence",
			".svg":  "Pour les graphiques vectoriels (logos, icônes, diagrammes)",
		},
		"video": {
			".mp4":  "Excellent choix pour une compatibilité universelle (H.264/H.265)",
			".webm": "Optimisé pour le web, compression efficace (VP9/AV1)",
			".mov":  "Bonne qualité, préféré pour les appareils Apple",
		},
		"document": {
			".pdf":  "Format portable universel, compatible avec tous les appareils",
			".docx": "Format Word pour l'édition, à convertir en PDF pour partage",
			".pptx": "Format PowerPoint pour présentations",
			".xlsx": "Format Excel pour données tabulaires et calculs",
			".txt":  "Texte brut simple et léger",
		},
	}
}

// --- VALIDATION DES FICHIERS ---

// Vérifie si un fichier est une image valide
func isValidImage(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))

	// Utiliser les formats définis dans la liste des formats supportés
	for _, format := range getImageFormatList() {
		if ext == format {
			return true
		}
	}

	return false
}

// Vérifie si un fichier est une vidéo valide
func isValidVideo(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))

	// Utiliser les formats définis dans la liste des formats supportés
	for _, format := range getVideoFormatList() {
		if ext == format {
			return true
		}
	}

	return false
}

// Vérifie si un fichier est un document valide
func isValidDocument(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))

	// Utiliser les formats définis dans la liste des formats supportés
	for _, format := range getDocumentFormatList() {
		if ext == format {
			return true
		}
	}

	return false
}

// Vérifie si un fichier est sous une taille maximale
func isUnderSize(f *multipart.FileHeader, max int64) (bool, string) {
	isValid := f.Size <= max
	var message string
	if !isValid {
		// Formater la taille en unité lisible
		fileSize := float64(f.Size) / (1024 * 1024) // en MB
		maxSize := float64(max) / (1024 * 1024)     // en MB
		message = fmt.Sprintf("Fichier trop volumineux: %.2f MB (maximum autorisé: %.2f MB)", fileSize, maxSize)
		log.Printf("❌ %s", message)

		// Vérifier si c'est un document et suggérer des alternatives
		ext := strings.ToLower(filepath.Ext(f.Filename))
		if ext == ".pdf" && f.Size > 10*1024*1024 {
			log.Printf("💡 Conseil: Considérez l'optimisation du PDF pour réduire sa taille")
		} else if ext == ".docx" || ext == ".pptx" || ext == ".xlsx" {
			log.Printf("💡 Conseil: Convertir en PDF pourrait réduire la taille du fichier")
		}
	} else {
		log.Printf("✅ Vérification taille fichier: %d bytes (max %d) -> OK", f.Size, max)
	}
	return isValid, message
}

// Vérifie si un fichier est potentiellement dangereux
func isSuspiciousFile(filename string) bool {
	// Liste d'extensions potentiellement dangereuses
	dangerousExtensions := map[string]bool{
		// Exécutables
		".exe": true,
		".bat": true,
		".cmd": true,
		".sh":  true,
		".com": true,
		".dll": true,
		".msi": true,
		".bin": true,
		".app": true,
		".dmg": true,

		// Scripts
		".php":  true,
		".js":   true,
		".vbs":  true,
		".ps1":  true,
		".py":   true,
		".rb":   true,
		".pl":   true,
		".asp":  true,
		".aspx": true,
		".jsp":  true,
		".cgi":  true,

		// Archives potentiellement dangereuses
		".jar": true,
		".war": true,
		".iso": true,

		// Macros et autres
		".scr": true,
		".reg": true,
		".inf": true,
		".hta": true,
	}

	// Vérifier l'extension
	ext := strings.ToLower(filepath.Ext(filename))
	if dangerousExtensions[ext] {
		log.Printf("⚠️ Extension de fichier potentiellement dangereuse détectée: %s", ext)
		return true
	}

	// Vérifier les doubles extensions (exemple: image.jpg.exe)
	nameParts := strings.Split(strings.ToLower(filename), ".")
	if len(nameParts) > 2 {
		// Ignorer la première partie (nom de base)
		for i := 1; i < len(nameParts)-1; i++ {
			extCandidate := "." + nameParts[len(nameParts)-1]
			if dangerousExtensions[extCandidate] {
				log.Printf("⚠️ Détection d'extension double potentiellement dangereuse: %s", filename)
				return true
			}
		}
	}

	// Vérifier les noms de fichiers suspects
	suspiciousPatterns := []string{
		"virus", "malware", "hack", "crack", "keygen", "pirate",
		"trojan", "exploit", "backdoor", "rootkit", "ransom",
	}

	lowerFilename := strings.ToLower(filename)
	for _, pattern := range suspiciousPatterns {
		if strings.Contains(lowerFilename, pattern) {
			log.Printf("⚠️ Nom de fichier suspect détecté: %s (contient '%s')", filename, pattern)
			return true
		}
	}

	return false
}

// Nettoyer le nom de fichier et éviter les injections
func sanitizeFileName(filename string) string {
	// Remplacer les caractères potentiellement problématiques
	sanitized := strings.Map(func(r rune) rune {
		if r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' || r == '-' || r == '_' || r == '.' {
			return r
		}
		return '_'
	}, filename)

	// S'assurer que le nom ne commence pas par un point (fichier caché)
	if strings.HasPrefix(sanitized, ".") {
		sanitized = "_" + sanitized[1:]
	}

	return sanitized
}

// Vérifier si le type MIME correspond à l'extension déclarée
func validateMimeType(file *multipart.FileHeader) bool {
	ext := strings.ToLower(filepath.Ext(file.Filename))
	mimeType := file.Header.Get("Content-Type")

	// Si le document_type semble suspect, on tente une détection plus précise
	if strings.Contains(mimeType, "application/octet-stream") {
		detectedType, err := detectMimeType(file)
		if err == nil && detectedType != mimeType {
			log.Printf("ℹ️ Type MIME détecté différent: %s au lieu de %s", detectedType, mimeType)
			mimeType = detectedType
		}
	}

	// Mapper des extensions aux types MIME attendus
	mimeMap := map[string][]string{
		// Images
		".jpg":  {"image/jpeg", "image/jpg"},
		".jpeg": {"image/jpeg", "image/jpg"},
		".png":  {"image/png"},
		".gif":  {"image/gif"},
		".webp": {"image/webp"},
		".svg":  {"image/svg+xml", "image/svg"},

		// Vidéos
		".mp4":  {"video/mp4", "application/mp4"},
		".webm": {"video/webm"},
		".mov":  {"video/quicktime"},
		".avi":  {"video/x-msvideo", "video/avi"},
		".mkv":  {"video/x-matroska"},
		".m4v":  {"video/x-m4v", "video/mp4"},

		// Documents
		".pdf":  {"application/pdf"},
		".doc":  {"application/msword"},
		".docx": {"application/vnd.openxmlformats-officedocument.wordprocessingml.document"},
		".xls":  {"application/vnd.ms-excel"},
		".xlsx": {"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"},
		".ppt":  {"application/vnd.ms-powerpoint"},
		".pptx": {"application/vnd.openxmlformats-officedocument.presentationml.presentation"},
		".txt":  {"text/plain"},
		".csv":  {"text/csv", "application/csv"},
		".md":   {"text/markdown", "text/plain"},
	}

	// Types spéciaux qui peuvent utiliser application/octet-stream
	specialBinaryTypes := map[string]bool{
		".docx": true,
		".xlsx": true,
		".pptx": true,
		".zip":  true,
		".mov":  true,
		".mp4":  true,
	}

	// Vérifier si le type MIME correspond à l'un des types attendus pour l'extension
	if validTypes, exists := mimeMap[ext]; exists {
		for _, validType := range validTypes {
			if strings.Contains(mimeType, validType) {
				return true
			}
		}

		// Cas spécial: certains types peuvent être envoyés comme application/octet-stream
		if specialBinaryTypes[ext] && strings.Contains(mimeType, "application/octet-stream") {
			log.Printf("ℹ️ Type MIME générique accepté pour %s: %s", ext, mimeType)
			return true
		}

		// Si on arrive ici, le type MIME ne correspond pas à l'extension
		log.Printf("⚠️ Type MIME suspect: %s ne correspond pas à l'extension %s", mimeType, ext)
		return false
	}

	// Pour les extensions non répertoriées, on accepte mais on journalise
	log.Printf("ℹ️ Extension non répertoriée: %s avec type MIME %s", ext, mimeType)
	return true
}

// Fonction pour déterminer le type MIME à partir d'un échantillon de fichier
// Utilise la signature magic number pour identifier le vrai type du fichier
func detectMimeType(file *multipart.FileHeader) (string, error) {
	// Ouvrir le fichier
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer func(src multipart.File) {
		err := src.Close()
		if err != nil {

		}
	}(src)

	// Lire les 512 premiers octets pour la détection du type
	buffer := make([]byte, 512)
	n, err := src.Read(buffer)
	if err != nil && err != io.EOF {
		return "", err
	}

	// Utiliser la fonction DetectContentType du package http
	documentType := http.DetectContentType(buffer[:n])

	// Vérifier des signatures spécifiques pour plus de précision
	if bytes.HasPrefix(buffer, []byte("%PDF")) {
		return "application/pdf", nil
	}

	if bytes.HasPrefix(buffer, []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}) {
		return "image/png", nil
	}

	if bytes.HasPrefix(buffer, []byte{0xFF, 0xD8}) {
		return "image/jpeg", nil
	}

	if bytes.HasPrefix(buffer, []byte("GIF87a")) || bytes.HasPrefix(buffer, []byte("GIF89a")) {
		return "image/gif", nil
	}

	// Check pour les fichiers MS Office (DOCX, XLSX, PPTX sont des archives ZIP)
	if bytes.HasPrefix(buffer, []byte{0x50, 0x4B, 0x03, 0x04}) {
		// C'est un ZIP, pourrait être DOCX/XLSX/PPTX
		extension := strings.ToLower(filepath.Ext(file.Filename))
		switch extension {
		case ".docx":
			return "application/vnd.openxmlformats-officedocument.wordprocessingml.document", nil
		case ".xlsx":
			return "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", nil
		case ".pptx":
			return "application/vnd.openxmlformats-officedocument.presentationml.presentation", nil
		case ".zip":
			return "application/zip", nil
		}
	}

	// Vérifier les exécutables Windows (commencent par MZ)
	if bytes.HasPrefix(buffer, []byte{0x4D, 0x5A}) {
		log.Printf("⚠️ ALERTE: Signature d'exécutable Windows (MZ) détectée dans le fichier %s", file.Filename)
		return "application/x-msdownload", nil
	}

	// Vérifier les scripts
	if bytes.Contains(buffer[:n], []byte("<?php")) {
		log.Printf("⚠️ ALERTE: Code PHP détecté dans le fichier %s", file.Filename)
		return "text/x-php", nil
	}

	if bytes.Contains(buffer[:n], []byte("<script")) {
		log.Printf("⚠️ ALERTE: Script JavaScript détecté dans le fichier %s", file.Filename)
		return "text/javascript", nil
	}

	return documentType, nil
}

// --- GESTION DES FICHIERS ---

// Enregistrer un fichier média sur le serveur
func saveFile(userID uint, f *multipart.FileHeader) (string, string, int64, error) {
	log.Printf("💾 Début sauvegarde fichier: %s (taille: %d bytes)", f.Filename, f.Size)

	// Vérifier si le fichier est potentiellement dangereux
	if isSuspiciousFile(f.Filename) {
		log.Printf("⚠️ Tentative d'upload d'un fichier potentiellement dangereux: %s", f.Filename)
		return "", "", 0, errors.New("format de fichier non autorisé pour des raisons de sécurité")
	}

	// Nettoyer le nom du fichier pour éviter les injections
	cleanFilename := sanitizeFileName(f.Filename)
	if cleanFilename != f.Filename {
		log.Printf("ℹ️ Nom de fichier nettoyé: %s -> %s", f.Filename, cleanFilename)
	}

	// Déterminer le type de fichier et le sous-dossier approprié
	var subDir string
	ext := strings.ToLower(filepath.Ext(cleanFilename))

	switch {
	case isValidImage(cleanFilename):
		subDir = "uploads/images"
	case isValidVideo(cleanFilename):
		subDir = "uploads/videos"
	case isValidDocument(cleanFilename):
		subDir = "uploads/documents"
		// Traitement spécial pour les PDF (journalisation)
		if strings.ToLower(ext) == ".pdf" {
			log.Printf("📄 Traitement de document PDF: %s", cleanFilename)
		}
	default:
		log.Printf("❌ Type de fichier non pris en charge: %s", ext)
		return "", "", 0, errors.New("type de fichier non pris en charge sur la plateforme")
	}

	// Créer le sous-dossier s'il n'existe pas
	if err := os.MkdirAll(subDir, 0750); err != nil {
		log.Printf("❌ Erreur création dossier %s: %v", subDir, err)
		return "", "", 0, fmt.Errorf("erreur système lors de la création du dossier: %v", err)
	}

	// Vérifier les limites de taille selon le type de fichier
	var maxSize int64
	var typeFichier string

	switch {
	case isValidImage(cleanFilename):
		maxSize = 10 * 1024 * 1024 // 10 MB pour les images
		typeFichier = "image"
	case isValidVideo(cleanFilename):
		maxSize = 100 * 1024 * 1024 // 100 MB pour les vidéos
		typeFichier = "vidéo"
	case isValidDocument(cleanFilename):
		maxSize = 20 * 1024 * 1024 // 20 MB pour les documents
		typeFichier = "document"
	default:
		maxSize = 5 * 1024 * 1024 // 5 MB pour les autres types
		typeFichier = "fichier"
	}

	if f.Size > maxSize {
		log.Printf("❌ %s trop volumineux: %.2f MB (max %.2f MB)",
			typeFichier, float64(f.Size)/(1024*1024), float64(maxSize)/(1024*1024))
		return "", "", 0, fmt.Errorf("%s trop volumineux (maximum %.2f MB autorisés)",
			typeFichier, float64(maxSize)/(1024*1024))
	}

	// Générer un nom de fichier unique avec timestamp pour éviter les collisions
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	unique := uuid.New().String()
	path := filepath.Join(subDir, fmt.Sprintf("user_%d_%d_%s%s", userID, timestamp, unique, ext))
	log.Printf("📁 Chemin de destination: %s", path)

	// Ouvrir le fichier source
	src, err := f.Open()
	if err != nil {
		log.Printf("❌ Erreur ouverture fichier source: %v", err)
		return "", "", 0, fmt.Errorf("impossible d'ouvrir le fichier source: %v", err)
	}
	defer func(src multipart.File) {
		err := src.Close()
		if err != nil {

		}
	}(src)

	// Créer le fichier destination
	dst, err := os.Create(path)
	if err != nil {
		log.Printf("❌ Erreur création fichier destination: %v", err)
		return "", "", 0, fmt.Errorf("impossible de créer le fichier de destination: %v", err)
	}
	defer func(dst *os.File) {
		err := dst.Close()
		if err != nil {

		}
	}(dst)

	// Copier les données avec une vérification supplémentaire
	log.Printf("⏳ Copie des données en cours...")
	var bytesWritten int64
	buf := make([]byte, 32*1024) // Buffer de 32KB pour optimiser la copie
	for {
		n, readErr := src.Read(buf)
		if n > 0 {
			n, writeErr := dst.Write(buf[:n])
			if writeErr != nil {
				log.Printf("❌ Erreur d'écriture après %d bytes: %v", bytesWritten, writeErr)
				err := os.Remove(path)
				if err != nil {
					return "", "", 0, err
				}
				return "", "", 0, fmt.Errorf("erreur lors de l'écriture du fichier: %v", writeErr)
			}
			bytesWritten += int64(n)
		}
		if readErr != nil {
			if readErr == io.EOF {
				break
			}
			log.Printf("❌ Erreur de lecture après %d bytes: %v", bytesWritten, readErr)
			err := os.Remove(path)
			if err != nil {
				return "", "", 0, err
			}
			return "", "", 0, fmt.Errorf("erreur lors de la lecture du fichier: %v", readErr)
		}
	}

	// Vérifier que le nombre d'octets écrits correspond à la taille du fichier
	if bytesWritten != f.Size {
		log.Printf("⚠️ Avertissement: Taille du fichier écrit (%d) différente de la taille attendue (%d)",
			bytesWritten, f.Size)
	}

	log.Printf("✅ Fichier enregistré avec succès: %d bytes écrits", bytesWritten)
	return path, cleanFilename, bytesWritten, nil
}

// Fonction pour vérifier la sécurité du dossier uploads
func checkUploadsDirectorySecurity() error {
	uploadsDir := "uploads"

	// 1. Vérifier que le dossier existe
	info, err := os.Stat(uploadsDir)
	if os.IsNotExist(err) {
		log.Printf("📁 Le dossier uploads n'existe pas, tentative de création...")
		if err := os.MkdirAll(uploadsDir, 0750); err != nil {
			return fmt.Errorf("impossible de créer le dossier uploads: %v", err)
		}
		log.Printf("✅ Dossier uploads créé avec succès")
		return nil
	} else if err != nil {
		return fmt.Errorf("erreur lors de la vérification du dossier uploads: %v", err)
	}

	// 2. Vérifier que c'est bien un dossier
	if !info.IsDir() {
		return fmt.Errorf("uploads existe mais n'est pas un dossier")
	}

	// 3. Vérifier les permissions
	mode := info.Mode()
	log.Printf("📁 Dossier uploads avec permissions: %v", mode)

	// 4. Vérifier si on peut écrire dans le dossier
	testFile := filepath.Join(uploadsDir, "test_write_permission.tmp")
	f, err := os.Create(testFile)
	if err != nil {
		return fmt.Errorf("impossible d'écrire dans le dossier uploads: %v", err)
	}
	err = f.Close()
	if err != nil {
		return err
	}
	err = os.Remove(testFile)
	if err != nil {
		return err
	}
	log.Printf("✅ Permissions d'écriture dans uploads vérifiées")

	// 5. Vérifier les permissions sur les sous-dossiers
	subDirs := []string{"images", "videos", "documents", "autres", "temp"}
	for _, subDir := range subDirs {
		fullPath := filepath.Join(uploadsDir, subDir)
		if err := os.MkdirAll(fullPath, 0750); err != nil {
			return fmt.Errorf("impossible de créer le sous-dossier %s: %v", subDir, err)
		}
		log.Printf("✅ Sous-dossier %s vérifié", subDir)
	}

	return nil
}

// Analyser les informations d'un document
func getDocumentInfo(file *multipart.FileHeader) DocumentInfo {
	ext := strings.ToLower(filepath.Ext(file.Filename))

	info := DocumentInfo{
		FileSize:     file.Size,
		DocumentType: file.Header.Get("document_type"),
		Format:       ext,
		Category:     getDocumentCategory(file.Filename),
		IsPDF:        ext == ".pdf",
		IsBinary:     isDocumentBinary(file.Filename),
	}

	return info
}

// Déterminer la catégorie d'un document selon son extension
func getDocumentCategory(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))

	switch ext {
	case ".pdf":
		return "Document PDF"
	case ".doc", ".docx", ".odt", ".rtf":
		return "Traitement de texte"
	case ".ppt", ".pptx", ".odp":
		return "Présentation"
	case ".xls", ".xlsx", ".ods", ".csv":
		return "Tableur"
	case ".txt", ".md":
		return "Texte brut"
	case ".tex":
		return "Document LaTeX"
	case ".epub":
		return "Livre électronique"
	default:
		return "Document divers"
	}
}

// Vérifier si un document est binaire (non texte)
func isDocumentBinary(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))

	// Liste des formats texte courants
	textFormats := map[string]bool{
		".txt":  true,
		".csv":  true,
		".md":   true,
		".html": true,
		".xml":  true,
		".json": true,
		".tex":  true,
	}

	return !textFormats[ext]
}

// Valider que le document respecte les critères de qualité éducative
func validateEducationalDocumentQuality(file *multipart.FileHeader) (bool, string) {
	// 1. Vérifier la taille du fichier
	const maxDocSize = 20 * 1024 * 1024 // 20 MB
	if file.Size > maxDocSize {
		return false, fmt.Sprintf(
			"Le document dépasse la taille maximale recommandée de 20 MB (%.2f MB)",
			float64(file.Size)/(1024*1024))
	}

	// 2. Vérifier le format
	if !isValidDocument(file.Filename) {
		return false, fmt.Sprintf("Le format %s n'est pas pris en charge. Formats acceptés: %v",
			strings.ToLower(filepath.Ext(file.Filename)),
			getDocumentFormatList())
	}

	// 3. Préférer PDF pour les documents finaux
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".pdf" && file.Size > 5*1024*1024 {
		// Si ce n'est pas un PDF et que c'est relativement gros, suggérer la conversion
		log.Printf("ℹ️ Document non-PDF volumineux détecté: %s (%.2f MB)",
			file.Filename, float64(file.Size)/(1024*1024))
	}

	// 4. Vérifier le type MIME
	if !validateMimeType(file) {
		log.Printf("⚠️ Le type MIME du document ne correspond pas à l'extension: %s", file.Header.Get("document_type"))
		// On ne bloque pas l'upload mais on log un avertissement
	}

	return true, ""
}

// Vérifie si un fichier est une image valide (exporté)
func IsValidImage(filename string) bool {
	return isValidImage(filename)
}

// Vérifie si un fichier est une vidéo valide (exporté)
func IsValidVideo(filename string) bool {
	return isValidVideo(filename)
}

// Vérifie si un fichier est un document valide (exporté)
func IsValidDocument(filename string) bool {
	return isValidDocument(filename)
}

// Vérifie si un fichier est sous une taille maximale (exporté, version simple)
func IsUnderSize(f *multipart.FileHeader, max int64) bool {
	return f.Size <= max
}
