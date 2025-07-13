#!/bin/bash

# Test complet du contrôle d'accès aux posts payants
echo "🧪 Test du système de contrôle d'accès aux posts payants"
echo "======================================================="

API_BASE="http://4.178.177.89"

# Tokens des utilisateurs existants
USER5_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTIyNTg2NTAsInVzZXJfaWQiOjV9.JPehDp5zB6wZIUmsmh1eW4pgTM4fqMvvM2-Pub9A6TY"
USER6_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTIyNTg3MjUsInVzZXJfaWQiOjZ9.b3znW5wP1WmnRxp25ntqi93ZRILuUrv2EBG884TDvSk"

echo "📋 1. Test de la liste des posts avec User5 (creator_test4)"
echo "-----------------------------------------------------------"
POSTS_USER5=$(curl -s "$API_BASE/api/posts?limit=3" -H "Authorization: Bearer $USER5_TOKEN")
echo "Réponse User5: $POSTS_USER5"

echo ""
echo "📋 2. Test de la liste des posts avec User6 (subscriber)"
echo "-------------------------------------------------------"
POSTS_USER6=$(curl -s "$API_BASE/api/posts?limit=3" -H "Authorization: Bearer $USER6_TOKEN")
echo "Réponse User6: $POSTS_USER6"

echo ""
echo "🔍 3. Test d'accès à un post spécifique (Post 38)"
echo "------------------------------------------------"
echo "Accès avec User5 (creator_test4):"
POST38_USER5=$(curl -s "$API_BASE/api/posts/38" -H "Authorization: Bearer $USER5_TOKEN")
echo "$POST38_USER5"

echo ""
echo "Accès avec User6 (subscriber):"
POST38_USER6=$(curl -s "$API_BASE/api/posts/38" -H "Authorization: Bearer $USER6_TOKEN")
echo "$POST38_USER6"

echo ""
echo "🔍 4. Test d'accès à un autre post (Post 37)"
echo "--------------------------------------------"
echo "Accès avec User6 (subscriber):"
POST37_USER6=$(curl -s "$API_BASE/api/posts/37" -H "Authorization: Bearer $USER6_TOKEN")
echo "$POST37_USER6"

echo ""
echo "📊 5. Analyse des résultats"
echo "---------------------------"

# Vérifier que les champs is_paid_only et has_access sont présents
if echo "$POSTS_USER6" | grep -q '"is_paid_only"'; then
    echo "✅ Champ 'is_paid_only' présent dans les réponses"
else
    echo "❌ Champ 'is_paid_only' manquant"
fi

if echo "$POSTS_USER6" | grep -q '"has_access"'; then
    echo "✅ Champ 'has_access' présent dans les réponses"
else
    echo "❌ Champ 'has_access' manquant"
fi

# Vérifier la logique d'accès
if echo "$POST38_USER6" | grep -q '"has_access":true'; then
    echo "✅ User6 a accès au post 38 (normal, post gratuit)"
elif echo "$POST38_USER6" | grep -q '"has_access":false'; then
    echo "⚠️  User6 n'a pas accès au post 38 (vérifie s'il est payant ou la logique d'abonnement)"
fi

echo ""
echo "🏁 Test terminé"
echo "==============="
