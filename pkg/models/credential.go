package models

type Credential struct {
	UserID   int64  `json:"user_id"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	UserType string `json:"user_type"`
	Token    string `json:"token"`
}
