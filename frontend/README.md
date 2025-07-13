# ThinkShare - Application de partage de connaissances

![ThinkShare Logo](assets/images/logo.png)

## Table des matières

- [Introduction](#introduction)
- [Fonctionnalités](#fonctionnalités)
- [Technologies utilisées](#technologies-utilisées)
- [Installation](#installation)
- [Structure du projet](#structure-du-projet)
- [Fonctionnalités détaillées](#fonctionnalités-détaillées)
- [Architecture](#architecture)
- [Contribuer](#contribuer)
- [Licence](#licence)

## Introduction

ThinkShare est une plateforme collaborative permettant aux utilisateurs de partager des idées, des documents et des connaissances dans un environnement professionnel. L'application permet de créer des posts, d'interagir avec d'autres utilisateurs et de suivre les tendances et l'activité de la plateforme.

## Fonctionnalités

### Authentification et gestion des utilisateurs
- Inscription et connexion sécurisée
- Authentification par e-mail/mot de passe
- Récupération de mot de passe
- Profils utilisateurs personnalisables

### Gestion de contenu
- Création de posts avec texte formaté
- Support pour différents types de documents (rapports, mémos, notes)
- Upload et partage de médias (images, vidéos, documents)
- Commentaires et discussions sur les posts

### Recherche et exploration
- Recherche avancée de contenu
- Filtres par type de document, date, et popularité
- Suggestions de contenu personnalisées
- Exploration des tendances

### Messagerie
- Conversations privées entre utilisateurs
- Partage de fichiers dans les conversations
- Notifications en temps réel

### Dashboard Admin
- Statistiques d'utilisation globales
- Graphique de posts par jour
- Graphique des médias (images/vidéos/documents)
- Affichage des posts populaires et impopulaires
- Monitoring de l'activité des utilisateurs

## Technologies utilisées

- **Frontend**: Flutter (Dart)
- **État**: Provider
- **UI/UX**: Material Design
- **Graphiques**: fl_chart
- **Formatage de dates**: intl
- **Navigation**: go_router

## Installation

### Prérequis
- Flutter SDK (version 3.10.0 ou supérieure)
- Dart SDK (version 3.0.0 ou supérieure)
- Un éditeur de code (VS Code recommandé)
- Un appareil Android/iOS ou émulateur

### Étapes d'installation

1. Cloner le répertoire
   ```bash
   git clone https://github.com/khalidaBerki/thinkshare.git
   cd thinkshare/frontend
   ```

2. Installer les dépendances
   ```bash
   flutter pub get
   ```

3. Lancer l'application
   ```bash
   flutter run -d chrome  # Pour le web
   flutter run            # Pour mobile
   ```

## Structure du projet

```
lib/
├── core/
│   ├── config/             # Configuration de l'application
│   ├── constants/          # Constantes (couleurs, textes, etc.)
│   ├── errors/             # Gestion des erreurs
│   ├── routing/            # Configuration des routes
│   └── utils/              # Utilitaires généraux
├── features/
│   ├── admin_dashboard/    # Dashboard administrateur
│   ├── auth/               # Authentification
│   ├── chat/               # Messagerie
│   ├── home/               # Écran d'accueil et posts
│   ├── profile/            # Profil utilisateur
│   └── search/             # Recherche
└── main.dart               # Point d'entrée de l'application
```

## Fonctionnalités détaillées

### Dashboard Admin
- **KPIs**: Total des posts, commentaires et conversations
- **Graphique temporel**: Visualisation des posts par jour
- **Médias**: Graphique des médias partagés (images/vidéos/documents)
- **Top/Flop**: Liste des posts les plus/moins populaires
- **Aperçu du contenu**: Prévisualisation des posts récents

### Création de Posts
- **Types de documents**: Support pour différents formats professionnels
- **Média**: Upload d'images, vidéos et documents
- **Visibilité**: Options de confidentialité (public/privé)
- **Formatage**: Mise en forme du texte

### Profil Utilisateur
- **Statistiques personnelles**: Activité, posts populaires
- **Personnalisation**: Photo de profil, bio, informations professionnelles
- **Historique**: Visualisation des actions passées

### Recherche
- **Filtres avancés**: Par date, type, popularité
- **Suggestions**: Contenu recommandé basé sur les intérêts

## Architecture

L'application suit une architecture **Feature-First** avec séparation claire des responsabilités:

- **Presentation**: Widgets, écrans, UI
- **Domain**: Logique métier, entités, use cases
- **Data**: Sources de données, repositories, modèles

La gestion d'état est assurée par **Provider**, permettant une mise à jour réactive de l'interface utilisateur.

## Contribuer

1. Forkez le projet
2. Créez une branche pour votre fonctionnalité (`git checkout -b feature/amazing-feature`)
3. Committez vos changements (`git commit -m 'Add some amazing feature'`)
4. Poussez vers la branche (`git push origin feature/amazing-feature`)
5. Ouvrez une Pull Request

---

© 2024 ThinkShare. Tous droits réservés.