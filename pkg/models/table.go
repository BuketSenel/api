package models

type Table struct {
	ID    int64  `json:"id"`
	ResID int64  `json:"RID"`
	QR    string `json:"qr"`
}
