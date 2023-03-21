package models

import "time"

type Order struct {
	ID          string    `json:"order_id"`
	OrderItemID string    `json:"order_item_id"`
	UserID      int64     `json:"user_id"`
	ResID       int64     `json:"rest_id"`
	TableID     int64     `json:"table_id"`
	Details     string    `json:"details"`
	Status      string    `json:"order_status"`
	Created     time.Time `json:"CREATED_AT"`
	Updated     time.Time `json:"UPDATED_AT"`
}
