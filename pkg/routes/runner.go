package routes

import "github.com/gin-gonic/gin"

var route = gin.Default()

func StartRouting() {
	getRoutes()
	err := route.Run(":5000")
	if err != nil {
		return
	}
}

func getRoutes() {
	restaurant := route.Group("/restaurants")
	user := route.Group("/users")
	category := route.Group("/categories")
	userRegister := route.Group("/user/register")
	login := route.Group("/login")

	restaurantRoute(restaurant)
	userRoute(user)
	categoryRoute(category)
	loginRoute(login)
	userRegisterRoute(userRegister)
}
