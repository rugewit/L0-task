package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"subscriber/model"
	"subscriber/services/interfaces"
)

type OrderController struct {
	orderService interfaces.OrderService
	cacheService interfaces.Cache
}

func NewOrderController(orderService interfaces.OrderService, cache interfaces.Cache) *OrderController {
	return &OrderController{orderService, cache}
}

func RegisterRoutes(r *gin.Engine, orderService interfaces.OrderService, cache interfaces.Cache) {
	orderController := NewOrderController(orderService, cache)

	routes := r.Group("/order")
	routes.GET("/", orderController.GetOrders)
	routes.GET("/:id", orderController.GetOrder)
	routes.GET("/entire-cache", orderController.GetCache)
}

func (controller *OrderController) GetOrder(c *gin.Context) {
	id := c.Param("id")
	var order model.Order
	var err error
	if order, err = controller.orderService.Get(id, context.Background(), controller.cacheService); err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}
	c.JSON(http.StatusOK, &order)
}

func (controller *OrderController) GetOrders(c *gin.Context) {
	orders := make([]model.Order, 0)
	var err error

	if orders, err = controller.orderService.GetAll(context.Background()); err != nil {
		c.AbortWithError(http.StatusNotFound, err)
	}
	c.JSON(http.StatusOK, &orders)
}

func (controller *OrderController) GetCache(c *gin.Context) {
	orders := make([]model.Order, 0)

	got := controller.cacheService.GetAll()
	for _, val := range got {
		orders = append(orders, val.(model.Order))
	}
	c.JSON(http.StatusOK, &orders)
}
