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
	ProductId   int64     `json:"prod_id"`
	ProductName string    `json:"prod_name"`
	Price       float32   `json:"price"`
	Quantity    string    `json:"prod_count"`
	CreatedAt   time.Time `json:"order_created_at"`
	UpdatedAt   time.Time `json:"order_updated_at"`
}
