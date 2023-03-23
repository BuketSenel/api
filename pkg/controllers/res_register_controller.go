package controllers

import (
	"database/sql"
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

	if custQuery.RestName == "" || custQuery.Address == "" || custQuery.District == "" || custQuery.City == "" || custQuery.Country == "" || custQuery.Email == "" || custQuery.Password == "" || custQuery.UserPhone == "" {
		fmt.Println("Registration Error!")
		return false, gin.H{"status": http.StatusBadRequest, "message": "Please fill all the fields"}
	}

	result, err := SaveRestaurant(custQuery, c)
	if !result {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return false, gin.H{"status": http.StatusBadRequest, "message": "Registration Error!"}
	}

	return true, gin.H{"status": http.StatusOK, "message": "Registration Error!"}
}

func SaveRestaurant(cq models.CustomQuery, c *gin.Context) (bool, error) {
	credential := models.Credential{}
	db, err := sql.Open("mysql", conf.Name+":"+conf.Password+"@tcp("+conf.Db+":3306)/selfservice")
	if err != nil {
		return false, err
	}
	hashedPass := PasswordHash(fmt.Sprint(cq.Password))

	query_restaurant := "INSERT INTO restaurants (rest_name, address, district, city, country, rest_phone) values (?,?,?,?,?,?,?)"
	results, err := db.ExecContext(c, query_restaurant, string(cq.RestName), string(cq.Address), string(cq.District), string(cq.City), string(cq.Country), string(cq.RestPhone))
	if results == nil || err != nil {
		return false, err
	}

	resID, err := db.Query("SELECT rest_id FROM restaurants WHERE rest_name = ?", string(cq.RestName))
	var restID int
	for resID.Next() {
		err = resID.Scan(&restID)
	}
	if err != nil {
		return false, err
	}
	query_user := "INSERT INTO users (user_name, password, user_phone, email, rest_id, type) values (?,?,?,?,?,?)"
	result_user, err := db.ExecContext(c, query_user, string(cq.UserName), hashedPass, string(cq.UserPhone), string(cq.Email), restID, "Manager")
	if result_user == nil || err != nil {
		return false, err
	}

	credential.Email = string(cq.Email)
	credential.Password = hashedPass
	query := "INSERT INTO credentials (email, password) values (?,?)"
	results, err = db.ExecContext(c, query, string(cq.Email), hashedPass)

	if results == nil || err != nil {
		return false, err
	}

	defer db.Close()
	return true, nil
}

func PasswordHash(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("Connection Error!", err.Error())
	}
	return string(hash)
}
