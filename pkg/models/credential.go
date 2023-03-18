package models

type Credential struct {
	UserID   int64  `json:"user_id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	UserType string `json:"user_type"`
	Token    string `json:"token"`
}
