package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/SelfServiceCo/api/pkg/models"
	"github.com/gin-gonic/gin"
	qrcode "github.com/skip2/go-qrcode"
)

func GetRestaurantTables(resID int64, tableId int64) ([]models.Table, gin.H) {

	db := CreateConnection()

	if db == nil {
		return nil, gin.H{"status": http.StatusBadRequest, "message": "DB Connection Error! Get Restaurant Tables"}
	}

	if tableId != 0 {
		results, err := db.Query("SELECT rest_id, table_no, waiter_id FROM tables WHERE rest_id = ? AND table_no = ?", resID, tableId)
		CloseConnection(db)

		if err != nil {
			return nil, gin.H{"status": http.StatusBadRequest, "message": "Query Error! Get Restaurant Tables"}
		}
		tables := []models.Table{}
		for results.Next() {
			var table models.Table
			err = results.Scan(&table.RestID, &table.TableNo, &table.WaiterID)
			if err != nil {
				return nil, gin.H{"status": http.StatusBadRequest, "message": "Scan Error! Get Restaurant Tables"}
			}
			tables = append(tables, table)
		}

		return tables, gin.H{"status": http.StatusOK, "message": "success", "data": tables}
	}
	results, err := db.Query("SELECT rest_id, table_no, waiter_id FROM tables WHERE rest_id = ?", resID)
	CloseConnection(db)
	if err != nil {
		return nil, gin.H{"status": http.StatusBadRequest, "message": "Query Error! Get Restaurant Tables", "error": err.Error()}
	}
	tables := []models.Table{}
	for results.Next() {
		var table models.Table
		err = results.Scan(&table.RestID, &table.TableNo, &table.WaiterID)
		if err != nil {
			return nil, gin.H{"status": http.StatusBadRequest, "message": "Scan Error! Get Restaurant Tables"}
		}
		tables = append(tables, table)
	}

	return tables, gin.H{"status": http.StatusOK, "message": "success", "data": tables}
}

func GetTable(tableId int64) (models.Table, gin.H) {
	var table models.Table
	db := CreateConnection()

	if db == nil {
		return table, gin.H{"status": http.StatusBadRequest, "message": "DB Connection Error! Get Table"}
	}

	results, err := db.Query("SELECT * FROM tables WHERE table_id = ?", tableId)
	CloseConnection(db)

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

	db := CreateConnection()

	if db == nil {
		return nil, gin.H{"status": http.StatusBadRequest, "message": "DB Connection Error! Orders By Table"}
	}

	results, err := db.Query("SELECT * FROM orders WHERE table_id = ?", tableId)
	CloseConnection(db)

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

	db := CreateConnection()

	if db == nil {
		return false, gin.H{"status": http.StatusBadRequest, "message": "DB Connection Error! Add Table"}
	}

	var header = gin.H{}
	domain, header := CreateQRCode(table.TableNo, table.RestID)
	if header["status"] != http.StatusOK {
		return false, gin.H{"status": http.StatusBadRequest, "message": "QR Code Error! Add Table", "error": header, "data": domain}
	}

	_, err = db.Exec("INSERT INTO tables (table_no, rest_id, qr) VALUES (?, ?, ?)", table.TableNo, table.RestID, domain)
	CloseConnection(db)

	if err != nil {
		return false, gin.H{"status": http.StatusBadRequest, "message": "DB Query Error! Add Table", "error": err.Error()}
	}

	return true, gin.H{"status": http.StatusOK, "message": "success", "qr": domain}
}

func CreateQRCode(tid int64, rid int64) (string, gin.H) {
	var table = models.Table{}
	var err error
	rest_id := strconv.FormatInt(rid, 10)
	table_no := strconv.FormatInt(tid, 10)
	qr_string := map[string]int64{"table_no": tid, "rest_id": rid}
	qr_json, _ := json.Marshal(qr_string)
	table.QR, err = qrcode.New(string(qr_json), qrcode.Medium)
	if err != nil {
		return "", gin.H{"status": http.StatusBadRequest, "message": "QR Code Error! Create QR Code"}
	}
	if err != nil {
		return "", gin.H{"status": http.StatusBadRequest, "message": "Could not retrive the file from the request!"}
	}
	filePath := "/www/qr-codes/" + rest_id + "/" + table_no + ".png"
	if _, err := os.Stat("/www/qr-codes/" + rest_id); os.IsNotExist(err) {
		os.Mkdir("/www/qr-codes/"+rest_id, 0777)
	}
	osFile, err := os.Create(filePath)
	if err != nil {
		return "", gin.H{"status": http.StatusBadRequest, "message": "Could not create the file!", "error": err}
	}
	err = table.QR.WriteFile(256, osFile.Name())
	if err != nil {
		return "", gin.H{"status": http.StatusBadRequest, "message": "QR Code Write File Error! Create QR Code", "error": err.Error()}
	}
	domain, header := UploadToS3(filePath)
	if header["status"] != http.StatusOK {
		return "", gin.H{"status": http.StatusBadRequest, "message": "QR Code Upload To S3 Error! Create QR Code", "error": header}
	}
	return domain, gin.H{"status": http.StatusOK, "message": "success"}
}

func GetQRCode(rid int64, tid int64) (string, gin.H) {
	db := CreateConnection()

	if db == nil {
		return "", gin.H{"status": http.StatusBadRequest, "message": "DB Connection Error! Get QR Code"}
	}
	results, err := db.Query("SELECT qr FROM tables WHERE table_no = (?) AND rest_id = (?)", tid, rid)
	if err != nil {
		return "", gin.H{"status": http.StatusBadRequest, "message": "DB Query Error! Get QR Code"}
	}
	CloseConnection(db)
	table := models.Table{}
	for results.Next() {
		err = results.Scan(&table.QRString)
		if err != nil {
			return "", gin.H{"status": http.StatusBadRequest, "message": "DB Scan Error! Get Users"}
		}
	}
	return table.QRString, gin.H{"status": http.StatusOK, "message": "success", "table": table}
}

func EditTable(c *gin.Context) (bool, gin.H) {
	table := models.Table{}
	err := c.BindJSON(&table)

	if err != nil {
		return false, gin.H{"status": http.StatusBadRequest, "message": "Bind Error! Edit Table"}
	}

	db := CreateConnection()

	if db == nil {
		return false, gin.H{"status": http.StatusBadRequest, "message": "DB Connection Error! Edit Table"}
	}

	_, err = db.Exec("UPDATE tables SET table_no = ? WHERE table_no = ? and rest_id = ?", table.NewTableNo, table.TableNo, table.RestID)
	CloseConnection(db)

	if err != nil {
		fmt.Println("Err", err.Error())
		return false, gin.H{"status": http.StatusBadRequest, "message": "DB Query Error! Edit Table", "error": err.Error()}
	}

	return true, gin.H{"status": http.StatusOK, "message": "success"}
}

func DeleteTable(c *gin.Context) (bool, gin.H) {
	table := models.Table{}
	err := c.BindJSON(&table)
	if err != nil {
		return false, gin.H{"status": http.StatusBadRequest, "message": "Bind Error! Delete Table"}
	}

	db := CreateConnection()

	if db == nil {
		return false, gin.H{"status": http.StatusBadRequest, "message": "DB Connection Error! Delete Table"}
	}

	results, err := db.Query("DELETE FROM tables WHERE rest_id = ? AND table_no = ?", table.RestID, table.TableNo)
	CloseConnection(db)
	if err != nil {
		return false, gin.H{"status": http.StatusBadRequest, "message": "Delete Error! Delete Table"}
	}
	return true, gin.H{"status": http.StatusOK, "message": "Table deleted!", "data": results}
}
