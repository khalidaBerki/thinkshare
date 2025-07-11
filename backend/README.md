# ThinkShare API — Backend

API backend pour ThinkShare, un réseau social collaboratif avec gestion des utilisateurs, posts, commentaires, likes, messagerie privée, abonnements et paiements Stripe.

---

## Démarrage rapide

### 1. **Prérequis**

- Go 1.21+
- PostgreSQL (Azure ou local)
- Stripe (pour les paiements)
- [swaggo/swag](https://github.com/swaggo/swag) pour la doc Swagger

### 2. **Variables d’environnement**

Crée un fichier `.env` ou configure dans ton shell :

```sh
PORT=
GIN_MODE=
JWT_SECRET=
PGHOST=
PGUSER=
PGPORT=
PGDATABASE=
PGPASSWORD=
PGSSLMODE=
STRIPE_SECRET_KEY=
STRIPE_WEBHOOK_SECRET=
```

### 3. **Installation des dépendances**

```sh
go mod tidy
```

### 4. **Générer la documentation Swagger**

```sh
swag init
```

### 5. **Lancer le serveur**

```sh
go run main.go
```

---

## Documentation API

Swagger est disponible sur :  
`http://localhost:8080/swagger/index.html`

---

## Structure des dossiers

```
backend/
│
├── internal/
│   ├── auth/         # Authentification, JWT, OAuth
│   ├── user/         # Utilisateurs
│   ├── post/         # Posts et médias
│   ├── comment/      # Commentaires
│   ├── like/         # Likes
│   ├── message/      # Messagerie privée
│   ├── subscription/ # Abonnements/followers
│   ├── media/        # Gestion des fichiers médias
│   ├── payment/      # Paiements Stripe
│   └── db/           # Connexion DB
│
├── uploads/          # Fichiers uploadés (images, docs, vidéos)
├── docs/             # Documentation Swagger auto-générée
├── main.go           # Point d’entrée du serveur
└── go.mod
```

---

## Authentification

- JWT pour toutes les routes protégées (`Authorization: Bearer <token>`)
- OAuth Google disponible

---

## Principales routes

### Utilisateur

- `POST /register` — Inscription
- `POST /login` — Connexion (retourne token + user_id)
- `GET /api/profile` — Profil utilisateur connecté
- `PUT /api/profile` — Modifier son profil
- `GET /api/users/{id}` — Profil public d’un utilisateur

### Posts

- `POST /api/posts` — Créer un post (texte + médias)
- `GET /api/posts` — Tous les posts (scroll infini)
- `GET /api/posts/user/{id}` — Posts d’un utilisateur
- `GET /api/posts/{id}` — Détail d’un post
- `PUT /api/posts/{id}` — Modifier un post
- `DELETE /api/posts/{id}` — Supprimer un post

### Commentaires

- `POST /api/comments` — Ajouter un commentaire
- `GET /api/comments/{postID}` — Commentaires d’un post
- `PUT /api/comments/{id}` — Modifier un commentaire
- `DELETE /api/comments/{id}` — Supprimer un commentaire

### Likes

- `POST /api/likes/posts/{postID}` — Like/unlike un post
- `GET /api/likes/posts/{postID}` — Stats de likes d’un post

### Messagerie

- `POST /api/messages` — Envoyer un message privé
- `GET /api/messages/conversations` — Liste des conversations
- `GET /api/messages/{otherUserID}` — Conversation avec un utilisateur
- `PATCH /api/messages/{senderID}/read` — Marquer comme lu
- `PUT /api/messages/{id}` — Modifier un message
- `DELETE /api/messages/{id}` — Supprimer un message

### Abonnements

- `POST /api/subscribe` — S’abonner à un utilisateur
- `POST /api/unsubscribe` — Se désabonner
- `GET /api/followers/{id}` — Voir les abonnés
- `GET /api/subscriptions` — Voir ses abonnements

### Médias

- `GET /api/media/{id}` — Récupérer un média
- `DELETE /api/media/{id}` — Supprimer un média
- `GET /api/media/post/{postID}` — Médias d’un post
- `PUT /api/media/{id}/metadata` — Modifier les métadonnées
- `POST /api/media/cleanup` — Nettoyer les médias orphelins

### Paiement Stripe

- `POST /api/payment/webhook` — Webhook Stripe (public)

---

## Notes

- Les fichiers uploadés doivent être accessibles via `/uploads/...`
- Les endpoints Stripe doivent être configurés avec les secrets corrects
- Les permissions sont gérées par middleware JWT

---

## Développement

- Pour activer les routes de debug, lance en mode `debug` (`GIN_MODE=debug`)
- Pour migrer la base, vérifie la logique dans `main.go` et `internal/db/`

---

## Swagger

- Les handlers sont annotés pour Swagger.
- Regénère la doc avec :  
  ```sh
  swag init
  ```

---

## Contact

Pour toute question, bug ou suggestion, contacte l’équipe ThinkShare.

---