package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func userRoute(rg *gin.RouterGroup) {
	rg.GET("/:name", func(c *gin.Context) {
		user := c.Params.ByName("name")
		c.JSON(http.StatusOK, gin.H{"user": user})
	})
}
