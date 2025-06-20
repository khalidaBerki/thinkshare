# Guide sur les contraintes de clé étrangère en PostgreSQL

## Problème : Violation de clé étrangère

Lorsque vous rencontrez cette erreur :

```
ERROR: update or delete on table "posts" violates foreign key constraint "fk_posts_media" on table "media" (SQLSTATE 23503)
```

Cela signifie que vous essayez de supprimer ou de mettre à jour une ligne dans la table `posts` qui est référencée par une ou plusieurs lignes dans la table `media` via la contrainte `fk_posts_media`.

## Explication des clés étrangères

Une clé étrangère est une contrainte d'intégrité qui garantit que les valeurs dans une colonne (ou un ensemble de colonnes) correspondent aux valeurs apparaissant dans une autre table. Elle empêche les actions qui détruiraient ces liens.

Dans votre cas :
- La table `media` a une colonne `PostID` qui référence la colonne `ID` de la table `posts`
- La contrainte `fk_posts_media` empêche la suppression d'un post tant que des entrées média y font référence

## Solutions possibles

### 1. Supprimer d'abord les médias (Solution implémentée)

La solution que nous avons implémentée consiste à supprimer d'abord les entrées de la table `media` avant de supprimer le post :

```go
// 1. D'abord supprimer les entrées de médias dans la base de données
if err := tx.Delete(&post.Media).Error; err != nil {
    tx.Rollback()
    return err
}

// 2. Ensuite supprimer le post
if err := tx.Delete(&post).Error; err != nil {
    tx.Rollback()
    return err
}
```

### 2. Configurer les suppressions en cascade

Une alternative consiste à modifier la définition de votre contrainte de clé étrangère pour ajouter `ON DELETE CASCADE`. Cela signifie que lorsqu'un post est supprimé, toutes les entrées média associées sont automatiquement supprimées.

Pour modifier une contrainte existante en PostgreSQL :

```sql
-- Supprimer la contrainte existante
ALTER TABLE media DROP CONSTRAINT fk_posts_media;

-- Recréer la contrainte avec CASCADE
ALTER TABLE media ADD CONSTRAINT fk_posts_media 
    FOREIGN KEY (post_id) REFERENCES posts(id) 
    ON DELETE CASCADE;
```

Si vous utilisez GORM pour gérer vos migrations, vous pouvez spécifier ce comportement dans votre modèle :

```go
type Post struct {
    ID     uint
    // ...
    Media  []Media `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE;"`
}
```

## Avantages et inconvénients des deux approches

### Suppression explicite (actuelle)

**Avantages :**
- Contrôle total du processus de suppression
- Possibilité de gérer les erreurs à chaque étape

**Inconvénients :**
- Nécessite plus de code à maintenir
- Plus facile d'introduire des bugs

### Suppressions en cascade

**Avantages :**
- Code plus simple et plus propre
- Gestion automatique par la base de données

**Inconvénients :**
- Moins de contrôle sur le processus
- Potentiellement plus difficile à déboguer

## Recommandation

Si votre application n'a pas besoin d'une logique spécifique entre la suppression d'un post et de ses médias, la solution CASCADE est généralement préférable. Elle est plus robuste et plus simple à maintenir.

Cependant, dans votre cas, puisque vous devez également supprimer les fichiers physiques, la suppression explicite (notre solution actuelle) est probablement plus appropriée.
