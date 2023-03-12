package routes

import (
	"net/http"

	"github.com/SelfServiceCo/api/pkg/controllers"
	"github.com/SelfServiceCo/api/pkg/middleware"
	"github.com/gin-gonic/gin"
)

func categoryRoute(rg *gin.RouterGroup) {
	catGroup := rg.Group("/")
	catGroup.Use(middleware.CORSMiddleware())
	catGroup.GET("", func(c *gin.Context) {
		categories := controllers.GetCategories()
		if categories == nil || len(categories) == 0 {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound,
				gin.H{"Error: ": "Invalid starting Index on search filter!"})
			c.Abort()
		} else {
			c.JSON(http.StatusOK, gin.H{
				"size":   len(categories),
				"offset": "0",
				"limit":  "25",
				"items":  categories,
			})
		}
	})
}
