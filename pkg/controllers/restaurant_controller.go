package controllers

import (
	"database/sql"
	"fmt"
	"github.com/SelfServiceCo/api/pkg/models"
	_ "github.com/go-sql-driver/mysql"
)

func GetRestaurant(id int64) []models.Restaurant {
	db, err := sql.Open("mysql", "")

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
			&rest.City, &rest.Country, &rest.Phone, &rest.Tags)
		if err != nil {
			panic(err.Error())
		}
		restaurant = append(restaurant, rest)
	}

	return restaurant
}

func GetTopRestaurants() []models.Restaurant {
	db, err := sql.Open("mysql", "")

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
			&rest.City, &rest.Country, &rest.Phone, &rest.Tags)
		if err != nil {
			panic(err.Error())
		}
		restaurants = append(restaurants, rest)
	}

	return restaurants
}
