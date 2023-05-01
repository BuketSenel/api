package models

import "time"

type CustomQuery struct {
	UserID             int64     `json:"user_id"`
	UserName           string    `json:"user_name"`
	Password           string    `json:"password"`
	UserPhone          string    `json:"user_phone"`
	Email              string    `json:"email"`
	RestID             int64     `json:"rest_id"`
	RestName           string    `json:"rest_name"`
	RestPhone          string    `json:"rest_phone"`
	Summary            string    `json:"summary"`
	Logo               string    `json:"logo"`
	Address            string    `json:"address"`
	District           string    `json:"district"`
	City               string    `json:"city"`
	Country            string    `json:"country"`
	Tags               string    `json:"tags"`
	RestCreatedAt      time.Time `json:"rest_created_at"`
	RestUpdatedAt      time.Time `json:"rest_updated_at"`
	UserType           string    `json:"type"`
	UserCreatedAt      time.Time `json:"user_created_at"`
	UserUpdatedAt      time.Time `json:"user_updated_at"`
	ProductID          int16     `json:"prod_id"`
	ProductName        string    `json:"prod_name"`
	ProductDescription string    `json:"prod_description"`
	ProductImage       string    `json:"prod_image"`
	CatID              int64     `json:"cat_id"`
	CatName            string    `json:"cat_name"`
	ParentCatID        int16     `json:"parent_cat_id"`
	Quantity           int64     `json:"count"`
	CatImage           string    `json:"cat_image"`
	Price              float32   `json:"price"`
	Currency           string    `json:"currency"`
	PrepDurationMin    int8      `json:"prep_dur_minute"`
	OrderID            string    `json:"order_id"`
	OrderItemID        int64     `json:"order_item_id"`
	TableID            int64     `json:"table_id"`
	QR                 string    `json:"qr"`
	Details            string    `json:"details"`
	Status             string    `json:"order_status"`
	OrderCreatedAt     time.Time `json:"order_created_at"`
	OrderUpdatedAt     time.Time `json:"order_updated_at"`
	OrderItemTotalQty  int64     `json:"total_quantity"`
}
