package controllers

import (
	"fmt"
	"net/http"

	"github.com/SelfServiceCo/api/pkg/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func RestaurantRegister(c *gin.Context) (bool, gin.H) {
	custQuery := models.CustomQuery{}

	if err := c.BindJSON(&custQuery); err != nil {
		c.AbortWithError(401, err)
	}

	if custQuery.RestName == "" || custQuery.Address == "" || custQuery.District == "" || custQuery.City == "" || custQuery.Country == "" || custQuery.Email == "" || custQuery.Password == "" || custQuery.RestPhone == "" || custQuery.UserPhone == "" {
		fmt.Println("Registration Error!")
		return false, gin.H{"status": http.StatusBadRequest, "message": "Please fill all the fields"}
	}

	result, err := SaveRestaurant(custQuery, c)
	if !result {
		return false, gin.H{"status": http.StatusBadRequest, "message": "Registration Error!", "error": err, "data": result}
	}

	return true, gin.H{"status": http.StatusOK, "message": "Registration Error!"}
}

func SaveRestaurant(cq models.CustomQuery, c *gin.Context) (bool, gin.H) {
	credential := models.Credential{}
	db := CreateConnection()
	if db == nil {
		return false, gin.H{"status": http.StatusBadGateway, "message": "Database connection! Save Restaurant"}
	}

	rows, err := db.Query("SELECT * FROM users WHERE email = ?", cq.Email)

	if err != nil {
		return false, gin.H{"status": http.StatusBadRequest, "message": "Database error! Save User", "error": err.Error()}
	}

	if rows.Next() {
		return false, gin.H{"message": "Email address is already registered!", "status": http.StatusBadRequest, "error": err, "data": rows.Next()}
	}

	hashedPass := PasswordHash(fmt.Sprint(cq.Password))

	query_restaurant := "INSERT INTO restaurants (rest_name, summary, address, district, city, country, rest_phone) values (?,?,?,?,?,?)"
	results, err := db.ExecContext(c, query_restaurant, string(cq.RestName), string(cq.Summary), string(cq.Address), string(cq.District), string(cq.City), string(cq.Country), string(cq.RestPhone))
	if results == nil || err != nil {
		return false, gin.H{"status": http.StatusBadRequest, "message": "Insertion error! Save Restaurant", "error": err.Error()}
	}

	resID, err := db.Query("SELECT rest_id FROM restaurants WHERE rest_name = ?", string(cq.RestName))
	var restID int
	for resID.Next() {
		err = resID.Scan(&restID)
	}
	if err != nil {
		return false, gin.H{"status": http.StatusBadRequest, "message": "Selection error! Save Restaurant", "error": err.Error()}
	}
	query_user := "INSERT INTO users (user_name, password, user_phone, email, rest_id, type) values (?,?,?,?,?,?)"
	result_user, err := db.ExecContext(c, query_user, string(cq.UserName), hashedPass, string(cq.UserPhone), string(cq.Email), restID, "Manager")
	if result_user == nil || err != nil {
		return false, gin.H{"status": http.StatusBadRequest, "message": "Insertion error! Save Restaurant", "error": err.Error()}
	}

	credential.Email = string(cq.Email)
	credential.Password = hashedPass
	query := "INSERT INTO credentials (email, password) values (?,?)"
	results, err = db.ExecContext(c, query, string(cq.Email), hashedPass)

	if results == nil || err != nil {
		return false, gin.H{"status": http.StatusBadRequest, "message": "Insertion error! Save Restaurant", "error": err.Error()}
	}

	CloseConnection(db)
	return true, gin.H{"status": http.StatusOK, "message": "Registration Successful!"}
}

func PasswordHash(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("Connection Error!", err.Error())
	}
	return string(hash)
}
