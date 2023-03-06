package controllers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/SelfServiceCo/api/pkg/models"
	"github.com/gin-gonic/gin"
)

func UserRegister(c *gin.Context) (bool, gin.H) {
	user := models.User{}

	user.Name = c.PostForm("name")
	user.Phone = c.PostForm("phone")
	user.Email = c.PostForm("email")
	user.Password = c.PostForm("password")
	user.Type = c.PostForm("type")

	if user.Name == "" || user.Phone == "" || user.Email == "" || user.Password == "" {
		return false, gin.H{"error": "Please fill all the fields"}
	}

	if c.PostForm("res_id") == "" {
		user.ResID = 0
	} else {
		user.ResID, _ = strconv.ParseInt(c.PostForm("res_id"), 10, 64)
	}

	result, err := SaveUser(user, c)

	if !result || err != nil {
		return false, gin.H{"error": "Registration Error"}
	}

	return true, gin.H{"message": "User registered!"}
}

func SaveUser(u models.User, c *gin.Context) (bool, gin.H) {
	credential := models.Credential{}
	db, err := sql.Open("mysql", conf.Name+":"+conf.Password+"@tcp("+conf.Db+":3306)/selfservice")
	if err != nil {
		message := gin.H{"status": http.StatusBadGateway, "message": err.Error()}
		return false, message
	}
	hashedPass := PasswordHash(fmt.Sprint(u.Password))

	query := "INSERT INTO users (name, password, phone, email, resID, type) values (?,?,?,? ?,?)"
	results, err := db.ExecContext(c, query, u.Name, hashedPass, u.Phone, u.Email, u.ResID, u.Type)
	if err != nil {
		message := gin.H{"status": http.StatusBadRequest, "message": err.Error()}
		return false, message
	}

	credential.Email = u.Email
	credential.Password = hashedPass
	query = "INSERT INTO credentials (email, password) values (?,?)"
	results, err = db.ExecContext(c, query, string(u.Email), hashedPass)
	if err != nil {
		message := gin.H{"status": http.StatusBadRequest, "message": err.Error()}
		return false, message
	}

	defer db.Close()
	fmt.Println("Results: ", results)
	return true, nil
}
