package routes

import (
	"github.com/gin-gonic/gin"
)

func categoryRoute(rg *gin.RouterGroup) {
	rg.Group("/categories")
	rg.GET("/", func(c *gin.Context) {

	})
}
