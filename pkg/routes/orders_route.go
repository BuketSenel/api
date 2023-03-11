package routes

import (
	"net/http"
	"strconv"

	"github.com/SelfServiceCo/api/pkg/controllers"
	"github.com/gin-gonic/gin"
)

func OrderRoute(rg *gin.RouterGroup) {
	orderGroup := rg.Group("/")

	orderGroup.GET("/:orderId", func(c *gin.Context) {
		orderID := c.Param("orderId")
		id, _ := strconv.ParseInt(orderID, 16, 64)
		order := controllers.GetOrder(id)
		if order == nil || len(order) == 0 {
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

	orderGroup.POST("", func(c *gin.Context) {
		order := controllers.CreateOrder(c)
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
