package models

type Table struct {
	TableNo    int64  `json:"table_no"`
	RestID     int64  `json:"rest_id"`
	WaiterID   int64  `json:"waiter_id"`
	QR         []byte `json:"qr"`
	NewTableNo int64  `json:"new_table_no"`
	WaiterName string `json:"waiter_name"`
}
