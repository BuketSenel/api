package controllers

import (
	"net/http"

	"github.com/SelfServiceCo/api/pkg/models"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func ProductsByCategories(cid int64, rid int64) ([]models.Product, gin.H) {
	db := CreateConnection()

	if db == nil {
		return nil, gin.H{"status": http.StatusBadRequest, "message": "DB Connection Error! Get Products By Categories"}
	}
	results, err := db.Query("SELECT * FROM products WHERE cat_id = ? AND rest_id = ?", cid, rid)
	CloseConnection(db)

	if err != nil {
		return nil, gin.H{"status": http.StatusBadRequest, "message": "Query Error! Get Products By Categories"}
	}

	products := []models.Product{}
	for results.Next() {
		var pro models.Product
		err = results.Scan(&pro.ID, &pro.Name, &pro.Description, &pro.CatID, &pro.ResID, &pro.Image, &pro.Price, &pro.Currency, &pro.PrepDurationMin)
		if err != nil {
			return nil, gin.H{"status": http.StatusBadRequest, "message": "Scan Error! Get Products By Categories"}
		}
		products = append(products, pro)
	}

	return products, gin.H{"status": http.StatusOK, "message": "OK"}
}

func ProductsByRestaurants(rid int64) ([]models.Product, gin.H) {
	db := CreateConnection()

	if db == nil {
		return nil, gin.H{"status": http.StatusBadRequest, "message": "DB Connection Error! Get Products By Restaurants"}
	}
	results, err := db.Query("SELECT * FROM products WHERE rest_id = ?", rid)
	CloseConnection(db)

	if err != nil {
		return nil, gin.H{"status": http.StatusBadRequest, "message": "Query Error! Get Products By Restaurants"}
	}

	products := []models.Product{}
	for results.Next() {
		var pro models.Product
		err = results.Scan(&pro.ID, &pro.Name, &pro.Description, &pro.CatID, &pro.ResID, &pro.Image, &pro.Price, &pro.Currency, &pro.PrepDurationMin)
		if err != nil {
			return nil, gin.H{"status": http.StatusBadRequest, "message": "Scan Error! Get Products By Restaurants"}
		}
		products = append(products, pro)
	}

	return products, gin.H{"status": http.StatusOK, "message": "OK"}
}

func CreateProduct(c *gin.Context) (bool, gin.H) {
	product := models.Product{}
	err := c.BindJSON(&product)
	if err != nil {
		return false, gin.H{"status": http.StatusBadRequest, "message": "Bind Error! Create Product"}
	}

	db := CreateConnection()

	if db == nil {
		return false, gin.H{"status": http.StatusBadRequest, "message": "DB Connection Error! Create Product"}
	}

	results, err := db.Query("INSERT INTO products (prod_name, prod_desc, cat_id, rest_id, prod_image, price, currency, prep_dur_minute) VALUES (?,?,?,?,?,?,?,?)", product.Name, product.Description, product.CatID, product.ResID, product.Image, product.Price, product.Currency, product.PrepDurationMin)
	CloseConnection(db)
	if err != nil {
		return false, gin.H{"status": http.StatusBadRequest, "message": "Insertion Error! Create Product"}
	}
	return true, gin.H{"status": http.StatusOK, "message": "Product created!", "data": results}
}

func EditProduct(c *gin.Context) (bool, gin.H) {
	product := models.Product{}
	err := c.BindJSON(&product)
	if err != nil {
		return false, gin.H{"status": http.StatusBadRequest, "message": "Bind Error! Edit Product"}
	}

	db := CreateConnection()

	if db == nil {
		return false, gin.H{"status": http.StatusBadRequest, "message": "DB Connection Error! Edit Product"}
	}

	results, err := db.Query("UPDATE products SET prod_name = ?, prod_desc = ?, price = ?, currency = ?, prep_dur_minute = ? WHERE prod_id = ?", product.Name, product.Description, product.Price, product.Currency, product.PrepDurationMin, product.ID)
	CloseConnection(db)
	if err != nil {
		return false, gin.H{"status": http.StatusBadRequest, "message": "Update Error! Edit Product"}
	}
	return true, gin.H{"status": http.StatusOK, "message": "Product updated!", "data": results}
}

func DeleteProduct(c *gin.Context) (bool, gin.H) {
	product := models.Product{}
	err := c.BindJSON(&product)
	if err != nil {
		return false, gin.H{"status": http.StatusBadRequest, "message": "Bind Error! Delete Product"}
	}

	db := CreateConnection()

	if db == nil {
		return false, gin.H{"status": http.StatusBadRequest, "message": "DB Connection Error! Delete Product"}
	}

	results, err := db.Query("DELETE FROM products WHERE prod_id = ?", product.ID)
	CloseConnection(db)
	if err != nil {
		return false, gin.H{"status": http.StatusBadRequest, "message": "Delete Error! Delete Product"}
	}
	return true, gin.H{"status": http.StatusOK, "message": "Product deleted!", "data": results}
}
