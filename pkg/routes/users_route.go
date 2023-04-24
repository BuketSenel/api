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

	userGroup.GET("/name", func(c *gin.Context) {
		user := c.Query("name")
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
			c.JSON(http.StatusOK, header)
		}
	})

	userGroup.GET("/orders", func(c *gin.Context) {
		uid := c.Query("userId")
		status := c.Query("status")
		userId, _ := strconv.ParseInt(uid, 10, 64)
		orders, header := controllers.GetOrdersByUser(userId, status)
		if orders == nil {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, header)
			c.Abort()
		} else {
			c.JSON(http.StatusOK, header)
		}
	})

	userGroup.GET("", func(c *gin.Context) {
		uid := c.Query("userId")
		userId, _ := strconv.ParseInt(uid, 10, 64)
		users, header := controllers.GetUser(userId)
		if users == nil {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusTeapot, header)
			c.Abort()
		} else {
			c.JSON(http.StatusOK, gin.H{
				"size":   len(users),
				"offset": "0",
				"limit":  "25",
				"items":  users,
			})
		}
	})

	userGroup.POST("/edit", func(c *gin.Context) {
		user, header := controllers.EditUser(c)
		if !user {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, header)
			c.Abort()
		} else {
			c.JSON(http.StatusOK, gin.H{
				"status":  "200",
				"message": "OK",
				"items":   user,
				"offset":  "0",
				"limit":   "25",
			},
			)
		}
	})
}
