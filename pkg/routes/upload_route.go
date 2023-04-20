package routes

import (
	"net/http"

	"github.com/SelfServiceCo/api/pkg/controllers"
	"github.com/SelfServiceCo/api/pkg/middleware"
	"github.com/gin-gonic/gin"
)

func uploadRoute(rg *gin.RouterGroup) {
	uploadGroup := rg.Group("")
	uploadGroup.Use(middleware.CORSMiddleware())
	uploadGroup.POST("", func(c *gin.Context) {
		filePath, header := controllers.UploadFile(c)
		if filePath == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusBadRequest, header)
			c.Abort()
		} else {
			c.JSON(http.StatusOK,
				gin.H{
					"status":  http.StatusOK,
					"message": "OK",
					"items":   filePath,
					"offset":  "0",
					"limit":   "25",
				},
			)
		}
	})
}
