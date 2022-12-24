package routes

import (
	"fmt"
	"github.com/SelfServiceCo/api/pkg/controllers"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func restaurantRoute(rg *gin.RouterGroup) {
	rg.Group("/restaurants")

	rg.GET("/:resId", func(c *gin.Context) {
		resID := c.Param("resId")
		id, _ := strconv.ParseInt(resID, 16, 64)
		restaurant := controllers.GetRestaurant(id)
		if err := c.BindJSON(&restaurant); err != nil || restaurant == nil || len(restaurant) == 0 {
			fmt.Println(err)
			fmt.Println(restaurant)
			c.AbortWithStatus(http.StatusNotFound)
		} else {
			c.AbortWithStatus(http.StatusNotFound)
		}
	})

	rg.GET("/", func(c *gin.Context) {
		restaurant := controllers.GetTopRestaurants()
		if restaurant == nil || len(restaurant) == 0 {
			c.AbortWithStatus(http.StatusNotFound)
		} else {
			c.IndentedJSON(http.StatusOK, restaurant)
		}
	})
}
