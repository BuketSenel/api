package models

import "time"

type Order struct {
	ID          string    `json:"order_id"`
	OrderItemID int64     `json:"order_item_id"`
	UserID      int64     `json:"user_id"`
	ResID       int64     `json:"rest_id"`
	TableID     int64     `json:"table_id"`
	Details     string    `json:"details"`
	Status      string    `json:"order_status"`
	PName       string    `json:"pname"`
	PDesc       string    `json:"pdesc"`
	Quantity    string    `json:"quantity"`
	Created     time.Time `json:"CREATED_AT"`
	Updated     time.Time `json:"UPDATED_AT"`
}
