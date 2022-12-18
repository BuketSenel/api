package models

import "time"

type Order struct {
	ID      int64
	UserID  int64
	ResID   int64
	Content string
	Status  string
	Created time.Time
}
