package controllers

import (
	"database/sql"
	"fmt"

	"github.com/SelfServiceCo/api/pkg/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Login(c *gin.Context) bool {
	cred := models.Credential{}
	var hashed string
	cred.Email = c.PostForm("email")
	cred.Password = c.PostForm("password")

	if cred.Email == "" || cred.Password == "" {
		fmt.Println("Empty Fields!")
		return false
	}

	db, err := sql.Open("mysql", conf.Name+":"+conf.Password+"@tcp("+conf.Db+":3306)/selfservice")
	if err != nil {
		fmt.Println("Connection Error!", err.Error())
		return false
	}

	err = db.QueryRow("SELECT password from credentials email = ?", cred.Email).Scan(&hashed)
	if err != nil {
		fmt.Println("Selection Error!", err.Error())
		return false
	}

	VerifyPassword(cred.Password, hashed)
	if err != nil {
		fmt.Println("Verification Error!", err.Error())
		return false
	}
	return true
}

func VerifyPassword(password string, hashed string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
}
