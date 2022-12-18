package routes

import "github.com/gin-gonic/gin"

var route = gin.Default()

func StartRouting() {
	getRoutes()
	err := route.Run(":8080")
	if err != nil {
		return
	}
}

func getRoutes() {
	restaurant := route.Group("/restaurant")
	user := route.Group("/user")
	restaurantRoute(restaurant)
	userRoute(user)
}
