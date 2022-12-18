package routes

import (
	"github.com/SelfServiceCo/api/pkg/controllers"
	"github.com/gin-gonic/gin"
)

func restaurantRoute(rg *gin.RouterGroup) {
	rg.GET("/:resId", func(c *gin.Context) {
		controllers.GetRestaurant()
	})
}
