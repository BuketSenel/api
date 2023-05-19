package models

import "database/sql"

type Product struct {
	ID              int16          `json:"prod_id"`
	Name            string         `json:"prod_name"`
	Description     string         `json:"prod_desc"`
	CatID           int64          `json:"cat_id"`
	ResID           int64          `json:"rest_id"`
	Quantity        int64          `json:"prod_count"`
	Image           sql.NullString `json:"prod_image"`
	Price           float32        `json:"price"`
	Currency        string         `json:"currency"`
	PrepDurationMin int8           `json:"prep_dur_minute"`
	Availability    bool           `json:"availability"`
}
