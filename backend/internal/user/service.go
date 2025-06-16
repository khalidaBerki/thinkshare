package user

func GetProfile(userID uint) (*User, error) {
	return GetUserByID(userID)
}

func UpdateProfile(userID uint, input UpdateUserInput) error {
	return UpdateUserProfile(userID, input)
}
