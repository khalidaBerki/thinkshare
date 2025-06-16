thinkshare/backend/
├── cmd/                   # Point d'entrée de l'application (main.go)
│   └── server/            # Serveur HTTP/API
│       └── main.go
│
├── internal/              # Code métier (privé à l'app)
│   ├── user/              # Logique utilisateur
│   │   ├── handler.go     # HTTP handlers
│   │   ├── service.go     # Règles métier
│   │   ├── repository.go  # Requêtes DB
│   │   └── model.go       # Structures données
│   ├── post/
│   ├── subscription/
│   ├── payment/
│   ├── messaging/
│   └── auth/
│
├── pkg/                   # Code réutilisable (auth, utils, jwt, etc.)
│   ├── config/            # Chargement config/env
│   ├── database/          # Connexion DB
│   ├── logger/            # Logs formatés
│   └── middleware/        # Middleware (auth, logging...)
│
├── api/                   # Fichiers Swagger/OpenAPI
│   └── docs.go
│
├── migrations/            # Migrations SQL
│   └── 001_init.sql
│
├── test/                  # Tests unitaires
│   └── user_test.go
│
├── go.mod
├── go.sum
└── README.md

PECeemi2025§
PECeemi2025§