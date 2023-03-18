package routes

import (
	"net/http"
	"strconv"

	"github.com/SelfServiceCo/api/pkg/controllers"
	"github.com/SelfServiceCo/api/pkg/middleware"
	"github.com/gin-gonic/gin"
)

func userRoute(rg *gin.RouterGroup) {
	userGroup := rg.Group("/")
	userGroup.Use(middleware.CORSMiddleware())

	userGroup.GET("/:name", func(c *gin.Context) {
		user := c.Params.ByName("name")
		c.JSON(http.StatusOK, gin.H{"user": user})
	})

	userGroup.POST("/register", func(c *gin.Context) {
		register, header := controllers.UserRegister(c)
		if !register {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, header)
			c.Abort()
		} else {
			c.Header("Content-Type", "application/json")
			c.Header("Authorization", "JWT")
			c.JSON(http.StatusOK, gin.H{
				"status":  "200",
				"message": "OK",
				"items":   register,
				"offset":  "0",
				"limit":   "25",
			},
			)
		}
	})

	userGroup.GET("/orders", func(c *gin.Context) {
		uid := c.Param("uid")
		userId, _ := strconv.ParseInt(uid, 16, 64)
		orders, header := controllers.GetOrdersByUser(userId)
		if len(orders) == 0 {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, header)
			c.Abort()
		} else {
			c.JSON(http.StatusOK,
				gin.H{
					"status":  http.StatusOK,
					"message": "OK",
					"size":    len(orders),
					"items":   orders,
					"offset":  "0",
					"limit":   "25",
				},
			)
		}
	})
}
