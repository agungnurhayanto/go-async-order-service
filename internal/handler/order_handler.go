package handler

import (
	"net/http"

	"github.com/agungnurhayanto/go-async-order-service/internal/service"
	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	service *service.OrderService
}

func NewOrderHandler(s *service.OrderService) *OrderHandler {
	return &OrderHandler{
		service: s,
	}
}

type createOrderRequest struct {
	Product  string `json:"product"`
	Quantity int    `json:"quantity"`
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req createOrderRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	order := h.service.CreateOrder(req.Product, req.Quantity)

	c.JSON(http.StatusAccepted, gin.H{
		"message": "order di terima",
		"order":   order,
	})
}

func (h *OrderHandler) GetOrders(c *gin.Context) {

	orders, err := h.service.GetOrders()

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, orders)
}

func (h *OrderHandler) GetOrderByID(c *gin.Context) {
	id := c.Param("id")

	order, err := h.service.GetOrderByID(id)

	if err != nil {
		c.JSON(404, gin.H{"error": "Order Tidak Di temukan"})
		return
	}

	c.JSON(200, order)
}
