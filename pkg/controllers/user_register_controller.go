package controllers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/SelfServiceCo/api/pkg/models"
	"github.com/gin-gonic/gin"
)

func UserRegister(c *gin.Context) (bool, error) {
	user := models.User{}

	user.Name = c.PostForm("name")
	user.Phone = c.PostForm("phone")
	user.Email = c.PostForm("email")
	user.Password = c.PostForm("password")
	user.Type = c.PostForm("type")

	if user.Name == "" || user.Phone == "" || user.Email == "" || user.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Please fill all the fields"})
		return false, nil
	}

	if c.PostForm("res_id") == "" {
		user.ResID = 0
	} else {
		user.ResID, _ = strconv.ParseInt(c.PostForm("res_id"), 10, 64)
	}

	result, err := SaveUser(user, c)

	if !result || err != nil {
		fmt.Println("Registration Error!")
		return false, err
	}
	c.JSON(http.StatusOK, gin.H{"message": "User registered!"})
	return true, nil
}

func SaveUser(u models.User, c *gin.Context) (bool, error) {
	credential := models.Credential{}
	db, err := sql.Open("mysql", conf.Name+":"+conf.Password+"@tcp("+conf.Db+":3306)/selfservice")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return false, err
	}
	hashedPass := PasswordHash(fmt.Sprint(u.Password))

	query := "INSERT INTO users (name, password, phone, email, resID, type) values (?,?,?,? ?,?)"
	results, err := db.ExecContext(c, query, u.Name, hashedPass, u.Phone, u.Email, u.ResID, u.Type)
	if err != nil {
		fmt.Println("Insertion Error!", err.Error())
		return false, err
	}

	credential.Email = u.Email
	credential.Password = hashedPass
	query = "INSERT INTO credentials (email, password) values (?,?)"
	results, err = db.ExecContext(c, query, string(u.Email), hashedPass)
	if err != nil {
		fmt.Println("Insertion Error!", err.Error())
		return false, err
	}

	defer db.Close()
	fmt.Println("Results: ", results)
	return true, nil
}
