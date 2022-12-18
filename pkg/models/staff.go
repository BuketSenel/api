package models

import "time"

type Staff struct {
	ID       int64
	ResID    int64
	Name     string
	Password string
	Phone    string
	Email    string
	Type     string
	Created  time.Time
}
