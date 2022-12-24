package models

type Restaurant struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Summary   string `json:"summary"`
	Logo      string `json:"logo"`
	Address   string `json:"address"`
	District  string `json:"district"`
	City      string `json:"city"`
	Country   string `json:"country"`
	Phone     string `json:"phone"`
	Tags      string `json:"tags"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
