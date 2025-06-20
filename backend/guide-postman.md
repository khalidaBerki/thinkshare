# Guide d'utilisation de Postman pour l'API

## Récupérer des Posts

### 1. Configuration générale

- **URL de base** : `http://localhost:8080`
- **Headers requis pour toute requête authentifiée** :
  - `Authorization: Bearer VOTRE_TOKEN_JWT`

### 2. Obtenir un Token JWT

1. Ouvrez Postman
2. Créez une requête POST vers : `http://localhost:8080/api/fake-login`
3. Cliquez sur "Send"
4. Copiez le token JWT de la réponse (vous en aurez besoin pour les autres requêtes)

![Exemple de requête pour obtenir un token](https://i.imgur.com/example1.png)

### 3. Récupérer tous les posts

1. Créez une requête GET vers : `http://localhost:8080/api/posts`
2. Dans l'onglet "Headers", ajoutez :
   - Key: `Authorization`
   - Value: `Bearer VOTRE_TOKEN_JWT` (remplacez par votre token)
3. Cliquez sur "Send"

**Paramètres optionnels :**
- `page` : Numéro de page (défaut: 1)
- `pageSize` : Nombre de posts par page (défaut: 10)
- `visibility` : Filtrer par visibilité (`public` ou `private`)

Exemple avec paramètres : `http://localhost:8080/api/posts?page=1&pageSize=5&visibility=public`

### 4. Récupérer un post spécifique

Pour récupérer un post par ID, il faut d'abord implémenter cette fonctionnalité car elle ne semble pas disponible actuellement.

Cependant, vous pouvez récupérer tous les posts et filtrer côté client.

### 5. Tester avec Postman

#### Configuration de l'environnement

1. Créez un environnement dans Postman (en haut à droite)
2. Ajoutez une variable `base_url` avec la valeur `http://localhost:8080`
3. Ajoutez une variable `token` (elle sera remplie automatiquement par le script ci-dessous)

#### Collection de requêtes recommandées

1. **Obtenir un token** (POST)
   - URL : `{{base_url}}/api/fake-login`
   - Dans l'onglet "Tests", ajoutez ce script pour stocker automatiquement le token :
   ```javascript
   if (pm.response.code === 200) {
       var jsonData = pm.response.json();
       pm.environment.set("token", jsonData.token);
       console.log("Token JWT enregistré dans la variable d'environnement");
   }
   ```

2. **Récupérer tous les posts** (GET)
   - URL : `{{base_url}}/api/posts`
   - Headers : `Authorization: Bearer {{token}}`

3. **Créer un post texte** (POST)
   - URL : `{{base_url}}/api/posts`
   - Headers : `Authorization: Bearer {{token}}`
   - Body : form-data
     - `content`: Contenu du post
     - `visibility`: public ou private

4. **Mettre à jour un post** (PUT)
   - URL : `{{base_url}}/api/posts/1` (remplacez 1 par l'ID du post)
   - Headers : `Authorization: Bearer {{token}}`
   - Body : raw (JSON)
   ```json
   {
       "content": "Contenu modifié",
       "visibility": "public"
   }
   ```

5. **Supprimer un post** (DELETE)
   - URL : `{{base_url}}/api/posts/1` (remplacez 1 par l'ID du post)
   - Headers : `Authorization: Bearer {{token}}`

## Exemples de requêtes cURL

### Obtenir un token
```bash
curl -X POST http://localhost:8080/api/fake-login
```

### Récupérer les posts
```bash
curl -X GET \
  -H "Authorization: Bearer VOTRE_TOKEN_JWT" \
  http://localhost:8080/api/posts
```

### Récupérer avec pagination
```bash
curl -X GET \
  -H "Authorization: Bearer VOTRE_TOKEN_JWT" \
  "http://localhost:8080/api/posts?page=1&pageSize=5"
```

## Dépannage des erreurs courantes

### Erreur lors de la suppression d'un post

Si vous rencontrez une erreur comme celle-ci lors de la suppression d'un post :

```json
{
    "error": "failed to delete media file: remove uploads\\user_1_25bf4e5f-d797-4e18-8822-0ac2d037a071.png: The system cannot find the file specified."
}
```

Cela signifie que le serveur essaie de supprimer un fichier média associé au post, mais que ce fichier n'existe plus sur le disque.

**Solutions :**

1. Avec la dernière version du code, cette erreur a été corrigée et le serveur ignorera désormais les fichiers manquants lors de la suppression d'un post.

2. Si vous rencontrez encore cette erreur, vérifiez que :
   - Le dossier `uploads` existe et est accessible
   - Les permissions sur le dossier `uploads` sont correctes
   - Aucun processus ne nettoie ou déplace les fichiers uploadés
