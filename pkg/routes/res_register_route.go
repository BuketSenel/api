package routes

import (
	"github.com/SelfServiceCo/api/pkg/controllers"
	"github.com/gin-gonic/gin"
	"net/http"
)

func resRegisterRoute(rg *gin.RouterGroup) {
	regGroup := rg.Group("/")

	regGroup.POST("/resRegister", func(c *gin.Context) {
		register := controllers.RestaurantRegister(c)
		if register {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound,
				gin.H{
					"status":    http.StatusNotFound,
					"message: ": "Registration failed!",
				})
			c.Abort()
		} else {
			c.JSON(http.StatusOK, gin.H{
				"status":  "200",
				"message": "OK",
				"items":   register,
				"offset":  "0",
				"limit":   "25",
			})
		}
	})
}
