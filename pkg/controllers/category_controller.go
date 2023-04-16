package controllers

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/SelfServiceCo/api/pkg/models"
	"github.com/gin-gonic/gin"
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

func CategoriesByRestaurant(rid int64) ([]models.Category, gin.H) {
	db, err := sql.Open("mysql", conf.Name+":"+conf.Password+"@tcp("+conf.Db+":3306)/selfservice")

	if err != nil {
		fmt.Println("Err", err.Error())
		return nil, gin.H{"status": http.StatusBadRequest, "message": "DB Connection Error! Get Categories By Restaurant"}
	}
	results, err := db.Query("SELECT * FROM categories WHERE rest_id = ?", rid)
	defer db.Close()

	if err != nil {
		return nil, gin.H{"status": http.StatusBadRequest, "message": "Query Error! Get Categories By Restaurant"}
	}

	categories := []models.Category{}
	for results.Next() {
		var cat models.Category
		err = results.Scan(&cat.ID, &cat.RID, &cat.Name, &cat.Description, &cat.Image, &cat.ParentCatID)
		if err != nil {
			return nil, gin.H{"status": http.StatusBadRequest, "message": "Scan Error! Get Categories By Restaurant"}
		}
		categories = append(categories, cat)
	}

	return categories, gin.H{"status": http.StatusOK, "message": "OK"}
}

func CreateCategory(c *gin.Context) (int64, gin.H) {
	category := models.Category{}
	err := c.BindJSON(&category)
	if err != nil {
		return 0, gin.H{"status": http.StatusBadRequest, "message": "Bind Error! Create Product Category"}
	}

	db, err := sql.Open("mysql", conf.Name+":"+conf.Password+"@tcp("+conf.Db+":3306)/selfservice")

	if err != nil {
		return 0, gin.H{"status": http.StatusBadRequest, "message": "DB Connection Error! Create Product Category"}
	}

	results, err := db.Query("INSERT INTO categories (cat_name, rest_id, cat_desc, cat_image, parent_cat_id) VALUES (?,?,?,?,?)", category.Name, category.RID, category.Description, category.Image, category.ParentCatID)

	if err != nil {
		return 0, gin.H{"status": http.StatusBadRequest, "message": "Insertion Error! Create Product Category"}
	}

	result, err := db.Query("SELECT max(cat_id) from categories")
	if err != nil {
		return 0, gin.H{"status": http.StatusBadRequest, "message": "Selection Error! Create Product Category"}
	}
	result.Next()
	err = result.Scan(&category.ID)

	if err != nil {
		return 0, gin.H{"status": http.StatusBadRequest, "message": "Scan Error! Create Product Category"}
	}
	defer db.Close()
	return category.ID, gin.H{"status": http.StatusOK, "message": "Product category created!", "data": results, "result": result, "cat_id": category.ID}
}

func CategoriesForDropdown(rid int64) []models.Category {
	db, err := sql.Open("mysql", conf.Name+":"+conf.Password+"@tcp("+conf.Db+":3306)/selfservice")

	if err != nil {
		fmt.Println("Err", err.Error())
		return nil
	}
	results, err := db.Query("SELECT DISTINCT cat_id, cat_name FROM categories WHERE rest_id = ? or rest_id = 0", rid)
	defer db.Close()

	if err != nil {
		fmt.Println("Err", err.Error())
		return nil
	}

	categories := []models.Category{}
	for results.Next() {
		var cat models.Category
		err = results.Scan(&cat.ID, &cat.Name)
		if err != nil {
			panic(err.Error())
		}
		categories = append(categories, cat)
	}

	return categories
}
