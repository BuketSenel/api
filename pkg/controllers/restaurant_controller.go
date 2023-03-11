package controllers

import (
	"database/sql"
	"fmt"

	"github.com/SelfServiceCo/api/pkg/drivers"
	"github.com/SelfServiceCo/api/pkg/models"
	_ "github.com/go-sql-driver/mysql"
)

var selfdb = "selfservicedb"
var conf = drivers.MysqlConfigLoad()

func GetRestaurant(id int64) []models.Restaurant {
	db, err := sql.Open("mysql", conf.Name+":"+conf.Password+"@tcp("+conf.Db+":3306)/selfservice")

	if err != nil {
		fmt.Println("Err", err.Error())
		return nil
	}

	results, err := db.Query("SELECT * FROM restaurant WHERE ID = ?", id)
	defer db.Close()

	if err != nil {
		fmt.Println("Err", err.Error())
		return nil
	}

	restaurant := []models.Restaurant{}
	for results.Next() {
		var rest models.Restaurant

		err = results.Scan(&rest.ID, &rest.Name, &rest.Summary, &rest.Logo, &rest.Address, &rest.District,
			&rest.City, &rest.Country, &rest.Phone, &rest.Tags, &rest.CreatedAt, &rest.UpdatedAt)
		if err != nil {
			panic(err.Error())
		}
		restaurant = append(restaurant, rest)
	}

	return restaurant
}

func GetTopRestaurants() []models.Restaurant {
	db, err := sql.Open("mysql", conf.Name+":"+conf.Password+"@tcp("+conf.Db+":3306)/selfservice")

	if err != nil {
		fmt.Println("Err", err.Error())
		return nil
	}

	results, err := db.Query("SELECT * FROM restaurant")
	defer db.Close()

	if err != nil {
		fmt.Println("Err", err.Error())
		return nil
	}

	restaurants := []models.Restaurant{}
	for results.Next() {
		var rest models.Restaurant

		err = results.Scan(&rest.ID, &rest.Name, &rest.Summary, &rest.Logo, &rest.Address, &rest.District,
			&rest.City, &rest.Country, &rest.Phone, &rest.Tags, &rest.CreatedAt, &rest.UpdatedAt)
		if err != nil {
			panic(err.Error())
		}
		restaurants = append(restaurants, rest)
	}

	return restaurants
}

func GetRestaurantStaff(rid int64) []models.User {
	db, err := sql.Open("mysql", conf.Name+":"+conf.Password+"@tcp("+conf.Db+":3306)/selfservice")

	if err != nil {
		fmt.Println("Err", err.Error())
		return nil
	}

	results, err := db.Query("SELECT * FROM users WHERE RID = ?", rid)
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
