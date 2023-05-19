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

	db := CreateConnection()

	if db == nil {
		return nil, gin.H{"status": http.StatusBadRequest, "message": "DB Connection Error!"}
	}

	results, err := db.Query("SELECT order_item_id, prod_name, prod_desc, table_id, prod_count, order_status, order_created_at FROM products JOIN orders ON `orders`.`prod_id` = products.prod_id WHERE `orders`.`rest_id` = (?) AND (`orders`.`order_status` = 'To do' OR  `orders`.`order_status` = 'In progress' OR `orders`.`order_status` = 'Completed')", rid)
	if err != nil {
		return nil, gin.H{"status": http.StatusBadRequest, "message": "Selection Error!"}
	}

	for results.Next() {
		cq := models.CustomQuery{}
		err = results.Scan(&cq.OrderItemID, &cq.ProductName, &cq.ProductDescription, &cq.TableID, &cq.Quantity, &cq.Status, &cq.OrderCreatedAt)
		if err != nil {
			return nil, gin.H{"status": http.StatusBadRequest, "message": "Scan Error!", "data": results, "Error": err.Error()}
		}
		customQuery = append(customQuery, cq)
	}
	CloseConnection(db)

	return &customQuery, gin.H{"status": "success", "data": customQuery}
}

func GetOrdersByUser(uid int64, status string, rid int64) ([]models.Order, gin.H) {
	orders := []models.Order{}

	db := CreateConnection()

	if db == nil {
		return orders, gin.H{"status": http.StatusBadRequest, "message": "DB Connection Error!"}
	}
	var results *sql.Rows
	var err error
	if status == "" {
		results, err = db.Query("SELECT user_id, order_id, order_item_id, prod_name, table_id, prod_count, order_status, orders.rest_id, products.price FROM products JOIN orders ON `orders`.`prod_id` = products.prod_id WHERE `orders`.`user_id` = (?) AND (`orders`.`order_status` = 'To do' OR  `orders`.`order_status` = 'In progress' OR `orders`.`order_status` = 'Completed') AND  `orders`.`rest_id` = (?)", uid, rid)
	} else {
		results, err = db.Query("SELECT user_id, order_id, order_item_id, prod_name, table_id, prod_count, order_status, orders.rest_id, products.price FROM products JOIN orders ON `orders`.`prod_id` = products.prod_id WHERE `orders`.`user_id` = (?) AND `orders`.`order_status` = (?) AND `orders`.`rest_id` = (?)", uid, status, rid)
	}
	CloseConnection(db)
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
	db := CreateConnection()

	if db == nil {
		return false, gin.H{"status": http.StatusBadRequest, "message": "DB Connection Error!"}
	}

	order := models.Order{}
	if err := c.BindJSON(&order); err != nil {
		return false, gin.H{"status": http.StatusBadRequest, "message": "JSON Bind Error! Change Order Status"}
	}

	result, err := db.Exec("UPDATE orders SET order_status = ? WHERE rest_id = ? AND order_item_id = ?", order.Status, order.ResID, order.OrderItemID)

	if err != nil {
		return false, gin.H{"status": http.StatusBadRequest, "message": "Update Error!"}
	}

	var availability sql.Result = nil
	if order.Status == "deny" {
		availability, err = db.Exec("UPDATE products SET availability = 0 WHERE rest_id = ? AND prod_id = ?", order.ResID, order.ProductId)
		if err != nil {
			return false, gin.H{"status": http.StatusBadRequest, "message": "Update Error!"}
		}
	}
	CloseConnection(db)

	return true, gin.H{"status": "success", "data": result, "availability": availability}
}

func GetOrder(oid int64, rid int64) (*[]models.CustomQuery, gin.H) {
	customQuery := []models.CustomQuery{}
	db := CreateConnection()
	if db == nil {
		return nil, gin.H{"status": http.StatusBadRequest, "message": "DB Connection Error!"}
	}

	results, err := db.Query("SELECT order_item_id, prod_name, prod_desc, table_id, prod_count, order_status FROM products JOIN orders ON orders.prod_id = products.prod_id WHERE `orders`.`rest_id` = (?) AND `orders`.`order_id` = (?)", rid, oid)
	if err != nil {
		return nil, gin.H{"status": http.StatusBadRequest, "message": "Selection Error!"}

	}

	CloseConnection(db)

	for results.Next() {
		cq := models.CustomQuery{}
		err = results.Scan(&cq.OrderItemID, &cq.ProductName, &cq.ProductDescription, &cq.TableID, &cq.Quantity, &cq.Status)
		if err != nil {
			return nil, gin.H{"status": http.StatusBadRequest, "message": "Get Order Query Error!"}
		}
		customQuery = append(customQuery, cq)
	}

	return &customQuery, gin.H{"status": "success", "data": customQuery}
}

func CreateOrder(c *gin.Context) (bool, gin.H) {
	orderRequest := models.Order{}
	products := []models.Product{}

	if err := c.BindJSON(&orderRequest); err != nil {
		return false, gin.H{"status": http.StatusBadRequest, "message": "Bind Error! Create Order", "error": err.Error()}
	}
	if err := json.Unmarshal([]byte(orderRequest.Details), &products); err != nil {
		return false, gin.H{"status": http.StatusBadRequest, "message": "Unmarshal Error! Create Order", "data": err.Error()}
	}
	db := CreateConnection()
	if db == nil {
		return false, gin.H{"status": http.StatusBadRequest, "message": "Database Connection Error!"}
	}

	for i := 0; i < len(products); i++ {
		results, err := db.Query("INSERT INTO orders (order_id, rest_id, table_id, user_id, prod_id, price, prod_count, order_status) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", orderRequest.ID, orderRequest.ResID, orderRequest.TableID, orderRequest.UserID, products[i].ID, products[i].Price, products[i].Quantity, orderRequest.Status)
		if err != nil {
			return false, gin.H{"status": http.StatusBadRequest, "message": "Insertion Error!", "data": err.Error(), "results": results}
		}
	}
	CloseConnection(db)

	return true, gin.H{"status": http.StatusOK, "data": "Order Created!"}
}

func GetPopularOrders(rid int64) (bool, gin.H) {
	orders := []models.CustomQuery{}
	db := CreateConnection()
	if db == nil {
		return false, gin.H{"status": http.StatusBadRequest, "message": "DB Connection Error!"}
	}
	results, err := db.Query("SELECT p.prod_id, p.prod_name, p.prod_desc, p.cat_id, p.prod_image, p.price, p.currency, p.prep_dur_minute, SUM(o.prod_count) AS total_quantity FROM orders o JOIN products p ON o.prod_id = p.prod_id WHERE o.rest_id = ? GROUP BY o.prod_id ORDER BY total_quantity DESC LIMIT 5", rid)
	if err != nil {
		return false, gin.H{"status": http.StatusBadRequest, "message": err.Error()}
	}
	CloseConnection(db)
	for results.Next() {
		order := models.CustomQuery{}
		err = results.Scan(&order.ProductID, &order.ProductName, &order.ProductDescription, &order.CatID, &order.ProductImage, &order.Price, &order.Currency, &order.PrepDurationMin, &order.OrderItemTotalQty)
		if err != nil {
			return false, gin.H{"status": http.StatusBadRequest, "message": "Get Popular Orders Query Error!"}
		}
		orders = append(orders, order)
	}
	return true, gin.H{"status": "success", "data": orders}
}

func GetDailyOrders(rid int64) (bool, gin.H) {
	orders := []models.CustomQuery{}
	db := CreateConnection()
	if db == nil {
		return false, gin.H{"status": http.StatusBadRequest, "message": "DB Connection Error! Get Daily Orders"}
	}
	results, err := db.Query("SELECT p.prod_id, p.prod_name, p.prod_desc, p.cat_id, p.prod_image, p.price, p.currency, p.prep_dur_minute, SUM(o.prod_count) AS total_quantity FROM orders o JOIN products p ON o.prod_id = p.prod_id WHERE o.rest_id = ? AND DATE(order_created_at) = CURDATE() GROUP BY o.prod_id ORDER BY total_quantity DESC", rid)
	if err != nil {
		return false, gin.H{"status": http.StatusBadRequest, "error": err.Error(), "message": "Selection Error! Get Daily Orders"}
	}
	CloseConnection(db)
	for results.Next() {
		order := models.CustomQuery{}
		err = results.Scan(&order.ProductID, &order.ProductName, &order.ProductDescription, &order.CatID, &order.ProductImage, &order.Price, &order.Currency, &order.PrepDurationMin, &order.OrderItemTotalQty)
		if err != nil {
			return false, gin.H{"status": http.StatusBadRequest, "message": "Get Daily Orders Query Error!", "error": err.Error()}
		}
		orders = append(orders, order)
	}
	return true, gin.H{"status": "success", "data": orders}
}

func GetWeeklyOrders(rid int64) (bool, gin.H) {
	orders := []models.CustomQuery{}
	db := CreateConnection()
	if db == nil {
		return false, gin.H{"status": http.StatusBadRequest, "message": "DB Connection Error! Get Weekly Orders"}
	}
	results, err := db.Query("SELECT p.prod_id, p.prod_name, p.prod_desc, p.cat_id, p.prod_image, p.price, p.currency, p.prep_dur_minute, SUM(o.prod_count) AS total_quantity FROM orders o JOIN products p ON o.prod_id = p.prod_id WHERE o.rest_id = ? AND order_created_at >= DATE_SUB(CURRENT_TIMESTAMP(), INTERVAL 7 DAY) GROUP BY o.prod_id ORDER BY total_quantity DESC", rid)
	if err != nil {
		return false, gin.H{"status": http.StatusBadRequest, "error": err.Error(), "message": "Selection Error! Get Weekly Orders"}
	}
	CloseConnection(db)
	for results.Next() {
		order := models.CustomQuery{}
		err = results.Scan(&order.ProductID, &order.ProductName, &order.ProductDescription, &order.CatID, &order.ProductImage, &order.Price, &order.Currency, &order.PrepDurationMin, &order.OrderItemTotalQty)
		if err != nil {
			return false, gin.H{"status": http.StatusBadRequest, "message": "Get Weekly Orders Query Error!", "error": err.Error()}
		}
		orders = append(orders, order)
	}
	return true, gin.H{"status": "success", "data": orders}
}

func GetMonthlyOrders(rid int64) (bool, gin.H) {
	orders := []models.CustomQuery{}
	db := CreateConnection()
	if db == nil {
		return false, gin.H{"status": http.StatusBadRequest, "message": "DB Connection Error! Get Monthly Orders"}
	}
	results, err := db.Query("SELECT p.prod_id, p.prod_name, p.prod_desc, p.cat_id, p.prod_image, p.price, p.currency, p.prep_dur_minute, SUM(o.prod_count) AS total_quantity FROM orders o JOIN products p ON o.prod_id = p.prod_id WHERE o.rest_id = ? AND order_created_at >= DATE_SUB(CURRENT_TIMESTAMP(), INTERVAL 30 DAY) GROUP BY o.prod_id ORDER BY total_quantity DESC", rid)
	if err != nil {
		return false, gin.H{"status": http.StatusBadRequest, "error": err.Error(), "message": "Selection Error! Get Monthly Orders"}
	}
	CloseConnection(db)
	for results.Next() {
		order := models.CustomQuery{}
		err = results.Scan(&order.ProductID, &order.ProductName, &order.ProductDescription, &order.CatID, &order.ProductImage, &order.Price, &order.Currency, &order.PrepDurationMin, &order.OrderItemTotalQty)
		if err != nil {
			return false, gin.H{"status": http.StatusBadRequest, "message": "Get Monthly Orders Query Error!", "error": err.Error()}
		}
		orders = append(orders, order)
	}
	return true, gin.H{"status": "success", "data": orders}
}
