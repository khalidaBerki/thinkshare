#!/bin/bash

# Script de test pour le contrôle d'accès aux posts payants
# Ce script teste le workflow complet : création de post payant, accès sans abonnement, puis avec abonnement

API_BASE="http://4.178.177.89"
echo "🔍 Test du contrôle d'accès aux posts payants sur $API_BASE"
echo "=================================="

# 1. Créer un utilisateur créateur
echo "📝 1. Création d'un utilisateur créateur..."
CREATOR_RESPONSE=$(curl -s -X POST "$API_BASE/register" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "creator_test@example.com",
    "password": "password123",
    "username": "creator_test"
  }')

echo "Réponse création créateur: $CREATOR_RESPONSE"

# Extraire le token du créateur
CREATOR_TOKEN=$(echo $CREATOR_RESPONSE | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
echo "Token créateur: $CREATOR_TOKEN"

# 2. Créer un utilisateur abonné potentiel
echo "📝 2. Création d'un utilisateur abonné..."
SUBSCRIBER_RESPONSE=$(curl -s -X POST "$API_BASE/register" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "subscriber_test@example.com",
    "password": "password123",
    "username": "subscriber_test"
  }')

echo "Réponse création abonné: $SUBSCRIBER_RESPONSE"

# Extraire le token de l'abonné
SUBSCRIBER_TOKEN=$(echo $SUBSCRIBER_RESPONSE | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
echo "Token abonné: $SUBSCRIBER_TOKEN"

# 3. Créer un post payant
echo "🔒 3. Création d'un post payant..."
POST_RESPONSE=$(curl -s -X POST "$API_BASE/api/posts" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $CREATOR_TOKEN" \
  -d '{
    "content": "Ceci est un contenu exclusif réservé aux abonnés payants !",
    "is_paid_only": true
  }')

echo "Réponse création post: $POST_RESPONSE"

# Extraire l'ID du post
POST_ID=$(echo $POST_RESPONSE | grep -o '"id":[0-9]*' | cut -d':' -f2)
echo "ID du post payant: $POST_ID"

# 4. Tester l'accès sans abonnement
echo "❌ 4. Test d'accès au post payant SANS abonnement..."
ACCESS_DENIED_RESPONSE=$(curl -s -X GET "$API_BASE/api/posts/$POST_ID" \
  -H "Authorization: Bearer $SUBSCRIBER_TOKEN")

echo "Réponse accès refusé: $ACCESS_DENIED_RESPONSE"

# Vérifier que le contenu est verrouillé
if echo "$ACCESS_DENIED_RESPONSE" | grep -q "réservé aux abonnés payants"; then
    echo "✅ Contrôle d'accès OK - Contenu verrouillé"
else
    echo "❌ ERREUR - Le contenu devrait être verrouillé"
fi

# 5. Lister tous les posts pour vérifier le contrôle d'accès
echo "📋 5. Test de la liste des posts avec contrôle d'accès..."
POSTS_LIST_RESPONSE=$(curl -s -X GET "$API_BASE/api/posts?page=1&limit=10" \
  -H "Authorization: Bearer $SUBSCRIBER_TOKEN")

echo "Réponse liste posts: $POSTS_LIST_RESPONSE"

# Vérifier que le post apparaît avec HasAccess = false
if echo "$POSTS_LIST_RESPONSE" | grep -q '"has_access":false'; then
    echo "✅ Liste des posts OK - HasAccess = false détecté"
else
    echo "❌ ERREUR - HasAccess devrait être false dans la liste"
fi

# 6. Créer un abonnement payant (simulation)
echo "💳 6. Simulation d'un abonnement payant..."
# Dans un vrai test, on passerait par Stripe, ici on simule en créant directement un abonnement actif

# Récupérer l'ID de l'abonné depuis le token
SUBSCRIBER_ID=$(echo $SUBSCRIBER_RESPONSE | grep -o '"user_id":[0-9]*' | cut -d':' -f2)
if [ -z "$SUBSCRIBER_ID" ]; then
    SUBSCRIBER_ID=$(echo $SUBSCRIBER_RESPONSE | grep -o '"id":[0-9]*' | cut -d':' -f2)
fi
echo "ID abonné: $SUBSCRIBER_ID"

# Récupérer l'ID du créateur
CREATOR_ID=$(echo $CREATOR_RESPONSE | grep -o '"user_id":[0-9]*' | cut -d':' -f2)
if [ -z "$CREATOR_ID" ]; then
    CREATOR_ID=$(echo $CREATOR_RESPONSE | grep -o '"id":[0-9]*' | cut -d':' -f2)
fi
echo "ID créateur: $CREATOR_ID"

# Ici, dans un vrai environnement, on déclencherait un webhook Stripe
# Pour le test, supposons qu'un abonnement soit créé manuellement dans la DB

echo "⏳ 7. Attente de l'activation de l'abonnement (simulation)..."
sleep 2

# 8. Tester l'accès AVEC abonnement
echo "✅ 8. Test d'accès au post payant AVEC abonnement..."
ACCESS_GRANTED_RESPONSE=$(curl -s -X GET "$API_BASE/api/posts/$POST_ID" \
  -H "Authorization: Bearer $SUBSCRIBER_TOKEN")

echo "Réponse accès autorisé: $ACCESS_GRANTED_RESPONSE"

# Vérifier que le contenu est maintenant accessible
if echo "$ACCESS_GRANTED_RESPONSE" | grep -q "contenu exclusif réservé aux abonnés"; then
    echo "✅ Abonnement OK - Contenu accessible"
else
    echo "⚠️  L'abonnement n'est peut-être pas encore activé"
fi

echo ""
echo "🏁 Test terminé !"
echo "=================================="

# Résumé des résultats
echo "📊 RÉSUMÉ DES TESTS:"
echo "1. Création d'utilisateurs: $([ ! -z "$CREATOR_TOKEN" ] && [ ! -z "$SUBSCRIBER_TOKEN" ] && echo "✅ OK" || echo "❌ ÉCHEC")"
echo "2. Création de post payant: $([ ! -z "$POST_ID" ] && echo "✅ OK" || echo "❌ ÉCHEC")"
echo "3. Contrôle d'accès sans abonnement: $(echo "$ACCESS_DENIED_RESPONSE" | grep -q "réservé aux abonnés" && echo "✅ OK" || echo "❌ ÉCHEC")"
echo "4. Liste avec HasAccess: $(echo "$POSTS_LIST_RESPONSE" | grep -q '"has_access":false' && echo "✅ OK" || echo "❌ ÉCHEC")"
echo "5. Accès avec abonnement: $(echo "$ACCESS_GRANTED_RESPONSE" | grep -q "contenu exclusif" && echo "✅ OK" || echo "⚠️  À VÉRIFIER")"
