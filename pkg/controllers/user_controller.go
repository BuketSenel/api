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

	results, err := db.Query("SELECT user_id, user_name FROM users WHERE rest_id = ? and type = 'waiter'", rid)
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
	var user models.User
	err := c.BindJSON(&user)
	if err != nil {
		return false, gin.H{"status": http.StatusBadRequest, "message": "Invalid JSON"}
	}

	db := CreateConnection()

	if db == nil {
		return false, gin.H{"status": http.StatusBadRequest, "message": "DB Connection Error! Edit User"}
	}
	result, err := db.Exec("UPDATE users SET user_name = ?, user_phone = ?, email = ?, type = ? WHERE user_id = ?", user.Name, user.Phone, user.Email, user.Type, user.ID)
	CloseConnection(db)

	if err != nil {
		return false, gin.H{"status": http.StatusBadRequest, "message": "DB Query Error! Edit User"}
	}

	return true, gin.H{"status": http.StatusOK, "message": "success", "data": user, "result": result}
}
