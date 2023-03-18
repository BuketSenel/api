package routes

import (
	"net/http"

	"github.com/SelfServiceCo/api/pkg/controllers"
	"github.com/gin-gonic/gin"
)

func orderRoute(rg *gin.RouterGroup) {
	orderGroup := rg.Group("/")

	orderGroup.POST("", func(c *gin.Context) {
		order, header := controllers.CreateOrder(c)
		if !order {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, header)
			c.Abort()
		} else {
			c.JSON(http.StatusOK,
				gin.H{
					"status":  http.StatusOK,
					"message": "OK",
					"items":   order,
					"offset":  "0",
					"limit":   "25",
				},
			)
		}
	})

}
