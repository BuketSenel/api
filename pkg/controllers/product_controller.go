package controllers

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/SelfServiceCo/api/pkg/models"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func ProductsByCategories(cid int64, rid int64) []models.Product {
	db, err := sql.Open("mysql", conf.Name+":"+conf.Password+"@tcp("+conf.Db+":3306)/selfservice")

	if err != nil {
		fmt.Println("Err", err.Error())
		return nil
	}
	results, err := db.Query("SELECT * FROM products WHERE cat_id = ? AND rest_id = ?", cid, rid)
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

func ProductsByRestaurants(rid int64) []models.Product {
	db, err := sql.Open("mysql", conf.Name+":"+conf.Password+"@tcp("+conf.Db+":3306)/selfservice")

	if err != nil {
		fmt.Println("Err", err.Error())
		return nil
	}
	results, err := db.Query("SELECT * FROM products WHERE rest_id = ?", rid)
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

func CreateProduct(c *gin.Context) (bool, gin.H) {
	product := models.Product{}
	err := c.BindJSON(&product)
	if err != nil {
		return false, gin.H{"status": http.StatusBadRequest, "message": "Bind Error! Create Product"}
	}

	db, err := sql.Open("mysql", conf.Name+":"+conf.Password+"@tcp("+conf.Db+":3306)/selfservice")

	if err != nil {
		return false, gin.H{"status": http.StatusBadRequest, "message": "DB Connection Error! Create Product"}
	}

	results, err := db.Query("INSERT INTO products (prod_name, prod_desc, cat_id, rest_id, prod_image, price, currency, prep_dur_minute) VALUES (?,?,?,?,?,?,?,?)", product.Name, product.Description, product.CatID, product.ResID, product.Image, product.Price, product.Currency, product.PrepDurationMin)
	defer db.Close()
	if err != nil {
		return false, gin.H{"status": http.StatusBadRequest, "message": "Insertion Error! Create Product"}
	}
	return true, gin.H{"status": http.StatusOK, "message": "Product created!", "data": results}
}
