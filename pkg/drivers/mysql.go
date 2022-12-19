package drivers

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func mysqlConnect() {
	sql.Open()
}
