package routes

import (
	"net/http"
	"strconv"

	"github.com/SelfServiceCo/api/pkg/middleware"
	"github.com/SelfServiceCo/api/pkg/controllers"
	"github.com/gin-gonic/gin"
)

func tableRoute(rg *gin.RouterGroup) {
	tableGroup := rg.Group("/")
	tableGroup.Use(middleware.CORSMiddleware())
	tableGroup.GET("/:tableId", func(c *gin.Context) {
		tableID := c.Param("tableId")
		id, _ := strconv.ParseInt(tableID, 16, 64)
		table, header := controllers.GetTable(id)
		if len(table) == 0 {
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

	tableGroup.GET("/:tableId/orders", func(c *gin.Context) {
		tableID := c.Param("tableId")
		id, _ := strconv.ParseInt(tableID, 16, 64)
		orders, header := controllers.OrdersByTable(id)
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
