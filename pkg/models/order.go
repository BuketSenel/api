package models

import "time"

type Order struct {
	ID      int64     `json:"ID"`
	UserID  int64     `json:"user_id"`
	ResID   int64     `json:"RID"`
	Content string    `json:"content"`
	Status  string    `json:"status"`
	Created time.Time `json:"created"`
}
