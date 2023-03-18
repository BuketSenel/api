package controllers

import (
	"database/sql"
	"net/http"

	"github.com/SelfServiceCo/api/pkg/models"
	"github.com/gin-gonic/gin"
)

func getUser(email string) (string, int64, int64, gin.H) {
	db, err := sql.Open("mysql", conf.Name+":"+conf.Password+"@tcp("+conf.Db+":3306)/selfservice")
	if err != nil {
		return "", 0, 0, gin.H{"status": http.StatusBadRequest, "message": "Connection Error!!"}
	}
	result, err := db.Query("SELECT role, UID, resId FROM users WHERE email = ?", email)
	if err != nil {
		return "", 0, 0, gin.H{"status": http.StatusBadRequest, "message": err}
	}
	defer db.Close()
	user := models.User{}
	err = result.Scan(&user.Type, &user.ID, &user.ResID)
	if err != nil {
		return "", 0, 0, gin.H{"status": http.StatusBadRequest, "message": err}
	}
	return user.Type, user.ID, user.ResID, gin.H{"status": http.StatusOK, "message": "Success"}
}
