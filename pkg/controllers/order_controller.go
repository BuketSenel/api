package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/SelfServiceCo/api/pkg/models"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func GetRestaurantOrders(rid int64) ([]models.Order, gin.H) {
	orders := []models.Order{}

	db, err := sql.Open("mysql", conf.Name+":"+conf.Password+"@tcp("+conf.Db+":3306)/selfservice?parseTime=true")

	if err != nil {
		return orders, gin.H{"status": http.StatusBadRequest, "message": "DB Connection Error!"}
	}

	results, err := db.Query("SELECT order_item_id, PNAME, PDESC, table_id, quantity FROM products JOIN orders ON `orders`.`product_id` = `products`.`ID` WHERE `orders`.`rest_id` = (?) AND `orders`.`order_status` != 'Deny'", rid)
	if err != nil {
		return orders, gin.H{"status": http.StatusBadRequest, "message": "Selection Error!"}
	}

	for results.Next() {
		order := models.Order{}
		err = results.Scan(&order.OrderItemID, &order.PName, &order.PDesc, &order.TableID, &order.Quantity)
		if err != nil {
			return orders, gin.H{"status": http.StatusBadRequest, "message": "Scan Error!", "data": results, "Error": err.Error()}
		}
		orders = append(orders, order)
	}
	defer db.Close()

	return orders, gin.H{"status": "success", "data": orders}
}

func GetOrdersByUser(uid int64) ([]models.Order, gin.H) {
	orders := []models.Order{}

	db, err := sql.Open("mysql", conf.Name+":"+conf.Password+"@tcp("+conf.Db+":3306)/selfservice?parseTime=true")

	if err != nil {
		return orders, gin.H{"status": http.StatusBadRequest, "message": "DB Connection Error!"}
	}

	results, err := db.Query("SELECT * FROM orders WHERE userId=?", uid)
	defer db.Close()
	if err != nil {
		return orders, gin.H{"status": http.StatusBadRequest, "message": "Selection Error!"}
	}

	for results.Next() {
		order := models.Order{}
		err = results.Scan(&order.ID, &order.UserID, &order.ResID, &order.TableID, &order.Details, &order.Status, &order.Created, &order.Updated)
		if err != nil {
			return orders, gin.H{"status": http.StatusBadRequest, "message": "Scan Error!", "data": results, "Error": err.Error()}
		}
		orders = append(orders, order)
	}

	return orders, gin.H{"status": "success", "data": orders}
}

func ChangeOrderStatus(c *gin.Context) (bool, gin.H) {
	db, err := sql.Open("mysql", conf.Name+":"+conf.Password+"@tcp("+conf.Db+":3306)/selfservice")

	if err != nil {
		return false, gin.H{"status": http.StatusBadRequest, "message": "DB Connection Error!"}
	}

	order := models.Order{}
	if err := c.BindJSON(&order); err != nil {
		return false, gin.H{"status": http.StatusBadRequest, "message": "JSON Bind Error!"}
	}

	result, err := db.Exec("UPDATE orders SET orderStatus = ? WHERE restId = ? AND ID = ?", order.Status, order.ResID, order.ID)
	defer db.Close()

	if err != nil {
		return false, gin.H{"status": http.StatusBadRequest, "message": "Update Error!"}
	}

	return true, gin.H{"status": "success", "data": result}
}

func GetOrder(oid int64, rid int64) ([]models.Order, gin.H) {
	order := []models.Order{}
	db, err := sql.Open("mysql", conf.Name+":"+conf.Password+"@tcp("+conf.Db+":3306)/selfservice")
	if err != nil {
		return order, gin.H{"status": http.StatusBadRequest, "message": "DB Connection Error!"}
	}

	results, err := db.Query("SELECT order_item_id, PNAME, PDESC, table_id, quantity, orderStatus FROM products JOIN orders ON `orders`.`product_id` = `products`.`ID` WHERE `orders`.`rest_id` = (?) AND `orders`.`order_id` = (?)", rid, oid)
	if err != nil {
		return order, gin.H{"status": http.StatusBadRequest, "message": "Selection Error!"}

	}

	for results.Next() {
		var ord models.Order
		err = results.Scan(&ord.OrderItemID, &ord.ProductName, &ord.ProductDesc, &ord.TableID, &ord.Quantity, &ord.Status)
		if err != nil {
			return order, gin.H{"status": http.StatusBadRequest, "message": "Get Order Query Error!"}
		}
		order = append(order, ord)
	}
	defer db.Close()
	return order, gin.H{"status": "success", "data": order}
}

func CreateOrder(c *gin.Context) (bool, gin.H) {
	orderRequest := models.Order{}
	products := []models.Product{}

	if err := c.BindJSON(&orderRequest); err != nil {
		c.AbortWithError(401, err)
	}
	if err := json.Unmarshal([]byte(orderRequest.Details), &products); err != nil {
		c.AbortWithError(401, err)
	}
	db, err := sql.Open("mysql", conf.Name+":"+conf.Password+"@tcp("+conf.Db+":3306)/selfservice")
	if err != nil {
		return false, gin.H{"status": http.StatusBadRequest, "message": "Database Connection Error!"}
	}

	for i := 0; i < len(products); i++ {
		results, err := db.Query("INSERT INTO orders (order_id, rest_id, table_id, user_id, product_id, price, quantity, order_status) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", orderRequest.ID, orderRequest.ResID, orderRequest.TableID, orderRequest.UserID, products[i].ID, products[i].Price, products[i].Quantity, orderRequest.Status)
		if err != nil {
			return false, gin.H{"status": http.StatusBadRequest, "message": "Insertion Error!", "data": err.Error(), "results": results}
		}
	}
	defer db.Close()
	return true, gin.H{"status": http.StatusOK, "data": "Order Created!"}
}
