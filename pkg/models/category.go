package models

type Category struct {
	ID          int64  `json:"id"`
	RID         int64  `json:"RID"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Image       string `json:"image"`
	ParentCatID int16  `json:"parentCategoryId"`
}
