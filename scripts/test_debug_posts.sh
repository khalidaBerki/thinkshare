#!/bin/bash

echo "🔒 Test corrigé : Post payant + Message de verrouillage"
echo "======================================================="

API_BASE="http://4.178.177.89"

# Utiliser les tokens existants que nous savons fonctionner
CREATOR_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTIyNTkyMzIsInVzZXJfaWQiOjd9.KbGLGdWtv9qNVc8x9_H_84i8cTs1ztER3nBo9oyZBf0"
SUBSCRIBER_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTIyNTkyMzMsInVzZXJfaWQiOjh9.FgzkvYO41YmxtHBkr6FDY7cKzBhFddmd-jQUKIZTXHo"

echo "🧪 Test avec utilisateurs créés (user_id: 7 et 8)"
echo "------------------------------------------------"

# Test direct de création de post payant avec URL corrigée
echo ""
echo "💰 1. Création d'un post PAYANT (test direct)"
echo "--------------------------------------------"

# Essai avec curl en mode verbose pour voir le problème
curl -v -X POST "$API_BASE/api/posts" \
  -H "Authorization: Bearer $CREATOR_TOKEN" \
  -F 'content=🔥 CONTENU EXCLUSIF PREMIUM 🔥 - Ceci est réservé aux abonnés payants seulement !' \
  -F 'visibility=public' \
  -F 'is_paid_only=true' 2>&1 | head -20

echo ""
echo "🆓 2. Création d'un post GRATUIT"
echo "-------------------------------"

curl -v -X POST "$API_BASE/api/posts" \
  -H "Authorization: Bearer $CREATOR_TOKEN" \
  -F 'content=📢 Contenu gratuit pour tous !' \
  -F 'visibility=public' \
  -F 'is_paid_only=false' 2>&1 | head -20

echo ""
echo "📋 3. Vérification des posts du créateur (user_id=7)"
echo "---------------------------------------------------"

curl -s "$API_BASE/api/posts/user/7" \
  -H "Authorization: Bearer $CREATOR_TOKEN"
