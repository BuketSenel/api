package models

import "time"

type Restaurant struct {
	ID        int64     `json:"rest_id"`
	Name      string    `json:"rest_name"`
	Summary   string    `json:"summary"`
	Logo      string    `json:"logo"`
	Address   string    `json:"address"`
	District  string    `json:"district"`
	City      string    `json:"city"`
	Country   string    `json:"country"`
	Phone     string    `json:"rest_phone"`
	Email     string    `json:"email"`
	Tags      string    `json:"tags"`
	CreatedAt time.Time `json:"rest_created_at"`
	UpdatedAt time.Time `json:"rest_updated_at"`
}
