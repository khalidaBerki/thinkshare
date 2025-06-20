#!/bin/bash

# Script pour vérifier et nettoyer les fichiers médias orphelins

# Couleurs pour les messages
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Variables
UPLOADS_DIR="uploads"
TOKEN=""
BASE_URL="http://localhost:8080"

echo -e "${BLUE}=============================================${NC}"
echo -e "${BLUE}      OUTIL DE VÉRIFICATION DES MÉDIAS      ${NC}"
echo -e "${BLUE}=============================================${NC}"

# Vérifier si le dossier uploads existe
if [ ! -d "$UPLOADS_DIR" ]; then
    echo -e "${YELLOW}Le dossier $UPLOADS_DIR n'existe pas. Création...${NC}"
    mkdir -p "$UPLOADS_DIR"
    echo -e "${GREEN}Dossier $UPLOADS_DIR créé.${NC}"
fi

# Compter les fichiers dans le dossier uploads
NB_FILES=$(ls -1 "$UPLOADS_DIR" 2>/dev/null | wc -l)

echo -e "${BLUE}Nombre de fichiers dans $UPLOADS_DIR: $NB_FILES${NC}"

# Obtenir un token pour les requêtes API
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

# Récupérer tous les posts et leurs médias
get_all_posts() {
    if [ -z "$TOKEN" ]; then
        echo -e "${RED}Erreur: Aucun token disponible. Veuillez d'abord obtenir un token.${NC}"
        return 1
    fi

    echo -e "\n${BLUE}Récupération de tous les posts...${NC}"

    # Récupérer le nombre total approximatif de posts (on prend une grande page)
    POSTS_COUNT_RESPONSE=$(curl -s -X GET \
        -H "Authorization: Bearer $TOKEN" \
        "${BASE_URL}/api/posts?pageSize=1000")

    # Vérifier si la réponse contient une erreur
    if echo "$POSTS_COUNT_RESPONSE" | grep -q "error"; then
        echo -e "${RED}Erreur: $POSTS_COUNT_RESPONSE${NC}"
        return 1
    fi

    # Compter les posts
    POSTS_COUNT=$(echo "$POSTS_COUNT_RESPONSE" | grep -o "ID" | wc -l)
    echo -e "${GREEN}Nombre total de posts: $POSTS_COUNT${NC}"

    # Récupérer tous les médias des posts
    echo "$POSTS_COUNT_RESPONSE" > all_posts.json
    echo -e "${GREEN}Données des posts sauvegardées dans all_posts.json${NC}"

    # Extraire les URLs des médias
    echo -e "\n${BLUE}Extraction des URLs des médias...${NC}"

    # Utiliser python si disponible pour extraire les URLs
    if command -v python3 &> /dev/null; then
        echo "$POSTS_COUNT_RESPONSE" | python3 -c "import sys, json; \
            data = json.load(sys.stdin); \
            media_urls = [m['MediaURL'] for post in data for m in post.get('Media', [])]; \
            print('\n'.join(media_urls))" > media_urls.txt
    elif command -v python &> /dev/null; then
        echo "$POSTS_COUNT_RESPONSE" | python -c "import sys, json; \
            data = json.load(sys.stdin); \
            media_urls = [m['MediaURL'] for post in data for m in post.get('Media', [])]; \
            print('\n'.join(media_urls))" > media_urls.txt
    else
        echo -e "${YELLOW}Python non trouvé. Extraction manuelle...${NC}"
        # Extraire les URLs avec grep (moins précis)
        echo "$POSTS_COUNT_RESPONSE" | grep -o '"MediaURL":"[^"]*"' | cut -d '"' -f 4 > media_urls.txt
    fi

    MEDIA_COUNT=$(wc -l < media_urls.txt)
    echo -e "${GREEN}$MEDIA_COUNT URLs de médias extraites et sauvegardées dans media_urls.txt${NC}"

    return 0
}

# Trouver les fichiers orphelins (fichiers dans uploads qui ne sont pas référencés dans les posts)
find_orphaned_files() {
    if [ ! -f "media_urls.txt" ]; then
        echo -e "${RED}Erreur: Fichier media_urls.txt non trouvé. Exécutez d'abord get_all_posts.${NC}"
        return 1
    fi

    echo -e "\n${BLUE}Recherche des fichiers orphelins...${NC}"

    # Créer une liste de tous les fichiers dans uploads
    find "$UPLOADS_DIR" -type f | sort > all_files.txt
    FILE_COUNT=$(wc -l < all_files.txt)
    echo -e "${BLUE}$FILE_COUNT fichiers trouvés dans le dossier $UPLOADS_DIR${NC}"

    # Comparer les fichiers avec les URLs des médias
    echo -e "${BLUE}Comparaison avec les médias référencés...${NC}"

    # Créer un fichier pour les orphelins
    > orphaned_files.txt

    while IFS= read -r file; do
        # Vérifier si le fichier est référencé dans media_urls.txt
        if ! grep -q "$file" media_urls.txt; then
            echo "$file" >> orphaned_files.txt
        fi
    done < all_files.txt

    ORPHAN_COUNT=$(wc -l < orphaned_files.txt)
    echo -e "${GREEN}$ORPHAN_COUNT fichiers orphelins trouvés et listés dans orphaned_files.txt${NC}"

    # Afficher les 10 premiers fichiers orphelins
    if [ "$ORPHAN_COUNT" -gt 0 ]; then
        echo -e "\n${YELLOW}Exemples de fichiers orphelins:${NC}"
        head -n 10 orphaned_files.txt

        if [ "$ORPHAN_COUNT" -gt 10 ]; then
            echo -e "${YELLOW}... et $(($ORPHAN_COUNT - 10)) autres${NC}"
        fi
    fi

    return 0
}

# Nettoyer les fichiers orphelins
clean_orphaned_files() {
    if [ ! -f "orphaned_files.txt" ]; then
        echo -e "${RED}Erreur: Fichier orphaned_files.txt non trouvé. Exécutez d'abord find_orphaned_files.${NC}"
        return 1
    fi

    ORPHAN_COUNT=$(wc -l < orphaned_files.txt)

    if [ "$ORPHAN_COUNT" -eq 0 ]; then
        echo -e "${GREEN}Aucun fichier orphelin à nettoyer.${NC}"
        return 0
    fi

    read -p "Voulez-vous supprimer $ORPHAN_COUNT fichiers orphelins? (o/n): " CONFIRM

    if [ "$CONFIRM" != "o" ] && [ "$CONFIRM" != "O" ]; then
        echo -e "${YELLOW}Nettoyage annulé.${NC}"
        return 0
    fi

    echo -e "\n${BLUE}Suppression des fichiers orphelins...${NC}"

    # Créer un répertoire de sauvegarde
    BACKUP_DIR="media_backup_$(date +%Y%m%d_%H%M%S)"
    mkdir -p "$BACKUP_DIR"
    echo -e "${BLUE}Sauvegarde des fichiers dans $BACKUP_DIR avant suppression${NC}"

    # Compteurs
    DELETED=0
    BACKUP=0
    ERROR=0

    while IFS= read -r file; do
        # Créer une copie de sauvegarde
        cp "$file" "$BACKUP_DIR/" 2>/dev/null
        if [ $? -eq 0 ]; then
            BACKUP=$((BACKUP+1))
        fi

        # Supprimer le fichier
        rm "$file" 2>/dev/null
        if [ $? -eq 0 ]; then
            DELETED=$((DELETED+1))
            echo -e "${GREEN}Supprimé: $file${NC}"
        else
            ERROR=$((ERROR+1))
            echo -e "${RED}Erreur lors de la suppression: $file${NC}"
        fi
    done < orphaned_files.txt

    echo -e "\n${GREEN}Résultat du nettoyage:${NC}"
    echo -e "${GREEN}- $DELETED fichiers supprimés${NC}"
    echo -e "${GREEN}- $BACKUP fichiers sauvegardés dans $BACKUP_DIR${NC}"
    if [ "$ERROR" -gt 0 ]; then
        echo -e "${RED}- $ERROR erreurs rencontrées${NC}"
    fi

    return 0
}

# Menu principal
while true; do
    echo -e "\n${BLUE}=============================================${NC}"
    echo -e "${BLUE}                MENU                      ${NC}"
    echo -e "${BLUE}=============================================${NC}"
    echo -e "1. Obtenir un token"
    echo -e "2. Récupérer tous les posts et leurs médias"
    echo -e "3. Trouver les fichiers médias orphelins"
    echo -e "4. Nettoyer les fichiers orphelins"
    echo -e "5. Vérifier les permissions du dossier uploads"
    echo -e "0. Quitter"
    echo -e "${BLUE}=============================================${NC}"

    read -p "Entrez votre choix: " CHOICE

    case $CHOICE in
        1) get_token ;;  
        2) get_all_posts ;;  
        3) find_orphaned_files ;;  
        4) clean_orphaned_files ;;  
        5) 
            echo -e "\n${BLUE}Vérification des permissions du dossier $UPLOADS_DIR...${NC}"
            ls -ld "$UPLOADS_DIR"
            echo -e "\n${BLUE}Tentative de création d'un fichier test...${NC}"
            if touch "$UPLOADS_DIR/test_permission" 2>/dev/null; then
                echo -e "${GREEN}Test réussi: Le dossier est accessible en écriture${NC}"
                rm "$UPLOADS_DIR/test_permission"
            else
                echo -e "${RED}Test échoué: Le dossier n'est pas accessible en écriture${NC}"
            fi
            ;;
        0) echo -e "${GREEN}Au revoir!${NC}"; exit 0 ;;  
        *) echo -e "${RED}Choix invalide${NC}" ;;  
    esac

    read -p "Appuyez sur Entrée pour continuer..."
done
