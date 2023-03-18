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
	restaurant := models.Restaurant{}
	if err := c.BindJSON(&restaurant); err != nil {
		c.AbortWithError(401, err)
	}

	if restaurant.Name == "" || restaurant.Address == "" || restaurant.District == "" || restaurant.City == "" || restaurant.Country == "" || restaurant.Email == "" || restaurant.Password == "" || restaurant.Phone == "" {
		fmt.Println("Registration Error!")
		return false, gin.H{"status": http.StatusBadRequest, "message": "Please fill all the fields"}
	}

	result, err := SaveRestaurant(restaurant, c)
	if !result {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return false, gin.H{"status": http.StatusBadRequest, "message": "Registration Error!"}
	}

	return true, gin.H{"status": http.StatusOK, "message": "Registration Error!"}
}

func SaveRestaurant(r models.Restaurant, c *gin.Context) (bool, error) {
	credential := models.Credential{}
	db, err := sql.Open("mysql", conf.Name+":"+conf.Password+"@tcp("+conf.Db+":3306)/selfservice")
	if err != nil {
		return false, err
	}
	hashedPass := PasswordHash(fmt.Sprint(r.Password))
	query_user := "INSERT INTO users (name, password, phone, email, resID, type) values (?,?,?,?,?,?)"
	result_user, err := db.ExecContext(c, query_user, string(r.Name), hashedPass, string(r.Address), string(r.District), string(r.City), string(r.Country), string(r.Phone), string(r.Email))
	if result_user == nil || err != nil {
		return false, err
	}
	query_restaurant := "INSERT INTO restaurants (name, address, district, city, country, phone, email) values (?,?,?,?,?,?,?,?)"
	results, err := db.ExecContext(c, query_restaurant, string(r.Name), string(r.Address), string(r.District), string(r.City), string(r.Country), string(r.Phone), string(r.Email))
	if results == nil || err != nil {
		return false, err
	}

	credential.Email = string(r.Email)
	credential.Password = hashedPass
	query = "INSERT INTO credentials (email, password) values (?,?)"
	results, err = db.ExecContext(c, query, string(r.Email), hashedPass)

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
