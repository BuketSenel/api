package routes

import (
	"github.com/SelfServiceCo/api/pkg/controllers"
	"github.com/gin-gonic/gin"
	"net/http"
)

func restaurantRoute(rg *gin.RouterGroup) {
	rg.Group("/restaurants")
	/*
		rg.GET("/:resId", func(c *gin.Context) {
			resID := c.Param("resId")
			id, _ := strconv.ParseInt(resID, 16, 64)
			restaurant := controllers.GetRestaurant(id)
			if err := c.BindJSON(&restaurant); err != nil || restaurant == nil || len(restaurant) == 0 {
				c.Header("Content-Type", "application/json")
				c.JSON(http.StatusNotFound,
					gin.H{"Error: ": "Invalid startingIndex on search filter!"})
				c.Abort()
			} else {
				c.Header("Content-Type", "application/json")
				c.IndentedJSON(http.StatusOK, restaurant)
			}
		})
	*/
	rg.GET("/", func(c *gin.Context) {
		restaurant := controllers.GetTopRestaurants()
		if restaurant == nil || len(restaurant) == 0 {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound,
				gin.H{"Error: ": "Invalid startingIndex on search filter!"})
			c.Abort()

		} else {
			c.IndentedJSON(http.StatusOK, restaurant)
		}
	})
}
