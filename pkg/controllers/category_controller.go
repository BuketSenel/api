package controllers

import (
	"database/sql"
	"fmt"

	"github.com/SelfServiceCo/api/pkg/models"
	_ "github.com/go-sql-driver/mysql"
)

func GetCategories() []models.Category {
	db, err := sql.Open("mysql", conf.Name+":"+conf.Password+"@tcp("+conf.Db+":3306)/selfservice")

	if err != nil {
		fmt.Println("Err", err.Error())
		return nil
	}
	results, err := db.Query("SELECT * FROM categories WHERE rest_id = ?", 0)
	defer db.Close()

	if err != nil {
		fmt.Println("Err", err.Error())
		return nil
	}

	categories := []models.Category{}
	for results.Next() {
		var cat models.Category
		err = results.Scan(&cat.ID, &cat.RID, &cat.Name, &cat.Description, &cat.Image, &cat.ParentCatID)
		if err != nil {
			panic(err.Error())
		}
		categories = append(categories, cat)
	}

	return categories
}

func CategoriesByRestaurant(rid int64) []models.Category {
	db, err := sql.Open("mysql", conf.Name+":"+conf.Password+"@tcp("+conf.Db+":3306)/selfservice")

	if err != nil {
		fmt.Println("Err", err.Error())
		return nil
	}
	results, err := db.Query("SELECT * FROM categories WHERE rest_id = ?", rid)
	defer db.Close()

	if err != nil {
		fmt.Println("Err", err.Error())
		return nil
	}

	categories := []models.Category{}
	for results.Next() {
		var cat models.Category
		err = results.Scan(&cat.ID, &cat.RID, &cat.Name, &cat.Description, &cat.Image, &cat.ParentCatID)
		if err != nil {
			panic(err.Error())
		}
		categories = append(categories, cat)
	}

	return categories
}
