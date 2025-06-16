//Mock temporairement le repository

package user

import "fmt"

var mockUser = User{
	ID:        1,
	Username:  "haithem95",
	FullName:  "Haithem Hammami",
	Bio:       "Développeur fullstack passionné.",
	AvatarURL: "https://cdn.thinkshare/avatar.jpg",
	Email:     "haithem@eemi.com",
}

func GetUserByID(id uint) (*User, error) {
	if id == mockUser.ID {
		return &mockUser, nil
	}
	return nil, ErrUserNotFound
}

func UpdateUserProfile(id uint, input UpdateUserInput) error {
	if id != mockUser.ID {
		return ErrUserNotFound
	}

	mockUser.FullName = input.FullName
	mockUser.Bio = input.Bio
	mockUser.AvatarURL = input.AvatarURL

	return nil
}

var ErrUserNotFound = fmt.Errorf("user not found")
