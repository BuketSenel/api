package controllers

import (
	"net/http"

	"github.com/SelfServiceCo/api/pkg/models"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func GetRestaurant(id int64) ([]models.Restaurant, gin.H) {
	db := CreateConnection()

	if db == nil {
		return nil, gin.H{"status": http.StatusBadRequest, "message": "DB Connection Error! Get Restaurant"}
	}

	results, err := db.Query("SELECT * FROM restaurants WHERE rest_id = ?", id)

	if err != nil {
		return nil, gin.H{"status": http.StatusBadRequest, "message": "Selection Error! Get Restaurant"}
	}

	restaurant := []models.Restaurant{}
	for results.Next() {
		rest := models.Restaurant{}
		err = results.Scan(&rest.ID, &rest.Name, &rest.Summary, &rest.Logo, &rest.Address, &rest.District, &rest.City, &rest.Country, &rest.Phone, &rest.CreatedAt, &rest.UpdatedAt)
		if err != nil {
			return nil, gin.H{"status": http.StatusBadRequest, "message": "Scan Error! Get Restaurant", "data": results, "Error": err.Error()}
		}
		restaurant = append(restaurant, rest)
	}
	CloseConnection(db)

	return restaurant, gin.H{"status": http.StatusOK, "message": restaurant}
}

func GetTopRestaurants() ([]models.Restaurant, gin.H) {
	db := CreateConnection()

	if db == nil {
		return nil, gin.H{"status": http.StatusBadRequest, "message": "DB Connection Error! Get Top Restaurant"}
	}

	results, err := db.Query("SELECT * FROM restaurants")

	if err != nil {
		return nil, gin.H{"status": http.StatusBadRequest, "message": "Selection Error! Get Top Restaurant", "error": err.Error()}
	}
	CloseConnection(db)
	restaurants := []models.Restaurant{}
	for results.Next() {
		var rest models.Restaurant

		err = results.Scan(&rest.ID, &rest.Name, &rest.Summary, &rest.Logo, &rest.Address, &rest.District, &rest.City, &rest.Country, &rest.Phone, &rest.CreatedAt, &rest.UpdatedAt)
		if err != nil {
			return nil, gin.H{"status": http.StatusBadRequest, "message": "Scan Error! Get Top Restaurant", "data": results, "Error": err.Error()}
		}
		restaurants = append(restaurants, rest)
	}

	return restaurants, gin.H{"status": http.StatusOK, "message": restaurants}
}

func GetRestaurantStaff(rid int64) ([]models.User, gin.H) {
	db := CreateConnection()

	if db == nil {
		return nil, gin.H{"status": http.StatusBadRequest, "message": "DB Connection Error! Get Staff"}
	}

	results, err := db.Query("SELECT user_id, user_name, user_phone, email, rest_id, type, user_created_at  FROM users WHERE rest_id = ? AND type != 'customer'", rid)
	CloseConnection(db)

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

func DeleteUser(c *gin.Context) (bool, gin.H) {
	user := models.User{}
	db := CreateConnection()

	if db == nil {
		return false, gin.H{"status": http.StatusBadRequest, "message": "DB Connection Error! Delete User"}
	}
	if err := c.BindJSON(&user); err != nil {
		c.AbortWithError(401, err)
	}

	if user.ID == 0 {
		return false, gin.H{"error": "Please fill all the fields"}
	}

	results, err := db.Query("DELETE FROM users WHERE user_id = ? and type != 'customer'", user.ID)
	CloseConnection(db)
	if err != nil {
		return false, gin.H{"status": http.StatusBadRequest, "message": "Delete Error! Delete User"}
	}

	return true, gin.H{"message": "Success", "data": results}
}

func GetWaiterTables(rid int64, waiterID int64) ([]models.Table, gin.H) {
	db := CreateConnection()

	if db == nil {
		return nil, gin.H{"status": http.StatusBadRequest, "message": "DB Connection Error! Get Tables"}
	}

	results, err := db.Query("SELECT table_no FROM tables WHERE rest_id = ? and waiter_id = ?", rid, waiterID)
	CloseConnection(db)

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
	db := CreateConnection()

	if db == nil {
		return nil, gin.H{"status": http.StatusBadRequest, "message": "DB Connection Error! Get Orders"}
	}

	results, err := db.Query("SELECT table_id, prod_name, prep_dur_minute, order_status, order_item_id FROM products JOIN orders ON `orders`.`prod_id` = `products`.prod_id WHERE `orders`.table_id = (?) and `orders`.order_status NOT IN ('paid', 'deny') and `orders`.rest_id = (?)", tableID, rid)
	CloseConnection(db)

	if err != nil {
		return nil, gin.H{"status": http.StatusBadRequest, "message": err.Error()}
	}
	customQuery := []models.CustomQuery{}
	for results.Next() {
		var cq models.CustomQuery

		err = results.Scan(&cq.TableID, &cq.ProductName, &cq.PrepDurationMin, &cq.Status, &cq.OrderItemID)
		if err != nil {
			return nil, gin.H{"status": http.StatusBadRequest, "message": "Scan Error! Get Orders", "data": results, "Error": err.Error()}
		}
		customQuery = append(customQuery, cq)
	}

	return customQuery, gin.H{"status": http.StatusOK, "orders": customQuery}
}

func TipWaiter(c *gin.Context) (bool, gin.H) {
	tip := models.Tip{}
	db := CreateConnection()

	if db == nil {
		return false, gin.H{"status": http.StatusBadRequest, "message": "DB Connection Error! Tipping Waiter"}
	}

	var err error
	if err = c.BindJSON(&tip); err != nil {
		return false, gin.H{"status": http.StatusBadRequest, "message": "Binding Error! Tipping Waiter"}

	}

	_, err = db.Query("INSERT INTO tips (user_id, waiter_id, rest_id, tip) VALUES (?, ?, ?, ?)", tip.UserID, tip.WaiterID, tip.RestID, tip.Tip)
	CloseConnection(db)
	if err != nil {
		return false, gin.H{"status": http.StatusBadRequest, "message": "Insert Error! Tipping Waiter"}
	}

	return true, gin.H{"message": "Success"}
}

func GetTips(rid int64, wid int64) ([]models.Tip, gin.H) {
	db := CreateConnection()

	if db == nil {
		return nil, gin.H{"status": http.StatusBadRequest, "message": "DB Connection Error! Get Tips"}
	}

	results, err := db.Query("SELECT tip FROM tips WHERE rest_id = ? AND waiter_id = ?", rid, wid)
	CloseConnection(db)

	if err != nil {
		return nil, gin.H{"status": http.StatusBadRequest, "message": "Selection Error! Get Tips"}
	}

	tips := []models.Tip{}
	for results.Next() {
		var tip models.Tip
		err = results.Scan(&tip.Tip)
		if err != nil {
			return nil, gin.H{"status": http.StatusBadRequest, "message": "Scan Error! Get Tips", "data": results, "Error": err.Error()}
		}
		tips = append(tips, tip)
	}

	return tips, gin.H{"status": http.StatusOK, "message": tips}
}

func GetWaitersByTable(rid int64, tid int64) ([]models.Table, gin.H) {
	db := CreateConnection()

	if db == nil {
		return nil, gin.H{"status": http.StatusBadRequest, "message": "DB Connection Error! Get Waiters"}
	}

	results, err := db.Query("SELECT waiter_id, user_name FROM tables JOIN users ON `tables`.`rest_id` = ? AND `tables`.`table_no` = ? AND `tables`.`waiter_id` = users.user_id", rid, tid)
	CloseConnection(db)

	if err != nil {
		return nil, gin.H{"status": http.StatusBadRequest, "message": "Selection Error! Get Waiters"}
	}

	waiters := []models.Table{}
	for results.Next() {
		var waiter models.Table
		err = results.Scan(&waiter.WaiterID, &waiter.WaiterName)
		if err != nil {
			return nil, gin.H{"status": http.StatusBadRequest, "message": "Scan Error! Get Waiters", "data": results, "Error": err.Error()}
		}
		waiters = append(waiters, waiter)
	}

	return waiters, gin.H{"status": http.StatusOK, "message": waiters}

}

func EditRestaurant(c *gin.Context) (bool, gin.H) {
	db := CreateConnection()

	if db == nil {
		return false, gin.H{"status": http.StatusBadRequest, "message": "DB Connection Error! Edit Restaurant"}
	}

	updateQuery := "UPDATE restaurants SET "
	args := make([]interface{}, 0)
	var data map[string]interface{}

	if err := c.BindJSON(&data); err != nil {
		return false, gin.H{"status": http.StatusBadRequest, "message": "Bind Error! Edit Restaurant", "error": err.Error()}
	}

	for key, value := range data {
		updateQuery += key + " = ?, "
		args = append(args, value)
	}
	updateQuery = updateQuery[:len(updateQuery)-2]

	updateQuery += " WHERE rest_id = ?"
	args = append(args, data["rest_id"])

	result, err := db.Exec(updateQuery, args...)
	if err != nil {
		return false, gin.H{"status": http.StatusBadRequest, "message": "Update Error! Edit Restaurant", "data": data, "error": err.Error()}
	}
	CloseConnection(db)

	if err != nil {
		return false, gin.H{"status": http.StatusBadRequest, "message": "DB Query Error! Edit Restaurant"}
	}

	return true, gin.H{"status": http.StatusOK, "message": "success", "data": data, "result": result}
}
