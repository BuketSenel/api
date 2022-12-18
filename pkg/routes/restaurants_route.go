package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func restaurantRoute(rg *gin.RouterGroup) {
	rg.GET("/:resId", func(c *gin.Context) {
		resId := c.Params.ByName("resId")
		c.JSON(http.StatusOK, gin.H{"resId": resId})
	})
}
