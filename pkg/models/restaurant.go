package models

import (
	"database/sql"
	"time"
)

type Restaurant struct {
	ID        int64          `json:"rest_id"`
	Name      string         `json:"rest_name"`
	Summary   string         `json:"summary"`
	Logo      sql.NullString `json:"logo"`
	Address   string         `json:"address"`
	District  string         `json:"district"`
	City      string         `json:"city"`
	Country   string         `json:"country"`
	Phone     string         `json:"rest_phone"`
	CreatedAt time.Time      `json:"rest_created_at"`
	UpdatedAt time.Time      `json:"rest_updated_at"`
}
