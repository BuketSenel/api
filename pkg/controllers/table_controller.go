package controllers

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/SelfServiceCo/api/pkg/models"
	"github.com/gin-gonic/gin"
)

func GetRestaurantTables(resID int64, tableId int64) []models.Table {

	db, err := sql.Open("mysql", conf.Name+":"+conf.Password+"@tcp("+conf.Db+":3306)/selfservice")

	if err != nil {
		fmt.Println("Err", err.Error())
		return nil
	}

	if tableId != 0 {
		results, err := db.Query("SELECT * FROM tables WHERE rest_id = ? AND table_id = ?", resID, tableId)
		defer db.Close()

		if err != nil {
			fmt.Println("Err", err.Error())
			return nil
		}
		tables := []models.Table{}
		for results.Next() {
			var table models.Table
			err = results.Scan(&table.TableNo, &table.RestID, &table.QR)
			if err != nil {
				panic(err.Error())
			}
			tables = append(tables, table)
		}

		return tables
	}
	results, err := db.Query("SELECT * FROM tables WHERE rest_id = ?", resID)
	defer db.Close()

	if err != nil {
		fmt.Println("Err", err.Error())
		return nil
	}
	tables := []models.Table{}
	for results.Next() {
		var table models.Table
		err = results.Scan(&table.TableNo, &table.RestID, &table.QR)
		if err != nil {
			panic(err.Error())
		}
		tables = append(tables, table)
	}

	return tables
}

func GetTable(tableId int64) (models.Table, gin.H) {
	var table models.Table
	db, err := sql.Open("mysql", conf.Name+":"+conf.Password+"@tcp("+conf.Db+":3306)/selfservice")

	if err != nil {
		return table, gin.H{"error": err.Error()}
	}

	results, err := db.Query("SELECT * FROM tables WHERE table_id = ?", tableId)
	defer db.Close()

	if err != nil {
		fmt.Println("Err", err.Error())
		return table, gin.H{"error": err.Error()}
	}

	for results.Next() {
		err = results.Scan(&table.TableNo, &table.RestID, &table.QR)
		if err != nil {
			return table, gin.H{"error": err.Error()}
		}
	}

	return table, gin.H{"status": http.StatusOK, "message": "success", "data": results}
}

func OrdersByTable(tableId int64) ([]models.CustomQuery, gin.H) {

	db, err := sql.Open("mysql", conf.Name+":"+conf.Password+"@tcp("+conf.Db+":3306)/selfservice")

	if err != nil {
		return nil, gin.H{"status": http.StatusBadRequest, "message": "DB Connection Error! Orders By Table"}
	}

	results, err := db.Query("SELECT * FROM orders WHERE table_id = ?", tableId)
	defer db.Close()

	if err != nil {
		return nil, gin.H{"status": http.StatusBadRequest, "message": "DB Query Error! Orders By Table"}
	}

	customQuery := []models.CustomQuery{}
	for results.Next() {
		var cq models.CustomQuery
		err = results.Scan(&cq.OrderItemID, &cq.TableID, &cq.ProductID, &cq.Status, &cq.Quantity)
		if err != nil {
			return nil, gin.H{"status": http.StatusBadRequest, "message": "DB Scan Error! Orders By Table"}
		}
		customQuery = append(customQuery, cq)
	}

	return customQuery, gin.H{"status": http.StatusOK, "message": "success", "data": results}
}

func AddTable(table models.Table) (bool, gin.H) {

	db, err := sql.Open("mysql", conf.Name+":"+conf.Password+"@tcp("+conf.Db+":3306)/selfservice")

	if err != nil {
		fmt.Println("Err", err.Error())
		return false, gin.H{"status": http.StatusBadRequest, "message": "DB Connection Error! Add Table"}
	}

	_, err = db.Exec("INSERT INTO tables (rest_id, qr) VALUES (?, ?)", table.RestID, table.QR)
	defer db.Close()

	if err != nil {
		fmt.Println("Err", err.Error())
		return false, gin.H{"status": http.StatusBadRequest, "message": "DB Query Error! Add Table"}
	}

	return true, gin.H{"status": http.StatusOK, "message": "success"}
}
