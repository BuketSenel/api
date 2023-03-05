package models

import "time"

type User struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Password  string    `json:"password"`
	Phone     string    `json:"phone"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
