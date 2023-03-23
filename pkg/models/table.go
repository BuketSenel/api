package models

type Table struct {
	ID    int64  `json:"table_id"`
	ResID int64  `json:"rest_id"`
	QR    string `json:"qr"`
}
