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
	restGroup.GET("", func(c *gin.Context) {
		resID := c.Query("resId")
		id, _ := strconv.ParseInt(resID, 10, 64)
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

	restGroup.GET("/categories", func(c *gin.Context) {
		resID := c.Query("resId")
		id, _ := strconv.ParseInt(resID, 10, 64)
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

	restGroup.GET("/categories/products", func(c *gin.Context) {
		resID := c.Query("resId")
		catID := c.Query("catId")
		rid, _ := strconv.ParseInt(resID, 10, 64)
		cid, _ := strconv.ParseInt(catID, 10, 64)
		products, header := controllers.ProductsByCategories(cid, rid)
		if len(products) == 0 {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, header)
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

	restGroup.GET("/products", func(c *gin.Context) {
		resID := c.Query("resId")
		rid, _ := strconv.ParseInt(resID, 10, 64)
		products, header := controllers.ProductsByRestaurants(rid)
		if len(products) == 0 {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, header)
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

	restGroup.GET("/top", func(c *gin.Context) {
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

	restGroup.GET("/orders", func(c *gin.Context) {
		resID := c.Query("resId")
		id, _ := strconv.ParseInt(resID, 10, 64)
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

	restGroup.GET("/order", func(c *gin.Context) {
		orderID := c.Query("orderId")
		resID := c.Query("resId")
		oid, _ := strconv.ParseInt(orderID, 10, 64)
		rid, _ := strconv.ParseInt(resID, 10, 64)
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

	restGroup.POST("/waiters/tips", func(c *gin.Context) {
		tipAdded, header := controllers.TipWaiter(c)
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

	restGroup.GET("/waiters/tips", func(c *gin.Context) {
		waiterID := c.Query("waiterId")
		resID := c.Query("resId")
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

	restGroup.GET("/tables/waiters", func(c *gin.Context) {
		tableID := c.Query("tableId")
		resID := c.Query("resId")
		tid, _ := strconv.ParseInt(tableID, 10, 64)
		rid, _ := strconv.ParseInt(resID, 10, 64)
		table, header := controllers.GetWaitersByTable(rid, tid)
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

	restGroup.GET("/tables/qr", func(c *gin.Context) {
		tableID := c.Query("tableId")
		resID := c.Query("resId")
		tid, _ := strconv.ParseInt(tableID, 10, 64)
		rid, _ := strconv.ParseInt(resID, 10, 64)
		qr, header := controllers.GetQRCode(rid, tid)
		if qr == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, header)
			c.Abort()
		} else {
			c.JSON(http.StatusOK,
				gin.H{
					"status":  http.StatusOK,
					"message": "OK",
					"size":    len(qr),
					"items":   qr,
					"offset":  "0",
					"limit":   "25",
				},
			)
		}
	})

	restGroup.GET("/staff", func(c *gin.Context) {
		resID := c.Query("resId")
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

	restGroup.GET("/tables", func(c *gin.Context) {
		resID := c.Query("resId")
		tableID := c.Query("tableId")
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

	restGroup.POST("/edit", func(c *gin.Context) {
		restaurant, header := controllers.EditRestaurant(c)
		if !restaurant {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, header)
			c.Abort()
		} else {
			c.JSON(http.StatusOK,
				gin.H{
					"status":  "200",
					"message": "OK",
					"items":   restaurant,
					"offset":  "0",
					"limit":   "25",
				},
			)
		}
	})

	restGroup.GET("/waiters/tables", func(c *gin.Context) {
		restID := c.Query("resId")
		waiterID := c.Query("waiterId")
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

	restGroup.GET("/tables/orders", func(c *gin.Context) {
		restID := c.Query("resId")
		tableID := c.Query("tableId")
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

	restGroup.GET("/waiters", func(c *gin.Context) {
		restID := c.Query("resId")
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

	restGroup.GET("/popular", func(c *gin.Context) {
		restID := c.Query("restId")
		rid, _ := strconv.ParseInt(restID, 10, 64)
		popular, header := controllers.GetPopularOrders(rid)
		if !popular {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusBadRequest, header)
			c.Abort()
		} else {
			c.JSON(http.StatusOK, header)
		}
	})
}
