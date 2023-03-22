package routes

import (
	"net/http"

	"github.com/SelfServiceCo/api/pkg/controllers"
	"github.com/gin-gonic/gin"
)

func orderRoute(rg *gin.RouterGroup) {
	orderGroup := rg.Group("/")

	orderGroup.POST("", func(c *gin.Context) {
		result, header := controllers.CreateOrder(c)
		if !result {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, header)
			c.Abort()
		} else {
			c.JSON(http.StatusOK,
				gin.H{
					"status":  http.StatusOK,
					"message": "OK",
					"items":   result,
					"offset":  "0",
					"limit":   "25",
				},
			)
		}
	})

}
