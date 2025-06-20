#!/bin/bash

# Script pour vérifier les permissions du dossier uploads

# Couleurs pour les messages
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}Vérification des permissions du dossier uploads${NC}"

# 1. Vérifier si le dossier existe
if [ ! -d "uploads" ]; then
    echo -e "${YELLOW}Le dossier uploads n'existe pas, tentative de création...${NC}"
    mkdir -p uploads
    if [ $? -ne 0 ]; then
        echo -e "${RED}ERREUR: Impossible de créer le dossier uploads${NC}"
        exit 1
    else
        echo -e "${GREEN}Dossier uploads créé avec succès${NC}"
    fi
else
    echo -e "${GREEN}Le dossier uploads existe${NC}"
fi

# 2. Vérifier les permissions du dossier
ls -la | grep uploads
echo ""

# 3. Tester l'écriture dans le dossier
echo -e "${BLUE}Test d'écriture dans le dossier uploads${NC}"
TEST_FILE="uploads/test_$(date +%s).txt"
echo "Test d'écriture" > "$TEST_FILE"
if [ $? -ne 0 ]; then
    echo -e "${RED}ERREUR: Impossible d'écrire dans le dossier uploads${NC}"
    exit 1
else
    echo -e "${GREEN}Écriture dans le dossier uploads réussie${NC}"
    rm "$TEST_FILE"
fi

# 4. Corriger les permissions si nécessaire
echo -e "${BLUE}Définition des permissions correctes${NC}"
chmod -R 755 uploads
if [ $? -ne 0 ]; then
    echo -e "${YELLOW}Note: Impossible de modifier les permissions. Essayez avec sudo.${NC}"
else
    echo -e "${GREEN}Permissions correctement définies${NC}"
fi

# 5. Afficher les permissions actuelles
echo -e "${BLUE}Permissions actuelles:${NC}"
ls -la | grep uploads

# 6. Tester avec un petit fichier vidéo
echo -e "\n${BLUE}Test d'écriture d'un petit fichier vidéo de test${NC}"
TEST_VIDEO="uploads/test_video_$(date +%s).mp4"

# Création d'un petit fichier binaire qui ressemble à un en-tête MP4
dd if=/dev/urandom of="$TEST_VIDEO" bs=1024 count=10 2>/dev/null

if [ $? -ne 0 ] || [ ! -f "$TEST_VIDEO" ]; then
    echo -e "${RED}ERREUR: Impossible de créer un fichier vidéo de test${NC}"
    exit 1
else
    echo -e "${GREEN}Fichier vidéo de test créé avec succès: $TEST_VIDEO${NC}"
    echo -e "${YELLOW}Taille: $(du -h "$TEST_VIDEO" | cut -f1)${NC}"

    # Supprimer le fichier de test
    rm "$TEST_VIDEO"
    echo -e "${GREEN}Fichier de test supprimé${NC}"
fi

echo -e "\n${GREEN}Vérification terminée. Le dossier uploads semble être correctement configuré.${NC}"
