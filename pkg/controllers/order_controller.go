package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/SelfServiceCo/api/pkg/models"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func GetRestaurantOrders(rid int64) (*[]models.CustomQuery, gin.H) {
	customQuery := []models.CustomQuery{}

	db, err := sql.Open("mysql", conf.Name+":"+conf.Password+"@tcp("+conf.Db+":3306)/selfservice?parseTime=true")

	if err != nil {
		return nil, gin.H{"status": http.StatusBadRequest, "message": "DB Connection Error!"}
	}

	results, err := db.Query("SELECT order_item_id, prod_name, prod_desc, table_id, quantity, order_status FROM products JOIN orders ON `orders`.`prod_id` = products.prod_id WHERE `orders`.`rest_id` = (?) AND (`orders`.`order_status` = 'To do' OR  `orders`.`order_status` = 'In progress' OR `orders`.`order_status` = 'Completed')", rid)
	if err != nil {
		return nil, gin.H{"status": http.StatusBadRequest, "message": "Selection Error!"}
	}

	for results.Next() {
		cq := models.CustomQuery{}
		err = results.Scan(&cq.OrderItemID, &cq.ProductName, &cq.ProductDescription, &cq.TableID, &cq.Quantity, &cq.Status)
		if err != nil {
			return nil, gin.H{"status": http.StatusBadRequest, "message": "Scan Error!", "data": results, "Error": err.Error()}
		}
		customQuery = append(customQuery, cq)
	}
	defer db.Close()

	return &customQuery, gin.H{"status": "success", "data": customQuery}
}

func GetOrdersByUser(uid int64, status string, rid int64) ([]models.Order, gin.H) {
	orders := []models.Order{}

	db, err := sql.Open("mysql", conf.Name+":"+conf.Password+"@tcp("+conf.Db+":3306)/selfservice?parseTime=true")

	if err != nil {
		return orders, gin.H{"status": http.StatusBadRequest, "message": "DB Connection Error!"}
	}
	var results *sql.Rows
	if status == "" {
		results, err = db.Query("SELECT user_id, order_id, order_item_id, prod_name, table_id, quantity, order_status, orders.rest_id, products.price FROM products JOIN orders ON `orders`.`prod_id` = products.prod_id WHERE `orders`.`user_id` = (?) AND (`orders`.`order_status` = 'To do' OR  `orders`.`order_status` = 'In progress' OR `orders`.`order_status` = 'Completed') AND  `orders`.`rest_id` = (?)", uid, rid)
	} else {
		results, err = db.Query("SELECT user_id, order_id, order_item_id, prod_name, table_id, quantity, order_status, orders.rest_id, products.price FROM products JOIN orders ON `orders`.`prod_id` = products.prod_id WHERE `orders`.`user_id` = (?) AND `orders`.`order_status` = (?) AND `orders`.`rest_id` = (?)", uid, status, rid)
	}
	defer db.Close()
	if err != nil {
		return orders, gin.H{"status": http.StatusBadRequest, "message": "Selection Error!"}
	}

	var count int64
	count = 0
	for results.Next() {
		order := models.Order{}
		err = results.Scan(&order.UserID, &order.ID, &order.OrderItemID, &order.ProductName, &order.TableID, &order.Quantity, &order.Status, &order.ResID, &order.Price)
		if err != nil {
			return orders, gin.H{"status": http.StatusBadRequest, "message": "Scan Error!", "data": results, "Error": err.Error()}
		}
		orders = append(orders, order)
		quantity, _ := strconv.ParseInt(order.Quantity, 10, 64)
		if quantity > 1 {
			count = count + (quantity - 1)
		}
		count++
	}

	return orders, gin.H{"status": "success", "data": orders, "results": results, "count": count}
}

func ChangeOrderStatus(c *gin.Context) (bool, gin.H) {
	db, err := sql.Open("mysql", conf.Name+":"+conf.Password+"@tcp("+conf.Db+":3306)/selfservice")

	if err != nil {
		return false, gin.H{"status": http.StatusBadRequest, "message": "DB Connection Error!"}
	}

	order := models.Order{}
	if err := c.BindJSON(&order); err != nil {
		return false, gin.H{"status": http.StatusBadRequest, "message": "JSON Bind Error! Change Order Status"}
	}

	result, err := db.Exec("UPDATE orders SET order_status = ? WHERE rest_id = ? AND order_item_id = ?", order.Status, order.ResID, order.OrderItemID)
	defer db.Close()

	if err != nil {
		return false, gin.H{"status": http.StatusBadRequest, "message": "Update Error!"}
	}

	return true, gin.H{"status": "success", "data": result}
}

func GetOrder(oid int64, rid int64) (*[]models.CustomQuery, gin.H) {
	customQuery := []models.CustomQuery{}
	db, err := sql.Open("mysql", conf.Name+":"+conf.Password+"@tcp("+conf.Db+":3306)/selfservice")
	if err != nil {
		return nil, gin.H{"status": http.StatusBadRequest, "message": "DB Connection Error!"}
	}

	results, err := db.Query("SELECT order_item_id, prod_name, prod_desc, table_id, quantity, order_status FROM products JOIN orders ON orders.prod_id = products.prod_id WHERE `orders`.`rest_id` = (?) AND `orders`.`order_id` = (?)", rid, oid)
	if err != nil {
		return nil, gin.H{"status": http.StatusBadRequest, "message": "Selection Error!"}

	}

	for results.Next() {
		cq := models.CustomQuery{}
		err = results.Scan(&cq.OrderItemID, &cq.ProductName, &cq.ProductDescription, &cq.TableID, &cq.Quantity, &cq.Status)
		if err != nil {
			return nil, gin.H{"status": http.StatusBadRequest, "message": "Get Order Query Error!"}
		}
		customQuery = append(customQuery, cq)
	}
	defer db.Close()
	return &customQuery, gin.H{"status": "success", "data": customQuery}
}

func CreateOrder(c *gin.Context) (bool, gin.H) {
	orderRequest := models.Order{}
	products := []models.Product{}

	if err := c.BindJSON(&orderRequest); err != nil {
		return false, gin.H{"status": http.StatusBadRequest, "message": "Bind Error! Create Order"}
	}
	if err := json.Unmarshal([]byte(orderRequest.Details), &products); err != nil {
		return false, gin.H{"status": http.StatusBadRequest, "message": "Unmarshal Error! Create Order"}
	}
	db, err := sql.Open("mysql", conf.Name+":"+conf.Password+"@tcp("+conf.Db+":3306)/selfservice")
	if err != nil {
		return false, gin.H{"status": http.StatusBadRequest, "message": "Database Connection Error!"}
	}

	for i := 0; i < len(products); i++ {
		results, err := db.Query("INSERT INTO orders (order_id, rest_id, table_id, user_id, prod_id, price, quantity, order_status) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", orderRequest.ID, orderRequest.ResID, orderRequest.TableID, orderRequest.UserID, products[i].ID, products[i].Price, products[i].Quantity, orderRequest.Status)
		if err != nil {
			return false, gin.H{"status": http.StatusBadRequest, "message": "Insertion Error!", "data": err.Error(), "results": results}
		}
	}
	defer db.Close()
	return true, gin.H{"status": http.StatusOK, "data": "Order Created!"}
}
