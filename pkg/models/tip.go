package models

import "time"

type Tip struct {
	UserID    int64     `json:"user_id"`
	WaiterID  int64     `json:"waiter_id"`
	Tip       float32   `json:"tip"`
	RestID    int64     `json:"rest_id"`
	CreatedAt time.Time `json:"tip_created_at"`
}
