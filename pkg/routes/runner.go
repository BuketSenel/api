package routes

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var route = gin.Default()

func StartRouting() {
	getRoutes()

	route.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "PATCH", "OPTIONS", "GET", "POST", "DELETE"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	err := route.Run(":5000")
	if err != nil {
		return
	}
}

func getRoutes() {
	restaurant := route.Group("/restaurants")
	user := route.Group("/users")
	category := route.Group("/categories")
	login := route.Group("/login")

	restaurantRoute(restaurant)
	userRoute(user)
	categoryRoute(category)
	loginRoute(login)
}
