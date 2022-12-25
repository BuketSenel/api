package models

type Product struct {
	ID              int16   `json:"id"`
	Name            string  `json:"name"`
	Description     string  `json:"description"`
	CatID           int64   `json:"categoryId"`
	ResID           int64   `json:"restorantId"`
	Image           string  `json:"image"`
	Price           float32 `json:"price"`
	Currency        string  `json:"currency"`
	PrepDurationMin int8    `json:"prepDurationMinute"`
}
