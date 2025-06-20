#!/bin/bash

# Script de débogage pour l'upload de vidéos

# Couleurs pour les messages
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# URL de base
BASE_URL="http://localhost:8080"

# Vérifier si curl est installé
if ! command -v curl &> /dev/null; then
    echo -e "${RED}Erreur: curl n'est pas installé. Veuillez l'installer pour utiliser ce script.${NC}"
    exit 1
 fi

# 1. Vérifier que le serveur est accessible
echo -e "${BLUE}Vérification de la connexion au serveur...${NC}"
SERVER_RESPONSE=$(curl -s -o /dev/null -w "%{http_code}" "${BASE_URL}")
if [ "$SERVER_RESPONSE" == "000" ]; then
    echo -e "${RED}Erreur: Impossible de se connecter au serveur. Assurez-vous qu'il est en cours d'exécution sur localhost:8080${NC}"
    exit 1
fi
echo -e "${GREEN}Serveur accessible!${NC}"

# 2. Vérifier le mode debug
echo -e "\n${BLUE}Vérification du mode debug...${NC}"
DEBUG_RESPONSE=$(curl -s "${BASE_URL}/debug/mode")
echo "$DEBUG_RESPONSE"

# 3. Obtenir un token JWT
echo -e "\n${BLUE}Obtention d'un token JWT...${NC}"
TOKEN_RESPONSE=$(curl -s -X POST "${BASE_URL}/api/fake-login")

# Extraire le token avec grep et cut (fonctionne sur la plupart des systèmes Unix)
TOKEN=$(echo "$TOKEN_RESPONSE" | grep -o '"token":"[^"]*"' | cut -d '"' -f 4)

if [ -z "$TOKEN" ]; then
    echo -e "${RED}Erreur: Impossible d'obtenir un token. Réponse: $TOKEN_RESPONSE${NC}"
    exit 1
fi

echo -e "${GREEN}Token obtenu avec succès!${NC}"

# 4. Demander le chemin du fichier vidéo
echo -e "\n${YELLOW}Entrez le chemin vers le fichier vidéo:${NC}"
read -p "> " VIDEO_PATH

if [ ! -f "$VIDEO_PATH" ]; then
    echo -e "${RED}Erreur: Le fichier n'existe pas!${NC}"
    exit 1
fi

# 5. Analyser le fichier vidéo
echo -e "\n${BLUE}Analyse du fichier vidéo...${NC}"
FILE_SIZE=$(stat -f%z "$VIDEO_PATH" 2>/dev/null || stat -c%s "$VIDEO_PATH" 2>/dev/null)
FILE_SIZE_MB=$(echo "scale=2; $FILE_SIZE / 1048576" | bc)
FILE_EXT=$(echo "$VIDEO_PATH" | awk -F. '{print tolower($NF)}')

echo -e "Nom du fichier: $(basename "$VIDEO_PATH")"
echo -e "Taille: ${FILE_SIZE_MB} MB"
echo -e "Extension: .${FILE_EXT}"

# Vérifier l'extension
if [[ "$FILE_EXT" != "mp4" && "$FILE_EXT" != "mov" && "$FILE_EXT" != "webm" ]]; then
    echo -e "${RED}Attention: L'extension '${FILE_EXT}' n'est pas dans la liste des formats acceptés (.mp4, .mov, .webm)${NC}"
fi

# Vérifier la taille
if (( $(echo "$FILE_SIZE_MB > 100" | bc -l) )); then
    echo -e "${RED}Attention: La taille du fichier (${FILE_SIZE_MB} MB) dépasse la limite de 100 MB${NC}"
fi

# 6. Upload de la vidéo
echo -e "\n${BLUE}Upload de la vidéo en cours...${NC}"
echo -e "Cela peut prendre du temps en fonction de la taille du fichier.\n"

UPLOAD_RESPONSE=$(curl -s -X POST \
    -H "Authorization: Bearer $TOKEN" \
    -F "content=Test depuis script debug" \
    -F "visibility=public" \
    -F "video=@$VIDEO_PATH" \
    "${BASE_URL}/api/posts")

# Vérifier si la réponse contient une erreur
if echo "$UPLOAD_RESPONSE" | grep -q "error"; then
    echo -e "${RED}L'upload a échoué!${NC}"
    echo "$UPLOAD_RESPONSE"
    exit 1
else
    echo -e "${GREEN}Upload réussi!${NC}"
    echo "$UPLOAD_RESPONSE"
fi

echo -e "\n${GREEN}Test terminé avec succès!${NC}"
