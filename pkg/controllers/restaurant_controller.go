package controllers

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/SelfServiceCo/api/pkg/drivers"
	"github.com/SelfServiceCo/api/pkg/models"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

// bu nedir?
var selfdb = "selfservicedb"
var conf = drivers.MysqlConfigLoad()

func GetRestaurant(id int64) ([]models.Restaurant, gin.H) {
	db, err := sql.Open("mysql", conf.Name+":"+conf.Password+"@tcp("+conf.Db+":3306)/selfservice?parseTime=true")

	if err != nil {
		return nil, gin.H{"status": http.StatusBadRequest, "message": "DB Connection Error! Get Restaurant"}
	}

	results, err := db.Query("SELECT * FROM restaurants WHERE rest_id = ?", id)

	if err != nil {
		return nil, gin.H{"status": http.StatusBadRequest, "message": "Selection Error! Get Restaurant"}
	}

	restaurant := []models.Restaurant{}
	for results.Next() {
		rest := models.Restaurant{}
		err = results.Scan(&rest.ID, &rest.Name, &rest.Summary, &rest.Logo, &rest.Address, &rest.District, &rest.City, &rest.Country, &rest.Phone, &rest.Tags, &rest.CreatedAt, &rest.UpdatedAt)
		if err != nil {
			return nil, gin.H{"status": http.StatusBadRequest, "message": "Scan Error! Get Restaurant", "data": results, "Error": err.Error()}
		}
		restaurant = append(restaurant, rest)
	}
	defer db.Close()

	return restaurant, gin.H{"status": http.StatusOK, "message": restaurant}
}

func GetTopRestaurants() ([]models.Restaurant, gin.H) {
	db, err := sql.Open("mysql", conf.Name+":"+conf.Password+"@tcp("+conf.Db+":3306)/selfservice?parseTime=true")

	if err != nil {
		return nil, gin.H{"status": http.StatusBadRequest, "message": "DB Connection Error! Get Top Restaurant"}
	}

	results, err := db.Query("SELECT * FROM restaurants")

	if err != nil {
		return nil, gin.H{"status": http.StatusBadRequest, "message": "Selection Error! Get Top Restaurant"}
	}

	restaurants := []models.Restaurant{}
	for results.Next() {
		var rest models.Restaurant

		err = results.Scan(&rest.ID, &rest.Name, &rest.Summary, &rest.Logo, &rest.Address, &rest.District, &rest.City, &rest.Country, &rest.Phone, &rest.Tags, &rest.CreatedAt, &rest.UpdatedAt)
		if err != nil {
			return nil, gin.H{"status": http.StatusBadRequest, "message": "Scan Error! Get Top Restaurant", "data": results, "Error": err.Error()}
		}
		restaurants = append(restaurants, rest)
	}
	defer db.Close()

	return restaurants, gin.H{"status": http.StatusOK, "message": restaurants}
}

func GetRestaurantStaff(rid int64) []models.User {
	db, err := sql.Open("mysql", conf.Name+":"+conf.Password+"@tcp("+conf.Db+":3306)/selfservice")

	if err != nil {
		fmt.Println("Err", err.Error())
		return nil
	}

	results, err := db.Query("SELECT * FROM users WHERE rest_id = ?", rid)
	defer db.Close()

	if err != nil {
		fmt.Println("Err", err.Error())
		return nil
	}

	staff := []models.User{}
	for results.Next() {
		var user models.User
		err = results.Scan(&user.ID, &user.Name, &user.Email, &user.ResID, &user.Type, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			panic(err.Error())
		}
		staff = append(staff, user)
	}

	return staff
}

func AddStaff(c *gin.Context) (bool, gin.H) {
	user := models.User{}

	if err := c.BindJSON(&user); err != nil {
		c.AbortWithError(401, err)
	}

	if user.Name == "" || user.Phone == "" || user.Email == "" || user.Password == "" || user.ResID == 0 || user.Type == "" {
		return false, gin.H{"error": "Please fill all the fields"}
	}

	result, err := SaveUser(user, c)

	if !result || err != nil {
		return result, err
	}

	return true, gin.H{"message": "Staff registered successfully"}
}

func DeleteStaff(c *gin.Context) (bool, gin.H) {
	user := models.User{}

	if err := c.BindJSON(&user); err != nil {
		c.AbortWithError(401, err)
	}

	if user.ID == 0 {
		return false, gin.H{"error": "Please fill all the fields"}
	}

	//result, err := DeleteUser(user, c)

	//if !result || err != nil {
	//	return result, err
	//}

	return true, gin.H{"message": "Staff deleted successfully"}
}
