package routes

import (
	"net/http"

	"github.com/SelfServiceCo/api/pkg/controllers"
	"github.com/SelfServiceCo/api/pkg/models"
	"github.com/gin-gonic/gin"
)

func OrderRoute(rg *gin.RouterGroup) {
	orderGroup := rg.Group("/")

	orderGroup.POST("", func(c *gin.Context) {
		var orderRequest models.Order
		err := c.BindJSON(&orderRequest)
		if err != nil {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound,
				gin.H{
					"status":    http.StatusNotFound,
					"message: ": "Order not found!",
				})
			c.Abort()
		}
		order := controllers.CreateOrder(orderRequest)
		if order == nil {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound,
				gin.H{
					"status":    http.StatusNotFound,
					"message: ": "Order not found!",
				})
			c.Abort()
		} else {
			c.JSON(http.StatusOK,
				gin.H{
					"status":  http.StatusOK,
					"message": "OK",
					"size":    len(order),
					"items":   order,
					"offset":  "0",
					"limit":   "25",
				},
			)
		}
	})

}
