package controllers

import (
	"database/sql"
	"net/http"

	"github.com/SelfServiceCo/api/pkg/models"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func GetRestaurantOrders(rid int64) ([]models.Order, gin.H) {
	orders := []models.Order{}

	db, err := sql.Open("mysql", conf.Name+":"+conf.Password+"@tcp("+conf.Db+":3306)/selfservice")

	if err != nil {
		return orders, gin.H{"status": http.StatusBadRequest, "message": "DB Connection Error!"}
	}

	results, err := db.Query("SELECT * FROM orders WHERE RID = ? AND orderStatus != 'DONE' AND orderStatus != 'DENY'", rid)
	defer db.Close()

	if err != nil {
		return orders, gin.H{"status": http.StatusBadRequest, "message": "Selection Error!"}
	}

	for results.Next() {
		var order models.Order
		err = results.Scan(&order.ID, &order.ResID, &order.UserID, &order.TableID, &order.Details, &order.Status)
		if err != nil {
			return orders, gin.H{"status": http.StatusBadRequest, "message": "Scan Error!"}
		}
		orders = append(orders, order)
	}

	return orders, gin.H{"status": "success", "data": orders}
}

func ChangeOrderStatus(oid int64, rid int64, status string) (bool, gin.H) {
	db, err := sql.Open("mysql", conf.Name+":"+conf.Password+"@tcp("+conf.Db+":3306)/selfservice")

	if err != nil {
		return false, gin.H{"status": http.StatusBadRequest, "message": "DB Connection Error!"}
	}

	result, err := db.Exec("UPDATE orders SET orderStatus = ? WHERE RID = ? AND OID = ?", status, rid, oid)
	defer db.Close()

	if err != nil {
		return false, gin.H{"status": http.StatusBadRequest, "message": "Update Error!"}
	}

	return true, gin.H{"status": "success", "data": result}
}

func GetOrdersByUser(uid int64) ([]models.Order, gin.H) {
	orders := []models.Order{}

	db, err := sql.Open("mysql", conf.Name+":"+conf.Password+"@tcp("+conf.Db+":3306)/selfservice")

	if err != nil {
		return orders, gin.H{"status": http.StatusBadRequest, "message": "DB Connection Error!"}
	}

	results, err := db.Query("SELECT * FROM orders WHERE UID = ?", uid)
	defer db.Close()

	if err != nil {
		return orders, gin.H{"status": http.StatusBadRequest, "message": "Selection Error!"}
	}

	for results.Next() {
		var order models.Order
		err = results.Scan(&order.ID, &order.ResID, &order.UserID, &order.TableID, &order.Details, &order.Status)
		if err != nil {
			return orders, gin.H{"status": http.StatusBadRequest, "message": "Scan Error!"}
		}
		orders = append(orders, order)
	}

	return orders, gin.H{"status": "success", "data": orders}
}

func GetOrder(oid int64, rid int64) ([]models.Order, gin.H) {
	order := []models.Order{}
	db, err := sql.Open("mysql", conf.Name+":"+conf.Password+"@tcp("+conf.Db+":3306)/selfservice")
	if err != nil {
		return order, gin.H{"status": http.StatusBadRequest, "message": "DB Connection Error!"}
	}

	results, err := db.Query("SELECT * FROM orders WHERE ID = ? AND RID = ?", oid, rid)
	if err != nil {
		return order, gin.H{"status": http.StatusBadRequest, "message": "Selection Error!"}

	}

	for results.Next() {
		var ord models.Order
		err = results.Scan(&ord.ID, &ord.UserID, &ord.ResID, &ord.TableID, &ord.Details, &ord.Status)
		if err != nil {
			return order, gin.H{"status": http.StatusBadRequest, "message": "Get Order Query Error!"}
		}
		order = append(order, ord)
	}
	return order, gin.H{"status": "success", "data": order}
}

func CreateOrder(c *gin.Context) (bool, gin.H) {
	var orderRequest models.Order
	if err := c.BindJSON(&orderRequest); err != nil {
		c.AbortWithError(401, err)
	}
	db, err := sql.Open("mysql", conf.Name+":"+conf.Password+"@tcp("+conf.Db+":3306)/selfservice")
	if err != nil {
		return false, gin.H{"status": http.StatusBadRequest, "message": "Database Connection Error!"}
	}
	results, err := db.Query("INSERT INTO orders (ID, user_id, RID, table_id, details, orderStatus) VALUES (?, ?, ?, ?, ?, ?)", 1, orderRequest.UserID, orderRequest.ResID, orderRequest.TableID, orderRequest.Details, orderRequest.Status)
	if err != nil {
		return false, gin.H{"status": http.StatusBadRequest, "message": "Insertion Error!"}
	}
	return true, gin.H{"status": "success", "data": results}
}
