package controllers

import (
	"database/sql"

	"github.com/SelfServiceCo/api/pkg/drivers"
)

var conf = drivers.MysqlConfigLoad()

func CreateConnection() *sql.DB {
	db, err := sql.Open("mysql", conf.Name+":"+conf.Password+"@tcp("+conf.Db+":3306)/selfservice?parseTime=true")
	if err != nil {
		return nil
	}
	return db
}

func CloseConnection(db *sql.DB) {
	db.Close()
}
