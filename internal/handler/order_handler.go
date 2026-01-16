package handler

import (
	"net/http"
	"strconv"

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

// Method UpdateOrder adalah handler HTTP untuk endpoint:
// PUT /orders/:id
// Method ini menjadi pintu masuk request update order dari client (Postman)
func (h *OrderHandler) UpdateOrder(c *gin.Context) {

	// Mengambil parameter "id" dari URL
	// Contoh URL: /orders/2
	// c.Param("id") akan mengembalikan string "2"
	id, err := strconv.ParseInt(
		c.Param("id"), // nilai id dari URL
		10,            // basis desimal
		64,            // hasil konversi ke int64
	)

	// Jika id tidak bisa diubah ke int64
	// (misalnya id berisi huruf)
	// maka request dianggap tidak valid
	if err != nil {
		// Mengirim response HTTP 400 (Bad Request)
		// dengan pesan error ke client
		c.JSON(400, gin.H{"error": "Id salah"})

		// Menghentikan eksekusi handler
		return
	}

	// Mendeklarasikan variable cmd
	// cmd bertipe UpdateOrderCommand dari package service
	// Variable ini akan menampung data JSON dari request body
	var cmd service.UpdateOrderCommand

	// Membaca dan mem-parsing body JSON request
	// JSON dari client akan di-mapping ke struct cmd
	// Contoh body:
	// {
	//   "product": "Laptop",
	//   "quantity": 10
	// }
	if err := c.ShouldBindBodyWithJSON(&cmd); err != nil {

		// Jika JSON tidak valid atau field tidak sesuai
		// kirim response HTTP 400 ke client
		c.JSON(400, gin.H{"error": err.Error()})

		// Hentikan eksekusi handler
		return
	}

	// Memanggil method UpdateOrder pada service
	// Handler menyerahkan logic bisnis sepenuhnya ke service
	order, err := h.service.UpdateOrder(
		id,  // ID order yang akan di-update
		cmd, // Data update (product dan quantity)
	)

	// Jika service mengembalikan error
	// (misalnya order sudah DONE)
	// maka kirim response error ke client
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Jika semua proses berhasil
	// kirim response HTTP 200 (OK)
	// berisi data order yang sudah di-update
	c.JSON(200, order)
}
