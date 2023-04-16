package controllers

import (
	"database/sql"
	"net/http"

	"github.com/SelfServiceCo/api/pkg/drivers"
	"github.com/SelfServiceCo/api/pkg/models"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var selfdb = "selfservicedb"
var conf = drivers.MysqlConfigLoad()

func GetRestaurant(id int64) ([]models.Restaurant, gin.H) {
	db, err := sql.Open("mysql", conf.Name+":"+conf.Password+"@tcp("+conf.Db+":3306)/selfservice?parseTime=true")

	if err != nil {
		return nil, gin.H{"status": http.StatusBadRequest, "message": "DB Connection Error! Get Restaurant"}
	}

	results, err := db.Query("SELECT * FROM restaurants WHERE rest_id = ?", id)

	if err != nil {
		return nil, gin.H{"status": http.StatusBadRequest, "message": "Selection Error! Get Restaurant"}
	}

	restaurant := []models.Restaurant{}
	for results.Next() {
		rest := models.Restaurant{}
		err = results.Scan(&rest.ID, &rest.Name, &rest.Summary, &rest.Logo, &rest.Address, &rest.District, &rest.City, &rest.Country, &rest.Phone, &rest.Tags, &rest.CreatedAt, &rest.UpdatedAt)
		if err != nil {
			return nil, gin.H{"status": http.StatusBadRequest, "message": "Scan Error! Get Restaurant", "data": results, "Error": err.Error()}
		}
		restaurant = append(restaurant, rest)
	}
	defer db.Close()

	return restaurant, gin.H{"status": http.StatusOK, "message": restaurant}
}

func GetTopRestaurants() ([]models.Restaurant, gin.H) {
	db, err := sql.Open("mysql", conf.Name+":"+conf.Password+"@tcp("+conf.Db+":3306)/selfservice?parseTime=true")

	if err != nil {
		return nil, gin.H{"status": http.StatusBadRequest, "message": "DB Connection Error! Get Top Restaurant"}
	}

	results, err := db.Query("SELECT * FROM restaurants")

	if err != nil {
		return nil, gin.H{"status": http.StatusBadRequest, "message": "Selection Error! Get Top Restaurant"}
	}

	restaurants := []models.Restaurant{}
	for results.Next() {
		var rest models.Restaurant

		err = results.Scan(&rest.ID, &rest.Name, &rest.Summary, &rest.Logo, &rest.Address, &rest.District, &rest.City, &rest.Country, &rest.Phone, &rest.Tags, &rest.CreatedAt, &rest.UpdatedAt)
		if err != nil {
			return nil, gin.H{"status": http.StatusBadRequest, "message": "Scan Error! Get Top Restaurant", "data": results, "Error": err.Error()}
		}
		restaurants = append(restaurants, rest)
	}
	defer db.Close()

	return restaurants, gin.H{"status": http.StatusOK, "message": restaurants}
}

func GetRestaurantStaff(rid int64) ([]models.User, gin.H) {
	db, err := sql.Open("mysql", conf.Name+":"+conf.Password+"@tcp("+conf.Db+":3306)/selfservice?parseTime=true")

	if err != nil {
		return nil, gin.H{"status": http.StatusBadRequest, "message": "DB Connection Error! Get Staff"}
	}

	results, err := db.Query("SELECT user_id, user_name, user_phone, email, rest_id, type, user_created_at  FROM users WHERE rest_id = ?", rid)
	defer db.Close()

	if err != nil {
		return nil, gin.H{"status": http.StatusBadRequest, "message": "Selection Error! Get Staff"}
	}

	staff := []models.User{}
	for results.Next() {
		var user models.User
		err = results.Scan(&user.ID, &user.Name, &user.Phone, &user.Email, &user.ResID, &user.Type, &user.CreatedAt)
		if err != nil {
			return nil, gin.H{"status": http.StatusBadRequest, "message": "Scan Error! Get Staff", "data": results, "Error": err.Error()}
		}
		staff = append(staff, user)
	}

	return staff, gin.H{"status": http.StatusOK, "message": staff}
}

func AddStaff(c *gin.Context) (bool, gin.H) {
	user := models.User{}

	if err := c.BindJSON(&user); err != nil {
		c.AbortWithError(401, err)
	}

	if user.Name == "" || user.Phone == "" || user.Email == "" || user.Password == "" || user.ResID == 0 || user.Type == "" {
		return false, gin.H{"error": "Please fill all the fields"}
	}

	result, err := SaveUser(user, c)

	if !result || err != nil {
		return result, err
	}

	return true, gin.H{"message": "Staff registered successfully"}
}

func DeleteStaff(c *gin.Context) (bool, gin.H) {
	user := models.User{}

	if err := c.BindJSON(&user); err != nil {
		c.AbortWithError(401, err)
	}

	if user.ID == 0 {
		return false, gin.H{"error": "Please fill all the fields"}
	}

	//result, err := DeleteUser(user, c)

	//if !result || err != nil {
	//	return result, err
	//}

	return true, gin.H{"message": "Staff deleted successfully"}
}

func GetWaiterTables(rid int64, waiterID int64) ([]models.Table, gin.H) {
	db, err := sql.Open("mysql", conf.Name+":"+conf.Password+"@tcp("+conf.Db+":3306)/selfservice?parseTime=true")

	if err != nil {
		return nil, gin.H{"status": http.StatusBadRequest, "message": "DB Connection Error! Get Tables"}
	}

	results, err := db.Query("SELECT table_no FROM tables WHERE rest_id = ? and waiter_id = ?", rid, waiterID)
	defer db.Close()

	if err != nil {
		return nil, gin.H{"status": http.StatusBadRequest, "message": "Selection Error! Get Tables"}
	}

	tables := []models.Table{}
	for results.Next() {
		var table models.Table
		err = results.Scan(&table.TableNo)
		if err != nil {
			return nil, gin.H{"status": http.StatusBadRequest, "message": "Scan Error! Get Tables", "data": results, "Error": err.Error()}
		}
		tables = append(tables, table)
	}

	return tables, gin.H{"status": http.StatusOK, "message": tables}
}

func GetWaiterOrdersByTable(rid int64, tableID int64) ([]models.CustomQuery, gin.H) {
	db, err := sql.Open("mysql", conf.Name+":"+conf.Password+"@tcp("+conf.Db+":3306)/selfservice?parseTime=true")

	if err != nil {
		return nil, gin.H{"status": http.StatusBadRequest, "message": "DB Connection Error! Get Orders"}
	}

	results, err := db.Query("SELECT prod_name, prep_dur_minute, order_status, order_item_id FROM products JOIN orders ON `orders`.`prod_id` = products.prod_id WHERE `orders`.table_id = (?) and `orders`.order_status NOT IN ('done', 'paid') and `orders`.rest_id = ?", tableID, rid)

	if err != nil {
		return nil, gin.H{"status": http.StatusBadRequest, "message": err.Error()}
	}
	customQuery := []models.CustomQuery{}
	for results.Next() {
		var cq models.CustomQuery

		err = results.Scan(&cq.ProductName, &cq.PrepDurationMin, &cq.Status, &cq.OrderItemID)
		if err != nil {
			return nil, gin.H{"status": http.StatusBadRequest, "message": "Scan Error! Get Orders", "data": results, "Error": err.Error()}
		}
		customQuery = append(customQuery, cq)
	}

	return customQuery, gin.H{"status": http.StatusOK, "orders": customQuery}
}
