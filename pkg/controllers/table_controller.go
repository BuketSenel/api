package controllers

import (
	"database/sql"
	"fmt"

	"github.com/SelfServiceCo/api/pkg/models"
)

func GetRestaurantTables(resID int64, tableId int64) []models.Table {

	db, err := sql.Open("mysql", conf.Name+":"+conf.Password+"@tcp("+conf.Db+":3306)/selfservice")

	if err != nil {
		fmt.Println("Err", err.Error())
		return nil
	}

	if tableId != 0 {
		results, err := db.Query("SELECT * FROM tables WHERE RID = ? AND id = ?", resID, tableId)
		defer db.Close()

		if err != nil {
			fmt.Println("Err", err.Error())
			return nil
		}
		tables := []models.Table{}
		for results.Next() {
			var table models.Table
			err = results.Scan(&table.ID, &table.ResID, &table.QR)
			if err != nil {
				panic(err.Error())
			}
			tables = append(tables, table)
		}

		return tables
	}
	results, err := db.Query("SELECT * FROM tables WHERE RID = ?", resID)
	defer db.Close()

	if err != nil {
		fmt.Println("Err", err.Error())
		return nil
	}
	tables := []models.Table{}
	for results.Next() {
		var table models.Table
		err = results.Scan(&table.ID, &table.ResID, &table.QR)
		if err != nil {
			panic(err.Error())
		}
		tables = append(tables, table)
	}

	return tables
}
