package models

import "time"

type User struct {
	ID        int64     `json:"user_id"`
	Name      string    `json:"user_name"`
	Password  string    `json:"password"`
	Phone     string    `json:"phone"`
	Email     string    `json:"email"`
	ResID     int64     `json:"rest_id"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"user_created_at"`
	UpdatedAt time.Time `json:"user_updated_at"`
}
