package users

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserLogin struct {
	accessToken string
	ID          int    `json:"id"`
	Email       string `json:"email"`
	Username    string `json:"username"`
}

type UserResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}
