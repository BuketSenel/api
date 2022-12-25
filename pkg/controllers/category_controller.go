package controllers

import (
	"database/sql"
	"fmt"
	"github.com/SelfServiceCo/api/pkg/models"
)

func GetCategories() []models.Category {
	db, err := sql.Open("mysql", conf.Name+":"+conf.Password+"@tcp("+conf.Db+":3306)/selfservice")

	if err != nil {
		fmt.Println("Err", err.Error())
		return nil
	}
	results, err := db.Query("SELECT * FROM category")
	defer db.Close()

	if err != nil {
		fmt.Println("Err", err.Error())
		return nil
	}

	categories := []models.Category{}
	for results.Next() {
		var cat models.Category
		err = results.Scan(&cat.ID, &cat.ParentCatID, &cat.Name, &cat.Description, &cat.Image)
		if err != nil {
			panic(err.Error())
		}
		categories = append(categories, cat)
	}

	return categories
}
