package controllers

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/SelfServiceCo/api/pkg/models"
	"github.com/gin-gonic/gin"
	qrcode "github.com/skip2/go-qrcode"
)

func GetRestaurantTables(resID int64, tableId int64) ([]models.Table, gin.H) {

	db, err := sql.Open("mysql", conf.Name+":"+conf.Password+"@tcp("+conf.Db+":3306)/selfservice")

	if err != nil {
		return nil, gin.H{"status": http.StatusBadRequest, "message": "DB Connection Error! Get Restaurant Tables"}
	}

	if tableId != 0 {
		results, err := db.Query("SELECT * FROM tables WHERE rest_id = ? AND table_no = ?", resID, tableId)
		defer db.Close()

		if err != nil {
			return nil, gin.H{"status": http.StatusBadRequest, "message": "Query Error! Get Restaurant Tables"}
		}
		tables := []models.Table{}
		for results.Next() {
			var table models.Table
			err = results.Scan(&table.RestID, &table.TableNo, &table.WaiterID, &table.QR)
			if err != nil {
				return nil, gin.H{"status": http.StatusBadRequest, "message": "Scan Error! Get Restaurant Tables"}
			}
			tables = append(tables, table)
		}

		return tables, gin.H{"status": http.StatusOK, "message": "success", "data": tables}
	}
	results, err := db.Query("SELECT * FROM tables WHERE rest_id = ?", resID)
	defer db.Close()

	if err != nil {
		return nil, gin.H{"status": http.StatusBadRequest, "message": "Query Error! Get Restaurant Tables", "error": err.Error()}
	}
	tables := []models.Table{}
	for results.Next() {
		var table models.Table
		err = results.Scan(&table.RestID, &table.TableNo, &table.WaiterID, &table.QR)
		if err != nil {
			return nil, gin.H{"status": http.StatusBadRequest, "message": "Scan Error! Get Restaurant Tables"}
		}
		tables = append(tables, table)
	}

	return tables, gin.H{"status": http.StatusOK, "message": "success", "data": tables}
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

func AddTable(c *gin.Context) (bool, gin.H) {
	table := models.Table{}
	err := c.BindJSON(&table)

	if err != nil {
		return false, gin.H{"status": http.StatusBadRequest, "message": "Bind Error! Create Product"}
	}

	db, err := sql.Open("mysql", conf.Name+":"+conf.Password+"@tcp("+conf.Db+":3306)/selfservice")

	if err != nil {
		fmt.Println("Err", err.Error())
		return false, gin.H{"status": http.StatusBadRequest, "message": "DB Connection Error! Add Table"}
	}

	var header = gin.H{}
	table.QR, header = CreateQRCode(table.TableNo, table.RestID)
	if header["status"] != http.StatusOK {
		return false, header
	}

	_, err = db.Exec("INSERT INTO tables (table_no, rest_id, qr) VALUES (?, ?, ?)", table.TableNo, table.RestID, table.QR)
	defer db.Close()

	if err != nil {
		fmt.Println("Err", err.Error())
		return false, gin.H{"status": http.StatusBadRequest, "message": "DB Query Error! Add Table", "error": err.Error()}
	}

	return true, gin.H{"status": http.StatusOK, "message": "success", "qr": table.QR}
}

func CreateQRCode(tid int64, rid int64) ([]byte, gin.H) {
	var table = models.Table{}
	var err error
	table.QR, err = qrcode.Encode(string(tid)+":"+string(rid), qrcode.Medium, 256)
	if err != nil {
		return nil, gin.H{"status": http.StatusBadRequest, "message": "QR Code Error! Create QR Code"}
	}
	return table.QR, gin.H{"status": http.StatusOK, "message": "success"}
}
