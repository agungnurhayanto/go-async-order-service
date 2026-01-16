package main

import (
	"log"

	"github.com/agungnurhayanto/go-async-order-service/internal/config"
	"github.com/agungnurhayanto/go-async-order-service/internal/handler"
	"github.com/agungnurhayanto/go-async-order-service/internal/service"
	"github.com/agungnurhayanto/go-async-order-service/internal/worker"
	"github.com/gin-gonic/gin"
)

func main() {
	// init gin

	config.ConnectDB()

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "OK",
		})
	})

	// =========================
	// INIT ASYNC ENGINE
	// =========================

	// init queue (buffer 100)
	worker.InitQueue(100)

	// start worker pool (3 worker)
	for i := 1; i <= 3; i++ {
		w := worker.Worker{ID: i}
		w.Start()
	}

	// =========================
	// INIT SERVICE & HANDLER
	// =========================

	orderService := service.NewOrderService()
	orderHandler := handler.NewOrderHandler(orderService)

	// routes
	r.POST("/orders", orderHandler.CreateOrder)
	r.GET("/orders", orderHandler.GetOrders)
	r.GET("/orders/:id", orderHandler.GetOrderByID)
	r.PUT("/orders/:id", orderHandler.UpdateOrder)

	// =========================
	// RUN SERVER
	// =========================

	log.Println("API running on :8090")
	if err := r.Run(":8090"); err != nil {
		log.Fatal(err)
	}
}
