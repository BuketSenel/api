package controllers

import (
	"database/sql"
	"fmt"
	"github.com/SelfServiceCo/api/pkg/models"
)

func ProductsByCategories(cid int64, rid int64) []models.Product {
	db, err := sql.Open("mysql", conf.Name+":"+conf.Password+"@tcp("+conf.Db+":3306)/selfservice")

	if err != nil {
		fmt.Println("Err", err.Error())
		return nil
	}
	results, err := db.Query("SELECT * FROM products WHERE CID = ? AND RID = ?", cid, rid)
	defer db.Close()

	if err != nil {
		fmt.Println("Err", err.Error())
		return nil
	}

	products := []models.Product{}
	for results.Next() {
		var pro models.Product
		err = results.Scan(&pro.ID, &pro.Name, &pro.Description, &pro.CatID, &pro.ResID, &pro.Image, &pro.Price, &pro.Currency, &pro.PrepDurationMin)
		if err != nil {
			panic(err.Error())
		}
		products = append(products, pro)
	}

	return products
}
