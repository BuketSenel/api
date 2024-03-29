package controllers

import (
	"net/http"

	"github.com/SelfServiceCo/api/pkg/models"
	"github.com/gin-gonic/gin"
)

func UserRegister(c *gin.Context) (bool, gin.H) {
	user := models.User{}

	if err := c.BindJSON(&user); err != nil {
		c.AbortWithError(401, err)
	}

	if user.Name == "" || user.Phone == "" || user.Email == "" || user.Password == "" {
		return false, gin.H{"error": "Please fill all the fields"}
	}

	result, err := SaveUser(user, c)

	if !result {
		return false, gin.H{"message": "Error", "status": http.StatusTeapot, "error": err, "data": result}
	}

	return true, gin.H{"message": "User registered successfully"}
}

func SaveUser(user models.User, c *gin.Context) (bool, gin.H) {
	credential := models.Credential{}
	db := CreateConnection()
	if db == nil {
		message := gin.H{"status": http.StatusBadGateway, "message": "DB Connection Error! Save User"}
		return false, message
	}

	rows, err := db.Query("SELECT * FROM users WHERE email = ? AND type = 'customer'", user.Email)

	if err != nil {
		return false, gin.H{"status": http.StatusBadRequest, "message": "Database error! Save User", "error": err.Error()}
	}

	if rows.Next() {
		return false, gin.H{"message": "Email address is already registered!", "status": http.StatusBadRequest, "error": err, "data": rows.Next()}
	}

	hashedPass := PasswordHash(user.Password)

	query := "INSERT INTO users (user_name, password, user_phone, email, rest_id, type) values (?,?,?,?,?,?)"
	results, err := db.ExecContext(c, query, user.Name, hashedPass, user.Phone, user.Email, user.ResID, user.Type)
	if err != nil {
		message := gin.H{"status": http.StatusBadRequest, "message": results}
		return false, message
	}

	credential.Email = user.Email
	credential.Password = hashedPass
	query = "INSERT INTO credentials (email, password) values (?,?)"
	results, err = db.ExecContext(c, query, string(user.Email), hashedPass)
	if err != nil {
		message := gin.H{"status": http.StatusBadRequest, "message": err.Error()}
		return false, message
	}
	CloseConnection(db)
	return true, gin.H{"message": "Success", "status": http.StatusOK, "data": results}
}
