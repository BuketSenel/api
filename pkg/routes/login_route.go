package routes

import (
	"github.com/SelfServiceCo/api/pkg/controllers"
	"github.com/gin-gonic/gin"
	"net/http"
)

func loginRoute(rg *gin.RouterGroup) {
	logGroup := rg.Group("/")

	logGroup.GET("", func(c *gin.Context) {
		logGroup := controllers.Login(c)
		if logGroup {
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
				"items":   logGroup,
				"offset":  "0",
				"limit":   "25",
			})
		}
	})
}
