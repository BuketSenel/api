package models

import "time"

type Restaurant struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Password  string    `json:"password"`
	Summary   string    `json:"summary"`
	Logo      string    `json:"logo"`
	Address   string    `json:"address"`
	District  string    `json:"district"`
	City      string    `json:"city"`
	Country   string    `json:"country"`
	Phone     string    `json:"phone"`
	Email     string    `json:"email"`
	Tags      string    `json:"tags"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
