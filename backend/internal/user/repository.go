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
	// Ajout de la mise à jour du prix si fourni
	if input.MonthlyPrice > 0 {
		updates["monthly_price"] = input.MonthlyPrice
	}
	// Ajout de la mise à jour du price_id Stripe si fourni
	if input.StripePriceID != "" {
		updates["stripe_price_id"] = input.StripePriceID
	}

	result = db.GormDB.Model(&user).Updates(updates)
	return result.Error
}

// UserRepository interface
type UserRepository interface {
	GetByID(id uint) (*User, error)
}

// repository implémentation
type repository struct{}

func NewRepository() UserRepository {
	return &repository{}
}

func (r *repository) GetByID(id uint) (*User, error) {
	return GetUserByID(id)
}
