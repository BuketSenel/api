package routes

import (
	"net/http"
	"strconv"

	"github.com/SelfServiceCo/api/pkg/controllers"
	"github.com/gin-gonic/gin"
)

func restaurantRoute(rg *gin.RouterGroup) {
	restGroup := rg.Group("/")

	restGroup.GET("/:resId", func(c *gin.Context) {
		resID := c.Param("resId")
		id, _ := strconv.ParseInt(resID, 16, 64)
		restaurant := controllers.GetRestaurant(id)
		if len(restaurant) == 0 {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound,
				gin.H{
					"status":    http.StatusNotFound,
					"message: ": "Restaurant not found!",
				})
			c.Abort()
		} else {
			c.JSON(http.StatusOK,
				gin.H{
					"status":  http.StatusOK,
					"message": "OK",
					"size":    len(restaurant),
					"items":   restaurant,
					"offset":  "0",
					"limit":   "25",
				},
			)
		}
	})

	restGroup.GET("/:resId/categories", func(c *gin.Context) {
		resID := c.Param("resId")
		id, _ := strconv.ParseInt(resID, 16, 64)
		categories := controllers.CategoriesByRestaurant(id)
		if len(categories) == 0 {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound,
				gin.H{
					"status":    http.StatusNotFound,
					"message: ": "No categories found!",
				})
			c.Abort()
		} else {
			c.JSON(http.StatusOK,
				gin.H{
					"status":  http.StatusOK,
					"message": "OK",
					"size":    len(categories),
					"items":   categories,
					"offset":  "0",
					"limit":   "25",
				},
			)
		}
	})

	restGroup.GET("/:resId/categories/:catID/products", func(c *gin.Context) {
		resID := c.Param("resId")
		catID := c.Param("catID")
		rid, _ := strconv.ParseInt(resID, 16, 64)
		cid, _ := strconv.ParseInt(catID, 16, 64)
		products := controllers.ProductsByCategories(cid, rid)
		if len(products) == 0 {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound,
				gin.H{
					"status":    http.StatusNotFound,
					"message: ": "No categories found!",
				},
			)
			c.Abort()
		} else {
			c.JSON(http.StatusOK,
				gin.H{
					"status":  http.StatusOK,
					"message": "OK",
					"size":    len(products),
					"items":   products,
					"offset":  "0",
					"limit":   "25",
				},
			)
		}
	})

	restGroup.GET("/:resId/products", func(c *gin.Context) {
		resID := c.Param("resId")
		rid, _ := strconv.ParseInt(resID, 16, 64)
		products := controllers.ProductsByRestaurants(rid)
		if len(products) == 0 {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound,
				gin.H{
					"status":    http.StatusNotFound,
					"message: ": "No categories found!",
				},
			)
			c.Abort()
		} else {
			c.JSON(http.StatusOK,
				gin.H{
					"status":  http.StatusOK,
					"message": "OK",
					"size":    len(products),
					"items":   products,
					"offset":  "0",
					"limit":   "25",
				},
			)
		}
	})

	restGroup.POST("/register", func(c *gin.Context) {
		register, header := controllers.RestaurantRegister(c)
		if !register {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, header)
			c.Abort()
		} else {
			c.JSON(http.StatusOK,
				gin.H{
					"status":  "200",
					"message": "OK",
					"items":   register,
					"offset":  "0",
					"limit":   "25",
				},
			)
		}
	})

	restGroup.GET("", func(c *gin.Context) {
		restaurant := controllers.GetTopRestaurants()
		if len(restaurant) == 0 {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound,
				gin.H{
					"status":    http.StatusNotFound,
					"message: ": "No restaurants retrieved!",
				},
			)
			c.Abort()
		} else {
			c.JSON(http.StatusOK,
				gin.H{
					"status":  "200",
					"message": "OK",
					"size":    len(restaurant),
					"items":   restaurant,
					"offset":  "0",
					"limit":   "25",
				},
			)
		}
	})

	restGroup.GET("/:resId/orders", func(c *gin.Context) {
		resID := c.Param("resId")
		id, _ := strconv.ParseInt(resID, 16, 64)
		orders := controllers.GetRestaurantOrders(id)
		if len(orders) == 0 {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound,
				gin.H{
					"status":    http.StatusNotFound,
					"message: ": "No restaurants retrieved!",
				},
			)
			c.Abort()
		} else {
			c.JSON(http.StatusOK,
				gin.H{
					"status":  "200",
					"message": "OK",
					"size":    len(orders),
					"items":   orders,
					"offset":  "0",
					"limit":   "25",
				},
			)
		}
	})

	restGroup.POST("/orders", func(c *gin.Context) {
		type ChangeOrder struct {
			orderId string `json:"orderId"`
			resId   string `json:"resId"`
			status  string `json:"status"`
		}
		order := new(ChangeOrder)
		c.BindJSON(&order)
		rid, _ := strconv.ParseInt(order.resId, 16, 64)
		oid, _ := strconv.ParseInt(order.orderId, 16, 64)
		orderChanged := controllers.ChangeOrderStatus(oid, rid, order.status)
		if !orderChanged {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound,
				gin.H{
					"status":    http.StatusNotFound,
					"message: ": "No restaurants retrieved!",
				},
			)
			c.Abort()
		} else {
			c.JSON(http.StatusOK,
				gin.H{
					"status":  "200",
					"message": "OK",
					"items":   orderChanged,
					"offset":  "0",
					"limit":   "25",
				},
			)
		}
	})

	restGroup.GET("/:resID/staff", func(c *gin.Context) {
		resID := c.Param("resID")
		rid, _ := strconv.ParseInt(resID, 16, 64)
		staff := controllers.GetRestaurantStaff(rid)
		if len(staff) == 0 {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound,
				gin.H{
					"status":    http.StatusNotFound,
					"message: ": "No staff found!",
				},
			)
			c.Abort()
		} else {
			c.JSON(http.StatusOK,
				gin.H{
					"status":  "200",
					"message": "OK",
					"size":    len(staff),
					"items":   staff,
					"offset":  "0",
					"limit":   "25",
				},
			)
		}
	})

	restGroup.GET("/:resID/tables", func(c *gin.Context) {
		resID := c.Param("resID")
		rid, _ := strconv.ParseInt(resID, 16, 64)
		tables := controllers.GetRestaurantTables(rid, 0)
		if len(tables) == 0 {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound,
				gin.H{
					"status":    http.StatusNotFound,
					"message: ": "No staff found!",
				},
			)
			c.Abort()
		} else {
			c.JSON(http.StatusOK,
				gin.H{
					"status":  "200",
					"message": "OK",
					"size":    len(tables),
					"items":   tables,
					"offset":  "0",
					"limit":   "25",
				},
			)
		}
	})

	restGroup.GET("/:resID/tables/:tableID", func(c *gin.Context) {
		resID := c.Param("resID")
		tableID := c.Param("tableID")
		rid, _ := strconv.ParseInt(resID, 16, 64)
		tid, _ := strconv.ParseInt(tableID, 16, 64)
		tables := controllers.GetRestaurantTables(rid, tid)
		if len(tables) == 0 {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound,
				gin.H{
					"status":    http.StatusNotFound,
					"message: ": "No staff found!",
				},
			)
			c.Abort()
		} else {
			c.JSON(http.StatusOK,
				gin.H{
					"status":  "200",
					"message": "OK",
					"size":    len(tables),
					"items":   tables,
					"offset":  "0",
					"limit":   "25",
				},
			)
		}
	})
}
