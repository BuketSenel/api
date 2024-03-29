package controllers

import (
	"fmt"
	"net/http"

	"github.com/SelfServiceCo/api/pkg/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Login(c *gin.Context) (string, gin.H) {
	cred := models.Credential{}
	c.BindJSON(&cred)
	if cred.Email == "" || cred.Password == "" {
		return "", gin.H{"status": http.StatusBadRequest, "message": "Please fill all the fields"}
	}

	db := CreateConnection()
	if db == nil {
		return "", gin.H{"status": http.StatusBadRequest, "message": "Connection Error!!"}
	}

	var hashed string
	err := db.QueryRow("SELECT password from credentials WHERE email = ?", cred.Email).Scan(&hashed)
	if err != nil {
		fmt.Println("Selection Error!", err.Error())
		return "", gin.H{"status": http.StatusBadRequest, "message": err}
	}

	err = VerifyPassword(cred.Password, hashed)
	if err != nil {
		return "", gin.H{"status": http.StatusBadRequest, "message": "Verification Error!"}
	}

	token, header := CreateJWTToken(cred.Email)
	if header["status"] != 200 {
		return "", header
	}
	results, err := db.Exec("UPDATE credentials SET token = ? WHERE email = ?", token, cred.Email)

	if err != nil {
		return "", gin.H{"status": http.StatusBadRequest, "message": err, "data": results}
	}
	CloseConnection(db)
	return token, gin.H{"status": http.StatusOK, "message": "Login Successful!"}
}

func VerifyPassword(password string, hashed string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
}
