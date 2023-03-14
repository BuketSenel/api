package controllers

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/SelfServiceCo/api/pkg/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Login(c *gin.Context) (string, gin.H) {
	cred := models.Credential{}
	var hashed string

	if err := c.BindJSON(&cred); err != nil {
		c.AbortWithError(401, err)
	}

	if cred.Email == "" || cred.Password == "" {
		return "", gin.H{"status": http.StatusBadRequest, "message": "Please fill all the fields"}
	}

	db, err := sql.Open("mysql", conf.Name+":"+conf.Password+"@tcp("+conf.Db+":3306)/selfservice")
	if err != nil {
		return "", gin.H{"status": http.StatusBadRequest, "message": "Connection Error!!"}
	}

	err = db.QueryRow("SELECT password from credentials WHERE email = ?", cred.Email).Scan(&hashed)
	if err != nil {
		fmt.Println("Selection Error!", err.Error())
		return "", gin.H{"status": http.StatusBadRequest, "message": err}
	}

	err = VerifyPassword(cred.Password, hashed)
	if err != nil {
		return "", gin.H{"status": http.StatusBadRequest, "message": "Verification Error!"}
	}

	token, err := CreateJWTToken()

	if err != nil {
		return "", gin.H{"status": http.StatusBadRequest, "message": "Token Creation Error!"}
	}

	results, err := db.Query("INSERT INTO credentials (token) VALUES (?) WHERE email = ?", token, cred.Email)
	if err != nil {
		fmt.Println("Selection Error!", err.Error())
		return "", gin.H{"status": http.StatusBadRequest, "message": err, "data": results}
	}

	return token, gin.H{"status": http.StatusOK, "message": "Login Successful!"}
}

func VerifyPassword(password string, hashed string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
}
