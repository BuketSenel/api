package models

type Category struct {
	ID          int16  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ParentCatID int16  `json:"parentCategoryId"`
	Image       string `json:"image"`
}
