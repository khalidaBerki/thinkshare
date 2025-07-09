package user

func GetProfile(userID uint) (*User, error) {
	return GetUserByID(userID)
}

func UpdateProfile(userID uint, input UpdateUserInput) error {
	return UpdateUserProfile(userID, input)
}

func SetUserMonthlyPrice(userID uint, price float64) error {
	return SetUserMonthlyPrice(userID, price)
}
