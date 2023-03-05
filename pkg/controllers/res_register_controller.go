package controllers

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/SelfServiceCo/api/pkg/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func RestaurantRegister(c *gin.Context) bool {
	restaurant := models.Restaurant{}

	if err := c.ShouldBindJSON(&restaurant); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return false
	}

	restaurant.Name = c.PostForm("name")
	restaurant.Address = c.PostForm("address")
	restaurant.District = c.PostForm("district")
	restaurant.City = c.PostForm("city")
	restaurant.Country = c.PostForm("country")
	restaurant.Email = c.PostForm("email")
	restaurant.Password = c.PostForm("password")
	restaurant.Phone = c.PostForm("phone")
	if !SaveRestaurant(restaurant, c) {
		fmt.Println("Registration Error!")
		return false
	}
	c.JSON(http.StatusOK, gin.H{"message": "Restaurant registered!"})
	return true
}

func SaveRestaurant(r models.Restaurant, c *gin.Context) bool {
	credential := models.Credential{}
	db, err := sql.Open("mysql", conf.Name+":"+conf.Password+"@tcp("+conf.Db+":3306)/selfservice")
	if err != nil {
		fmt.Println("Connection Error!", err.Error())
		return false
	}
	hashedPass := PasswordHash(fmt.Sprint(r.Password))

	query := "INSERT INTO restaurant (name, password, address, district, city, country, phone, email, created_at) values (?,?,?,?,?,?,?,?,?)"
	results, err := db.ExecContext(c, query, string(r.Name), hashedPass, string(r.Address), string(r.District), string(r.City), string(r.Country), string(r.Phone), string(r.Email), time.Now())
	if err != nil {
		fmt.Println("Insertion Error!", err.Error())
		return false
	}

	credential.Email = string(r.Email)
	credential.Password = hashedPass
	query = "INSERT INTO credentials (email, password) values (?,?)"
	results, err = db.ExecContext(c, query, string(r.Email), hashedPass)
	if err != nil {
		fmt.Println("Insertion Error!", err.Error())
		return false
	}

	defer db.Close()
	fmt.Println("Results: ", results)
	return true
}

func PasswordHash(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("Connection Error!", err.Error())
	}
	return string(hash)
}
