package routes

import (
	"github.com/SelfServiceCo/api/pkg/controllers"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func restaurantRoute(rg *gin.RouterGroup) {
	restGroup := rg.Group("/")

	restGroup.GET("/:resId", func(c *gin.Context) {
		resID := c.Param("resId")
		id, _ := strconv.ParseInt(resID, 16, 64)
		restaurant := controllers.GetRestaurant(id)
		if restaurant == nil || len(restaurant) == 0 {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound,
				gin.H{
					"status":    http.StatusNotFound,
					"message: ": "Restaurant not found!",
				})
			c.Abort()
		} else {
			c.JSON(http.StatusOK,
				gin.H{
					"status":  http.StatusOK,
					"message": "OK",
					"size":    len(restaurant),
					"items":   restaurant,
					"offset":  "0",
					"limit":   "25",
				})
		}
	})

	restGroup.GET("/:resId/categories", func(c *gin.Context) {
		resID := c.Param("resId")
		id, _ := strconv.ParseInt(resID, 16, 64)
		categories := controllers.CategoriesByRestaurant(id)
		if categories == nil || len(categories) == 0 {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound,
				gin.H{
					"status":    http.StatusNotFound,
					"message: ": "No categories found!",
				})
			c.Abort()
		} else {
			c.JSON(http.StatusOK,
				gin.H{
					"status":  http.StatusOK,
					"message": "OK",
					"size":    len(categories),
					"items":   categories,
					"offset":  "0",
					"limit":   "25",
				})
		}
	})

	restGroup.GET("", func(c *gin.Context) {
		restaurant := controllers.GetTopRestaurants()
		if restaurant == nil || len(restaurant) == 0 {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound,
				gin.H{
					"status":    http.StatusNotFound,
					"message: ": "No restaurants retrieved!",
				})
			c.Abort()
		} else {
			c.JSON(http.StatusOK, gin.H{
				"status":  "200",
				"message": "OK",
				"size":    len(restaurant),
				"items":   restaurant,
				"offset":  "0",
				"limit":   "25",
			})
		}
	})
}
