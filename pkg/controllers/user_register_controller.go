package controllers

import (
	"database/sql"
	"fmt"
	"github.com/SelfServiceCo/api/pkg/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func UserRegister(c *gin.Context) bool {
	user := models.User{}
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return false
	}

	user.Name = c.PostForm("name")
	user.Phone = c.PostForm("phone")
	user.Email = c.PostForm("email")
	user.CreatedAt = time.Now()

	if !SaveUser(user, c) {
		fmt.Println("Registration Error!")
		return false
	}
	c.JSON(http.StatusOK, gin.H{"message": "validated!"})
	return true
}

func SaveUser(u models.User, c *gin.Context) bool {
	credential := models.Credential{}
	db, err := sql.Open("mysql", conf.Name+":"+conf.Password+"@tcp("+conf.Db+":3306)/selfservice")
	if err != nil {
		fmt.Println("Connection Error!", err.Error())
		return false
	}
	hashedPass := PasswordHash(fmt.Sprint(u.Password))

	query := "INSERT INTO restaurant (name, password, phone, email, created_at) values (?,?,?,?,?)"
	results, err := db.ExecContext(c, query, string(u.Name), hashedPass, string(u.Phone), string(u.Email), u.CreatedAt)
	if err != nil {
		fmt.Println("Insertion Error!", err.Error())
		return false
	}

	credential.Email = string(u.Email)
	credential.Password = hashedPass
	query = "INSERT INTO credentials (email, password) values (?,?)"
	results, err = db.ExecContext(c, query, string(u.Email), hashedPass)
	if err != nil {
		fmt.Println("Insertion Error!", err.Error())
		return false
	}

	defer db.Close()
	fmt.Println("Results: ", results)
	return true
}
