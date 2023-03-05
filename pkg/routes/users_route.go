package routes

import (
	"net/http"

	"github.com/SelfServiceCo/api/pkg/controllers"
	"github.com/gin-gonic/gin"
)

func userRoute(rg *gin.RouterGroup) {
	rg.GET("/:name", func(c *gin.Context) {
		user := c.Params.ByName("name")
		c.JSON(http.StatusOK, gin.H{"user": user})
	})

	rg.POST("/register", func(c *gin.Context) {
		register, err := controllers.UserRegister(c)
		if !register || err != nil {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound,
				gin.H{
					"status":    http.StatusNotFound,
					"message: ": "User registration failed!",
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
