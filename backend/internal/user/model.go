package user

type User struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	FullName  string `json:"full_name"`
	Bio       string `json:"bio"`
	AvatarURL string `json:"avatar_url"`
	Email     string `json:"email"`
}

type UpdateUserInput struct {
	FullName  string `json:"full_name"`
	Bio       string `json:"bio"`
	AvatarURL string `json:"avatar_url"`
}
