#!/bin/bash

# Script pour récupérer, créer et gérer des posts via l'API

# Couleurs pour les messages
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# URL de base
BASE_URL="http://localhost:8080"

# Variables
TOKEN=""

show_menu() {
    clear
    echo -e "${BLUE}=============================================${NC}"
    echo -e "${BLUE}         GESTIONNAIRE DE POSTS              ${NC}"
    echo -e "${BLUE}=============================================${NC}"
    echo -e "1. Obtenir un token"
    echo -e "2. Afficher tous les posts"
    echo -e "3. Créer un post texte"
    echo -e "4. Créer un post avec image"
    echo -e "5. Créer un post avec vidéo"
    echo -e "6. Supprimer un post"
    echo -e "7. Afficher le token actuel"
    echo -e "8. Tester l'authentification"
    echo -e "0. Quitter"
    echo -e "${BLUE}=============================================${NC}"
    echo -e "Token actuel: ${TOKEN:0:15}..."
    echo -e "${BLUE}=============================================${NC}"
}

get_token() {
    echo -e "\n${BLUE}Obtention d'un token JWT...${NC}"
    TOKEN_RESPONSE=$(curl -s -X POST "${BASE_URL}/api/fake-login")

    # Extraire le token avec grep et cut
    TOKEN=$(echo "$TOKEN_RESPONSE" | grep -o '"token":"[^"]*"' | cut -d '"' -f 4)

    if [ -z "$TOKEN" ]; then
        echo -e "${RED}Erreur: Impossible d'obtenir un token. Réponse: $TOKEN_RESPONSE${NC}"
        return 1
    fi

    echo -e "${GREEN}Token obtenu avec succès!${NC}"
    return 0
}

show_posts() {
    if [ -z "$TOKEN" ]; then
        echo -e "${RED}Erreur: Aucun token disponible. Veuillez d'abord obtenir un token.${NC}"
        return 1
    fi

    # Options de pagination
    read -p "Page (défaut: 1): " PAGE
    read -p "Taille de page (défaut: 10): " PAGE_SIZE
    read -p "Visibilité (public/private/laisser vide pour tous): " VISIBILITY

    PAGE=${PAGE:-1}
    PAGE_SIZE=${PAGE_SIZE:-10}

    PARAMS="page=${PAGE}&pageSize=${PAGE_SIZE}"
    if [ ! -z "$VISIBILITY" ]; then
        PARAMS="${PARAMS}&visibility=${VISIBILITY}"
    fi

    echo -e "\n${BLUE}Récupération des posts...${NC}"
    POSTS_RESPONSE=$(curl -s -X GET \
        -H "Authorization: Bearer $TOKEN" \
        "${BASE_URL}/api/posts?${PARAMS}")

    # Vérifier si la réponse contient une erreur
    if echo "$POSTS_RESPONSE" | grep -q "error"; then
        echo -e "${RED}Erreur: $POSTS_RESPONSE${NC}"
        return 1
    fi

    # Vérifier si la liste est vide
    if [ "$POSTS_RESPONSE" = "[]" ] || [ "$POSTS_RESPONSE" = "null" ]; then
        echo -e "${YELLOW}Aucun post trouvé.${NC}"
        return 0
    fi

    # Formater et afficher les posts
    echo -e "${GREEN}Posts récupérés:${NC}"

    # Analyser le JSON pour extraire les données importantes
    echo "$POSTS_RESPONSE" | python3 -c "import sys, json; \
        data = json.load(sys.stdin); \
        for post in data: \
            print(f\"\\n{'-' * 50}\\nID: {post['ID']}\\nDate: {post['CreatedAt']}\\nVisibilité: {post['Visibility']}\\nContenu: {post['Content']}\\nMédia: {len(post.get('Media', []))} fichier(s)\")" || \
    echo "$POSTS_RESPONSE" | python -c "import sys, json; \
        data = json.load(sys.stdin); \
        for post in data: \
            print(f\"\\n{'-' * 50}\\nID: {post['ID']}\\nDate: {post['CreatedAt']}\\nVisibilité: {post['Visibility']}\\nContenu: {post['Content']}\\nMédia: {len(post.get('Media', []))} fichier(s)\")" || \
    echo "$POSTS_RESPONSE"

    echo -e "\n${BLUE}Total: $(echo "$POSTS_RESPONSE" | grep -o "ID" | wc -l) posts${NC}"
    return 0
}

create_text_post() {
    if [ -z "$TOKEN" ]; then
        echo -e "${RED}Erreur: Aucun token disponible. Veuillez d'abord obtenir un token.${NC}"
        return 1
    fi

    read -p "Contenu du post: " CONTENT
    read -p "Visibilité (public/private): " VISIBILITY

    if [ -z "$CONTENT" ]; then
        echo -e "${RED}Erreur: Le contenu ne peut pas être vide.${NC}"
        return 1
    fi

    if [ "$VISIBILITY" != "public" ] && [ "$VISIBILITY" != "private" ]; then
        echo -e "${RED}Erreur: La visibilité doit être 'public' ou 'private'.${NC}"
        return 1
    fi

    echo -e "\n${BLUE}Création du post...${NC}"
    CREATE_RESPONSE=$(curl -s -X POST \
        -H "Authorization: Bearer $TOKEN" \
        -F "content=$CONTENT" \
        -F "visibility=$VISIBILITY" \
        "${BASE_URL}/api/posts")

    # Vérifier si la réponse contient une erreur
    if echo "$CREATE_RESPONSE" | grep -q "error"; then
        echo -e "${RED}Erreur: $CREATE_RESPONSE${NC}"
        return 1
    else
        echo -e "${GREEN}Post créé avec succès!${NC}"
        echo "$CREATE_RESPONSE"
        return 0
    fi
}

create_image_post() {
    if [ -z "$TOKEN" ]; then
        echo -e "${RED}Erreur: Aucun token disponible. Veuillez d'abord obtenir un token.${NC}"
        return 1
    fi

    read -p "Chemin vers l'image: " IMAGE_PATH
    read -p "Contenu du post: " CONTENT
    read -p "Visibilité (public/private): " VISIBILITY

    if [ ! -f "$IMAGE_PATH" ]; then
        echo -e "${RED}Erreur: Le fichier n'existe pas!${NC}"
        return 1
    fi

    echo -e "\n${BLUE}Création du post avec image...${NC}"
    CREATE_RESPONSE=$(curl -s -X POST \
        -H "Authorization: Bearer $TOKEN" \
        -F "content=$CONTENT" \
        -F "visibility=$VISIBILITY" \
        -F "images=@$IMAGE_PATH" \
        "${BASE_URL}/api/posts")

    # Vérifier si la réponse contient une erreur
    if echo "$CREATE_RESPONSE" | grep -q "error"; then
        echo -e "${RED}Erreur: $CREATE_RESPONSE${NC}"
        return 1
    else
        echo -e "${GREEN}Post avec image créé avec succès!${NC}"
        echo "$CREATE_RESPONSE"
        return 0
    fi
}

create_video_post() {
    if [ -z "$TOKEN" ]; then
        echo -e "${RED}Erreur: Aucun token disponible. Veuillez d'abord obtenir un token.${NC}"
        return 1
    fi

    read -p "Chemin vers la vidéo: " VIDEO_PATH
    read -p "Contenu du post: " CONTENT
    read -p "Visibilité (public/private): " VISIBILITY

    if [ ! -f "$VIDEO_PATH" ]; then
        echo -e "${RED}Erreur: Le fichier n'existe pas!${NC}"
        return 1
    fi

    echo -e "\n${BLUE}Création du post avec vidéo...${NC}"
    CREATE_RESPONSE=$(curl -s -X POST \
        -H "Authorization: Bearer $TOKEN" \
        -F "content=$CONTENT" \
        -F "visibility=$VISIBILITY" \
        -F "video=@$VIDEO_PATH" \
        "${BASE_URL}/api/posts")

    # Vérifier si la réponse contient une erreur
    if echo "$CREATE_RESPONSE" | grep -q "error"; then
        echo -e "${RED}Erreur: $CREATE_RESPONSE${NC}"
        return 1
    else
        echo -e "${GREEN}Post avec vidéo créé avec succès!${NC}"
        echo "$CREATE_RESPONSE"
        return 0
    fi
}

delete_post() {
    if [ -z "$TOKEN" ]; then
        echo -e "${RED}Erreur: Aucun token disponible. Veuillez d'abord obtenir un token.${NC}"
        return 1
    fi

    read -p "ID du post à supprimer: " POST_ID

    if [ -z "$POST_ID" ]; then
        echo -e "${RED}Erreur: L'ID du post ne peut pas être vide.${NC}"
        return 1
    fi

    echo -e "\n${BLUE}Suppression du post...${NC}"
    DELETE_RESPONSE=$(curl -s -X DELETE \
        -H "Authorization: Bearer $TOKEN" \
        "${BASE_URL}/api/posts/$POST_ID" -w "\n%{http_code}")

    HTTP_CODE=$(echo "$DELETE_RESPONSE" | tail -n1)
    RESPONSE_BODY=$(echo "$DELETE_RESPONSE" | sed '$d')

    if [ "$HTTP_CODE" == "204" ]; then
        echo -e "${GREEN}Post supprimé avec succès!${NC}"
        return 0
    else
        echo -e "${RED}Erreur (HTTP $HTTP_CODE): $RESPONSE_BODY${NC}"
        return 1
    fi
}

show_token() {
    if [ -z "$TOKEN" ]; then
        echo -e "${RED}Aucun token disponible. Veuillez d'abord obtenir un token.${NC}"
    else
        echo -e "${GREEN}Token actuel:${NC}"
        echo "$TOKEN"
    fi
}

test_auth() {
    if [ -z "$TOKEN" ]; then
        echo -e "${RED}Erreur: Aucun token disponible. Veuillez d'abord obtenir un token.${NC}"
        return 1
    fi

    echo -e "\n${BLUE}Test d'authentification...${NC}"
    AUTH_RESPONSE=$(curl -s -X GET \
        -H "Authorization: Bearer $TOKEN" \
        "${BASE_URL}/api/test-auth")

    # Vérifier si la réponse contient une erreur
    if echo "$AUTH_RESPONSE" | grep -q "error"; then
        echo -e "${RED}Erreur d'authentification: $AUTH_RESPONSE${NC}"
        return 1
    else
        echo -e "${GREEN}Authentification réussie!${NC}"
        echo "$AUTH_RESPONSE"
        return 0
    fi
}

# Programme principal
while true; do
    show_menu
    read -p "Entrez votre choix: " CHOICE

    case $CHOICE in
        1) get_token; read -p "Appuyez sur Entrée pour continuer..." ;;
        2) show_posts; read -p "Appuyez sur Entrée pour continuer..." ;;
        3) create_text_post; read -p "Appuyez sur Entrée pour continuer..." ;;
        4) create_image_post; read -p "Appuyez sur Entrée pour continuer..." ;;
        5) create_video_post; read -p "Appuyez sur Entrée pour continuer..." ;;
        6) delete_post; read -p "Appuyez sur Entrée pour continuer..." ;;
        7) show_token; read -p "Appuyez sur Entrée pour continuer..." ;;
        8) test_auth; read -p "Appuyez sur Entrée pour continuer..." ;;
        0) echo -e "${GREEN}Au revoir!${NC}"; exit 0 ;;
        *) echo -e "${RED}Choix invalide${NC}"; read -p "Appuyez sur Entrée pour continuer..." ;;
    esac
done
