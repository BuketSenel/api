package drivers

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type DB struct {
	SQL *sql.DB
}

var dbConnection = &DB{}

const maxOpen = 10
const maxIdle = 5
const ttl = 5 * time.Minute

func MysqlConnect(instance string) (*DB, error) {
	d, err := newConnection(instance)
	if err != nil {
		panic(err)
	}

	d.SetMaxOpenConns(maxOpen)
	d.SetMaxIdleConns(maxIdle)
	d.SetConnMaxLifetime(maxOpen)

	dbConnection.SQL = d

	err = testDB(d)
	if err != nil {
		return nil, err
	}
	return dbConnection, nil
}

func newConnection(instance string) (*sql.DB, error) {
	db, err := sql.Open("mysql", instance)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func testDB(d *sql.DB) error {
	err := d.Ping()
	if err != nil {
		return err
	}
	return nil
}
