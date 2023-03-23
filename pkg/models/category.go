package models

type Category struct {
	ID          int64  `json:"cat_id"`
	RID         int64  `json:"rest_id"`
	Name        string `json:"cat_name"`
	Description string `json:"cat_description"`
	Image       string `json:"cat_image"`
	ParentCatID int16  `json:"parent_cat_id"`
}
