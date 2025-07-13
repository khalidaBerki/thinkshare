#!/bin/bash

echo "🔒 Test complet : Post payant + Message de verrouillage + Intégration Stripe"
echo "============================================================================"

API_BASE="http://4.178.177.89"

# Étape 1 : Créer un nouveau créateur pour tester
echo "👤 1. Création d'un nouveau créateur pour les tests"
echo "---------------------------------------------------"

CREATOR_RESPONSE=$(curl -s -X POST "$API_BASE/register" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "creator_payant@test.com",
    "password": "password123",
    "username": "creator_payant",
    "name": "Créateur Payant",
    "firstName": "Créateur"
  }')

echo "Créateur créé: $CREATOR_RESPONSE"

# Connexion du créateur
CREATOR_LOGIN=$(curl -s -X POST "$API_BASE/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "creator_payant@test.com",
    "password": "password123"
  }')

CREATOR_TOKEN=$(echo $CREATOR_LOGIN | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
echo "Token créateur: $CREATOR_TOKEN"

# Étape 2 : Créer un abonné
echo ""
echo "👥 2. Création d'un abonné pour les tests"
echo "-----------------------------------------"

SUBSCRIBER_RESPONSE=$(curl -s -X POST "$API_BASE/register" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "abonne_test@test.com",
    "password": "password123",
    "username": "abonne_test",
    "name": "Abonné Test",
    "firstName": "Abonné"
  }')

echo "Abonné créé: $SUBSCRIBER_RESPONSE"

# Connexion de l'abonné
SUBSCRIBER_LOGIN=$(curl -s -X POST "$API_BASE/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "abonne_test@test.com",
    "password": "password123"
  }')

SUBSCRIBER_TOKEN=$(echo $SUBSCRIBER_LOGIN | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
echo "Token abonné: $SUBSCRIBER_TOKEN"

# Étape 3 : Créer un post payant
echo ""
echo "💰 3. Création d'un post PAYANT"
echo "-------------------------------"

POST_PAYANT_RESPONSE=$(curl -s -X POST "$API_BASE/api/posts" \
  -H "Authorization: Bearer $CREATOR_TOKEN" \
  -F 'content=🔥 CONTENU EXCLUSIF PREMIUM 🔥 - Ce contenu incroyable est réservé aux abonnés payants ! Vous découvrirez des secrets exclusifs ici !' \
  -F 'visibility=public' \
  -F 'is_paid_only=true')

echo "Réponse création post payant: $POST_PAYANT_RESPONSE"

# Extraire l'ID du post (si la création a réussi)
POST_ID=$(echo $POST_PAYANT_RESPONSE | grep -o '"id":[0-9]*' | cut -d':' -f2)
echo "ID du post payant: $POST_ID"

# Étape 4 : Créer un post gratuit pour comparaison
echo ""
echo "🆓 4. Création d'un post GRATUIT pour comparaison"
echo "------------------------------------------------"

POST_GRATUIT_RESPONSE=$(curl -s -X POST "$API_BASE/api/posts" \
  -H "Authorization: Bearer $CREATOR_TOKEN" \
  -F 'content=📢 Contenu gratuit accessible à tous ! Ceci est un aperçu de mes publications.' \
  -F 'visibility=public' \
  -F 'is_paid_only=false')

echo "Réponse création post gratuit: $POST_GRATUIT_RESPONSE"

# Extraire l'ID du post gratuit
POST_GRATUIT_ID=$(echo $POST_GRATUIT_RESPONSE | grep -o '"id":[0-9]*' | cut -d':' -f2)
echo "ID du post gratuit: $POST_GRATUIT_ID"

# Étape 5 : Vérifier les posts du créateur
echo ""
echo "📋 5. Vérification des posts créés"
echo "----------------------------------"

POSTS_CREATOR=$(curl -s "$API_BASE/api/posts/user/$(echo $CREATOR_LOGIN | grep -o '"user_id":[0-9]*' | cut -d':' -f2 || echo $CREATOR_LOGIN | grep -o '"id":[0-9]*' | cut -d':' -f2)" \
  -H "Authorization: Bearer $CREATOR_TOKEN")

echo "Posts du créateur: $POSTS_CREATOR"

# Étape 6 : Test d'accès SANS abonnement (message de verrouillage attendu)
echo ""
echo "🚫 6. Test d'accès au post payant SANS abonnement"
echo "------------------------------------------------"

if [ ! -z "$POST_ID" ]; then
    ACCESS_DENIED=$(curl -s "$API_BASE/api/posts/$POST_ID" \
      -H "Authorization: Bearer $SUBSCRIBER_TOKEN")
    
    echo "Réponse accès sans abonnement: $ACCESS_DENIED"
    
    # Vérifier le message de verrouillage
    if echo "$ACCESS_DENIED" | grep -q "réservé aux abonnés payants"; then
        echo "✅ MESSAGE DE VERROUILLAGE DÉTECTÉ !"
    else
        echo "❌ Message de verrouillage non trouvé"
    fi
    
    # Vérifier has_access = false
    if echo "$ACCESS_DENIED" | grep -q '"has_access":false'; then
        echo "✅ has_access = false (correct)"
    else
        echo "❌ has_access devrait être false"
    fi
else
    echo "❌ Impossible de tester - Post payant non créé"
fi

# Étape 7 : Test d'accès au post gratuit (accès normal attendu)
echo ""
echo "✅ 7. Test d'accès au post gratuit (contrôle)"
echo "--------------------------------------------"

if [ ! -z "$POST_GRATUIT_ID" ]; then
    ACCESS_GRATUIT=$(curl -s "$API_BASE/api/posts/$POST_GRATUIT_ID" \
      -H "Authorization: Bearer $SUBSCRIBER_TOKEN")
    
    echo "Réponse accès post gratuit: $ACCESS_GRATUIT"
    
    # Vérifier l'accès libre
    if echo "$ACCESS_GRATUIT" | grep -q '"has_access":true'; then
        echo "✅ Accès gratuit = true (correct)"
    else
        echo "❌ Accès gratuit devrait être true"
    fi
fi

echo ""
echo "📊 RÉSUMÉ DES TESTS"
echo "==================="
echo "1. Création créateur: $([ ! -z "$CREATOR_TOKEN" ] && echo "✅ OK" || echo "❌ ÉCHEC")"
echo "2. Création abonné: $([ ! -z "$SUBSCRIBER_TOKEN" ] && echo "✅ OK" || echo "❌ ÉCHEC")"
echo "3. Post payant créé: $([ ! -z "$POST_ID" ] && echo "✅ OK" || echo "❌ ÉCHEC")"
echo "4. Post gratuit créé: $([ ! -z "$POST_GRATUIT_ID" ] && echo "✅ OK" || echo "❌ ÉCHEC")"
echo "5. Message de verrouillage: $(echo "$ACCESS_DENIED" | grep -q "réservé aux abonnés" && echo "✅ OK" || echo "❌ ÉCHEC")"

echo ""
echo "🎯 PROCHAINE ÉTAPE : Test intégration Stripe"
echo "============================================"
