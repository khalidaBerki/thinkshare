package user

import (
	"backend/internal/db"
	"errors"
)

var ErrUserNotFound = errors.New("utilisateur non trouvé")

func GetUserByID(id uint) (*User, error) {
	var user User
	result := db.GormDB.First(&user, id)
	if result.Error != nil {
		return nil, ErrUserNotFound
	}
	return &user, nil
}

func UpdateUserProfile(id uint, input UpdateUserInput) error {
	// Vérifie d'abord si l'utilisateur existe
	var user User
	result := db.GormDB.First(&user, id)
	if result.Error != nil {
		return ErrUserNotFound
	}

	// Met à jour uniquement les champs autorisés
	updates := map[string]interface{}{
		"full_name":  input.FullName,
		"bio":        input.Bio,
		"avatar_url": input.AvatarURL,
	}

	result = db.GormDB.Model(&user).Updates(updates)
	return result.Error
}
