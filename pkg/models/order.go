package models

import "time"

type Order struct {
	ID      int64     `json:"ID"`
	UserID  int64     `json:"user_id"`
	ResID   int64     `json:"RID"`
	TableID int64     `json:"table_id"`
	Details string    `json:"details"`
	Status  string    `json:"status"`
	Created time.Time `json:"created"`
}
