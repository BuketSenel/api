package controllers

import (
	"net/http"

	"github.com/SelfServiceCo/api/pkg/models"
	"github.com/gin-gonic/gin"
)

func getUser(email string) (string, int64, int64, gin.H) {
	db := CreateConnection()
	if db == nil {
		return "", 0, 0, gin.H{"status": http.StatusBadRequest, "message": "Connection Error!!"}
	}
	result, err := db.Query("SELECT type, user_id, rest_id FROM users WHERE email = ? AND 1=1", email)

	if err != nil {
		return "", 0, 0, gin.H{"status": http.StatusBadRequest, "message": err.Error()}
	}
	CloseConnection(db)
	user := models.User{}
	for result.Next() {
		err = result.Scan(&user.Type, &user.ID, &user.ResID)
	}
	if err != nil {
		return "", 0, 0, gin.H{"status": http.StatusBadRequest, "message": err.Error()}
	}
	return user.Type, user.ID, user.ResID, gin.H{"status": http.StatusOK, "message": "Success"}
}

func GetRestaurantWaiters(rid int64) ([]models.User, gin.H) {
	var users = []models.User{}
	db := CreateConnection()

	if db == nil {
		return nil, gin.H{"status": http.StatusBadRequest, "message": "DB Connection Error! Get Waiters"}
	}

	results, err := db.Query("SELECT COUNT(tables.waiter_id), user_id, user_name FROM users JOIN tables ON users.user_id = tables.waiter_id WHERE users.rest_id = ? and type = 'waiter' GROUP BY tables.waiter_id", rid)
	CloseConnection(db)

	if err != nil {
		return nil, gin.H{"status": http.StatusBadRequest, "message": "DB Query Error! Get Waiters"}
	}

	for results.Next() {
		var user models.User
		err = results.Scan(&user.ID, &user.Name)
		if err != nil {
			return nil, gin.H{"status": http.StatusBadRequest, "message": "DB Scan Error! Get Waiters"}
		}
		users = append(users, user)
	}

	return users, gin.H{"status": http.StatusOK, "message": "success", "data": results}
}

func AssignWaiter(c *gin.Context) (bool, gin.H) {
	var table models.Table
	err := c.BindJSON(&table)
	if err != nil {
		return false, gin.H{"status": http.StatusBadRequest, "message": "Invalid JSON"}
	}

	db := CreateConnection()

	if db == nil {
		return false, gin.H{"status": http.StatusBadRequest, "message": "DB Connection Error! Assign Waiter"}
	}
	result, err := db.Exec("UPDATE tables SET waiter_id = ? WHERE table_no = ? AND rest_id = ?", table.WaiterID, table.TableNo, table.RestID)
	CloseConnection(db)

	if err != nil {
		return false, gin.H{"status": http.StatusBadRequest, "message": "DB Query Error! Assign Waiter"}
	}

	return true, gin.H{"status": http.StatusOK, "message": "success", "data": table, "result": result}
}

func GetUser(uid int64) ([]models.User, gin.H) {
	var users = []models.User{}
	db := CreateConnection()

	if db == nil {
		return nil, gin.H{"status": http.StatusBadRequest, "message": "DB Connection Error! Get Users"}
	}

	results, err := db.Query("SELECT user_name, user_phone, email, rest_id, type FROM users WHERE user_id = ?", uid)
	CloseConnection(db)

	if err != nil {
		return nil, gin.H{"status": http.StatusBadRequest, "message": "DB Query Error! Get Users", "error": err.Error()}
	}

	for results.Next() {
		var user models.User
		err = results.Scan(&user.Name, &user.Phone, &user.Email, &user.ResID, &user.Type)
		if err != nil {
			return nil, gin.H{"status": http.StatusBadRequest, "message": "DB Scan Error! Get Users"}
		}
		users = append(users, user)
	}

	return users, gin.H{"status": http.StatusOK, "message": "success", "data": results}
}

func EditUser(c *gin.Context) (bool, gin.H) {
	db := CreateConnection()

	if db == nil {
		return false, gin.H{"status": http.StatusBadRequest, "message": "DB Connection Error! Edit User"}
	}

	updateQuery := "UPDATE users SET "
	args := make([]interface{}, 0)
	var data map[string]interface{}

	if err := c.BindJSON(&data); err != nil {
		return false, gin.H{"status": http.StatusBadRequest, "message": "Bind Error! Edit User", "error": err.Error()}
	}

	for key, value := range data {
		updateQuery += key + " = ?, "
		args = append(args, value)
	}
	updateQuery = updateQuery[:len(updateQuery)-2]

	updateQuery += " WHERE user_id = ?"
	args = append(args, data["user_id"])

	result, err := db.Exec(updateQuery, args...)
	if err != nil {
		return false, gin.H{"status": http.StatusBadRequest, "message": "Update Error! Edit Product", "data": data, "error": err.Error()}
	}
	CloseConnection(db)

	if err != nil {
		return false, gin.H{"status": http.StatusBadRequest, "message": "DB Query Error! Edit User"}
	}

	return true, gin.H{"status": http.StatusOK, "message": "success", "data": data, "result": result}
}
