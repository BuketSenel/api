package routes

import (
	"github.com/SelfServiceCo/api/pkg/controllers"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func restaurantRoute(rg *gin.RouterGroup) {
	rg.GET("/:resId", func(c *gin.Context) {
		resID := c.Param("resId")
		id, _ := strconv.ParseInt(resID, 16, 64)
		restaurant := controllers.GetRestaurant(id)
		if restaurant == nil || len(restaurant) == 0 {
			c.AbortWithStatus(http.StatusNotFound)
		} else {
			c.IndentedJSON(http.StatusOK, restaurant)
		}
	})

	rg.GET("/restaurants", func(context *gin.Context) {
		controllers.GetTopRestaurants()
	})
}
