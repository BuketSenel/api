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
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return false, err
	}

	user.Name = c.PostForm("name")
	user.Phone = c.PostForm("phone")
	user.Email = c.PostForm("email")
	user.Password = c.PostForm("password")

	if c.PostForm("res_id") == "" {
		user.ResID = 0
	} else {
		user.ResID, _ = strconv.ParseInt(c.PostForm("res_id"), 10, 64)
		if c.PostForm("res_id") == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Please fill all the fields"})
			return false, nil
		}
	}
	if user.Name == "" || user.Phone == "" || user.Email == "" || user.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Please fill all the fields"})
		return false, nil
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

	query := "INSERT INTO users (name, password, phone, email) values (?,?,?,?)"
	results, err := db.ExecContext(c, query, string(u.Name), hashedPass, string(u.Phone), string(u.Email), u.CreatedAt)
	if err != nil {
		fmt.Println("Insertion Error!", err.Error())
		return false, err
	}

	credential.Email = string(u.Email)
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
