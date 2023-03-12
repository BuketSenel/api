package routes

import (
	"net/http"

	"github.com/SelfServiceCo/api/pkg/controllers"
	"github.com/SelfServiceCo/api/pkg/middleware"
	"github.com/gin-gonic/gin"
)

func loginRoute(rg *gin.RouterGroup) {
	logGroup := rg.Group("")
	logGroup.Use(middleware.CORSMiddleware())

	logGroup.POST("", func(c *gin.Context) {
		logGroup, header := controllers.Login(c)
		if !logGroup {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, header)
			c.Abort()
		} else {
			c.JSON(http.StatusOK, gin.H{
				"status":  "200",
				"message": "OK",
				"items":   logGroup,
				"offset":  "0",
				"limit":   "25",
			},
			)
		}
	})
}
