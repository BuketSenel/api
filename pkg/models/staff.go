package models

import "time"

type Staff struct {
	ID        int64     `json:"id"`
	ResID     int64     `json:"RID"`
	Name      string    `json:"name"`
	Password  string    `json:"password"`
	Phone     string    `json:"phone"`
	Email     string    `json:"email"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
}
