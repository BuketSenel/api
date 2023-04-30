package models

import "github.com/skip2/go-qrcode"

type Table struct {
	TableNo    int64          `json:"table_no"`
	RestID     int64          `json:"rest_id"`
	WaiterID   int64          `json:"waiter_id"`
	QR         *qrcode.QRCode `json:"qr_code"`
	NewTableNo int64          `json:"new_table_no"`
	WaiterName string         `json:"user_name"`
	QRString   string         `json:"qr"`
}
