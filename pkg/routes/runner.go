package routes

import (
	"github.com/SelfServiceCo/api/pkg/middleware"
	"github.com/gin-gonic/gin"
)

var route = gin.Default()

func StartRouting() {
	getRoutes()
	route.Use(middleware.CORSMiddleware())
	err := route.Run(":5000")
	if err != nil {
		return
	}
}

func getRoutes() {
	route.MaxMultipartMemory = 15 << 20 // 15 MiB
	restaurant := route.Group("/restaurants")
	user := route.Group("/users")
	category := route.Group("/categories")
	login := route.Group("/login")
	order := route.Group("/orders")
	upload := route.Group("/upload")

	restaurantRoute(restaurant)
	userRoute(user)
	categoryRoute(category)
	loginRoute(login)
	orderRoute(order)
	uploadRoute(upload)
}
