package routes

import (
	"net/http"
	"strconv"

	"github.com/SelfServiceCo/api/pkg/controllers"
	"github.com/SelfServiceCo/api/pkg/middleware"
	"github.com/gin-gonic/gin"
)

func restaurantRoute(rg *gin.RouterGroup) {
	restGroup := rg.Group("/")
	restGroup.Use(middleware.CORSMiddleware())
	restGroup.GET("/:resId", func(c *gin.Context) {
		resID := c.Param("resId")
		id, _ := strconv.ParseInt(resID, 16, 64)
		restaurant, header := controllers.GetRestaurant(id)
		if len(restaurant) == 0 {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, header)
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
		categories, header := controllers.CategoriesByRestaurant(id)
		if len(categories) == 0 {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, header)
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
		restaurant, header := controllers.GetTopRestaurants()
		if restaurant == nil {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"status": header})
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
		orders, header := controllers.GetRestaurantOrders(id)
		if *orders == nil {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, header)
			c.Abort()
		} else {
			c.JSON(http.StatusOK,
				gin.H{
					"status":  "200",
					"message": "OK",
					"size":    len(*orders),
					"items":   *orders,
					"offset":  "0",
					"limit":   "25",
				},
			)
		}
	})

	restGroup.GET("/:resId/orders/:orderId", func(c *gin.Context) {
		orderID := c.Param("orderId")
		resID := c.Param("resId")
		oid, _ := strconv.ParseInt(orderID, 16, 64)
		rid, _ := strconv.ParseInt(resID, 16, 64)
		order, _ := controllers.GetOrder(oid, rid)
		if *order == nil {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound,
				gin.H{
					"status":    http.StatusNotFound,
					"message: ": "Order not found!",
				})
			c.Abort()
		} else {
			c.JSON(http.StatusOK,
				gin.H{
					"status":  http.StatusOK,
					"message": "OK",
					"size":    len(*order),
					"items":   *order,
					"offset":  "0",
					"limit":   "25",
				},
			)
		}
	})

	restGroup.POST("/orders/alter", func(c *gin.Context) {
		orderChanged, header := controllers.ChangeOrderStatus(c)
		if !orderChanged {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, header)
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

	restGroup.POST("/staff/add", func(c *gin.Context) {
		staffCreated, header := controllers.AddStaff(c)
		if !staffCreated {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, header)
			c.Abort()
		} else {
			c.JSON(http.StatusOK,
				gin.H{
					"status":  "200",
					"message": "OK",
					"items":   staffCreated,
					"offset":  "0",
					"limit":   "25",
				},
			)
		}
	})

	restGroup.POST("/staff/delete", func(c *gin.Context) {
		staffDeleted, header := controllers.DeleteUser(c)
		if !staffDeleted {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, header)
			c.Abort()
		} else {
			c.JSON(http.StatusOK,
				gin.H{
					"status":  "200",
					"message": "OK",
					"items":   staffDeleted,
					"offset":  "0",
					"limit":   "25",
				},
			)
		}
	})

	restGroup.POST("/waiter/tip", func(c *gin.Context) {
		tipAdded, header := controllers.TippingWaiter(c)
		if !tipAdded {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, header)
			c.Abort()
		} else {
			c.JSON(http.StatusOK,
				gin.H{
					"status":  "200",
					"message": "OK",
					"items":   tipAdded,
					"offset":  "0",
					"limit":   "25",
				},
			)
		}
	})

	restGroup.GET("/:resId/waiter/:waiterId", func(c *gin.Context) {
		waiterID := c.Param("waiterId")
		resID := c.Param("resId")
		wid, _ := strconv.ParseInt(waiterID, 10, 64)
		rid, _ := strconv.ParseInt(resID, 10, 64)
		tip, header := controllers.GetTips(wid, rid)
		if tip == nil {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, header)
			c.Abort()
		} else {
			c.JSON(http.StatusOK,
				gin.H{
					"status":  http.StatusOK,
					"message": "OK",
					"size":    len(tip),
					"items":   tip,
					"offset":  "0",
					"limit":   "25",
				},
			)
		}
	})

	restGroup.GET("/:resId/tables/:tableId", func(c *gin.Context) {
		tableID := c.Param("tableId")
		resID := c.Param("resId")
		tid, _ := strconv.ParseInt(tableID, 10, 64)
		rid, _ := strconv.ParseInt(resID, 10, 64)
		table, header := controllers.GetWaitersByTable(tid, rid)
		if table == nil {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, header)
			c.Abort()
		} else {
			c.JSON(http.StatusOK,
				gin.H{
					"status":  http.StatusOK,
					"message": "OK",
					"size":    len(table),
					"items":   table,
					"offset":  "0",
					"limit":   "25",
				},
			)
		}
	})

	restGroup.GET("/:resId/staff", func(c *gin.Context) {
		resID := c.Param("resId")
		rid, _ := strconv.ParseInt(resID, 10, 64)
		staff, header := controllers.GetRestaurantStaff(rid)
		if staff == nil {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, header)
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

	restGroup.GET("/:resId/tables", func(c *gin.Context) {
		resID := c.Param("resId")
		rid, _ := strconv.ParseInt(resID, 10, 64)
		tables, header := controllers.GetRestaurantTables(rid, 0)
		if tables == nil {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, header)
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

	restGroup.GET("/:resId/:tableID/tables", func(c *gin.Context) {
		resID := c.Param("resID")
		tableID := c.Param("tableID")
		rid, _ := strconv.ParseInt(resID, 10, 64)
		tid, _ := strconv.ParseInt(tableID, 10, 64)
		tables, header := controllers.GetRestaurantTables(rid, tid)
		if len(tables) == 0 {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusTeapot, header)
			c.Abort()
		} else {
			c.JSON(http.StatusOK, header)
		}
	})

	restGroup.POST("/products", func(c *gin.Context) {
		product, header := controllers.CreateProduct(c)
		if !product {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, header)
			c.Abort()
		} else {
			c.JSON(http.StatusOK,
				gin.H{
					"status":  "200",
					"message": "OK",
					"items":   product,
					"offset":  "0",
					"limit":   "25",
				},
			)
		}
	})

	restGroup.POST("/products/edit", func(c *gin.Context) {
		product, header := controllers.EditProduct(c)
		if !product {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, header)
			c.Abort()
		} else {
			c.JSON(http.StatusOK,
				gin.H{
					"status":  "200",
					"message": "OK",
					"items":   product,
					"offset":  "0",
					"limit":   "25",
				},
			)
		}
	})

	restGroup.POST("/products/delete", func(c *gin.Context) {
		product, header := controllers.DeleteProduct(c)
		if !product {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, header)
			c.Abort()
		} else {
			c.JSON(http.StatusOK,
				gin.H{
					"status":  "200",
					"message": "OK",
					"items":   product,
					"offset":  "0",
					"limit":   "25",
				},
			)
		}
	})

	restGroup.POST("/categories", func(c *gin.Context) {
		category, header := controllers.CreateCategory(c)
		if category == 0 {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, header)
			c.Abort()
		} else {
			c.JSON(http.StatusOK,
				gin.H{
					"status":  "200",
					"message": "OK",
					"items":   category,
					"offset":  "0",
					"limit":   "25",
				},
			)
		}
	})

	restGroup.GET("/tables/:restID/waiter/:waiterID", func(c *gin.Context) {
		restID := c.Param("restID")
		waiterID := c.Param("waiterID")
		rid, _ := strconv.ParseInt(restID, 10, 64)
		tid, _ := strconv.ParseInt(waiterID, 10, 64)
		tables, header := controllers.GetWaiterTables(rid, tid)
		if len(tables) == 0 {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, header)
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

	restGroup.GET("/tables/:restID/order/:tableID", func(c *gin.Context) {
		restID := c.Param("restID")
		tableID := c.Param("tableID")
		rid, _ := strconv.ParseInt(restID, 10, 64)
		tid, _ := strconv.ParseInt(tableID, 10, 64)
		orders, header := controllers.GetWaiterOrdersByTable(rid, tid)
		if len(orders) == 0 {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, header)
			c.Abort()
		} else {
			c.JSON(http.StatusOK,
				gin.H{
					"status":  "200",
					"message": "OK",
					"size":    len(orders),
					"orders":  orders,
					"offset":  "0",
					"limit":   "25",
				},
			)
		}
	})

	restGroup.GET("/waiters/:restID", func(c *gin.Context) {
		restID := c.Param("restID")
		rid, _ := strconv.ParseInt(restID, 10, 64)
		waiters, header := controllers.GetRestaurantWaiters(rid)
		if waiters == nil {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, header)
			c.Abort()
		} else {
			c.JSON(http.StatusOK,
				gin.H{
					"status":  "200",
					"message": "OK",
					"size":    len(waiters),
					"items":   waiters,
					"offset":  "0",
					"limit":   "25",
				},
			)
		}

	})

	restGroup.POST("/tables", func(c *gin.Context) {
		tables, header := controllers.AddTable(c)
		if !tables {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, header)
			c.Abort()
		} else {
			c.JSON(http.StatusOK, header)
		}

	})

	restGroup.POST("/tables/assign", func(c *gin.Context) {
		tables, header := controllers.AssignWaiter(c)
		if !tables {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, header)
			c.Abort()
		} else {
			c.JSON(http.StatusOK, header)
		}

	})

	restGroup.POST("/tables/edit", func(c *gin.Context) {
		tables, header := controllers.EditTable(c)
		if !tables {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, header)
			c.Abort()
		} else {
			c.JSON(http.StatusOK, header)
		}

	})
}
