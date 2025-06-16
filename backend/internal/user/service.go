package user

func GetProfile(userID int) (*User, error) {
	return GetUserByID(userID)
}

func UpdateProfile(userID int, input UpdateUserInput) error {
	return UpdateUserProfile(userID, input)
}
