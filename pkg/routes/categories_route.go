package routes

import (
	"net/http"
	"strconv"

	"github.com/SelfServiceCo/api/pkg/controllers"
	"github.com/SelfServiceCo/api/pkg/middleware"
	"github.com/gin-gonic/gin"
)

func categoryRoute(rg *gin.RouterGroup) {
	catGroup := rg.Group("/")
	catGroup.Use(middleware.CORSMiddleware())
	catGroup.GET("", func(c *gin.Context) {
		categories, header := controllers.GetCategories()
		if len(categories) == 0 {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, header)
			c.Abort()
		} else {
			c.JSON(http.StatusOK, gin.H{
				"size":   len(categories),
				"offset": "0",
				"limit":  "25",
				"items":  categories,
			})
		}
	})

	catGroup.POST("", func(c *gin.Context) {
		category, header := controllers.CreateCategory(c)
		if category == 0 {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, header)
			c.Abort()
		} else {
			c.JSON(http.StatusOK,
				gin.H{
					"status":  "200",
					"message": "OK",
					"items":   category,
					"offset":  "0",
					"limit":   "25",
				},
			)
		}
	})

	catGroup.GET("/dropdown", func(c *gin.Context) {
		resID := c.Query("resId")
		id, _ := strconv.ParseInt(resID, 10, 64)
		categories, header := controllers.CategoriesForDropdown(id)
		if len(categories) == 0 {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, header)
			c.Abort()
		} else {
			c.JSON(http.StatusOK, gin.H{
				"size":   len(categories),
				"offset": "0",
				"limit":  "25",
				"items":  categories,
			})
		}
	})

	catGroup.POST("/delete", func(c *gin.Context) {
		category, header := controllers.DeleteCategory(c)
		if !category {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, header)
			c.Abort()
		} else {
			c.JSON(http.StatusOK,
				gin.H{
					"status":  "200",
					"message": "OK",
					"items":   category,
					"offset":  "0",
					"limit":   "25",
				},
			)
		}
	})

	catGroup.POST("/edit", func(c *gin.Context) {
		category, header := controllers.EditCategory(c)
		if !category {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, header)
			c.Abort()
		} else {
			c.JSON(http.StatusOK,
				gin.H{
					"status":  "200",
					"message": "OK",
					"items":   category,
					"offset":  "0",
					"limit":   "25",
				},
			)
		}
	})
}
