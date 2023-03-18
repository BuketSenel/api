package models

import "time"

type Order struct {
	ID      string    `json:"ID"`
	UserID  int64     `json:"userId"`
	ResID   int64     `json:"restId"`
	TableID int64     `json:"tableId"`
	Details string    `json:"details"`
	Status  string    `json:"orderStatus"`
	Created time.Time `json:"created"`
}
