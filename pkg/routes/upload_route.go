package routes

import (
	"net/http"

	"github.com/SelfServiceCo/api/pkg/controllers"
	"github.com/SelfServiceCo/api/pkg/middleware"
	"github.com/gin-gonic/gin"
)

func uploadRoute(rg *gin.RouterGroup) {
	uploadGroup := rg.Group("/")
	uploadGroup.Use(middleware.CORSMiddleware())

	uploadGroup.POST("", func(c *gin.Context) {
		file, header := controllers.UploadFile(c)
		if !file {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, header)
			c.Abort()
		} else {
			c.JSON(http.StatusOK,
				gin.H{
					"status":  http.StatusOK,
					"message": "OK",
					"items":   file,
					"offset":  "0",
					"limit":   "25",
				},
			)
		}
	})
}
