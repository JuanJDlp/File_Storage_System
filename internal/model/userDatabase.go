package model

type UserDatabase struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}
