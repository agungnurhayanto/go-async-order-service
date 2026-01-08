package service

import (
	"time"

	"github.com/agungnurhayanto/go-async-order-service/internal/entity"
	"github.com/agungnurhayanto/go-async-order-service/internal/worker"
)

type OrderService struct {
}

func NewOrderService() *OrderService {
	return &OrderService{}
}

func (s *OrderService) CreateOrder(product string, qty int) entity.Order {
	order := entity.Order{
		Product:   product,
		Quantity:  qty,
		Status:    entity.OrderPending,
		CreatedAt: time.Now(),
	}

	worker.OrderQueue <- worker.OrderJob{
		Order: order}

	return order
}
