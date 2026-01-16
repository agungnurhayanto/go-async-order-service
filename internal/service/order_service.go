package service

import (
	"time"

	"github.com/agungnurhayanto/go-async-order-service/internal/config"
	"github.com/agungnurhayanto/go-async-order-service/internal/entity"
	"github.com/agungnurhayanto/go-async-order-service/internal/worker"
)

type OrderService struct {
}

func NewOrderService() *OrderService {
	return &OrderService{}
}

func (s *OrderService) CreateOrder(product string, qty int) entity.Order {

	// insert ke database

	query := `INSERT INTO orders (product, quantity, status) VALUES (?,?,?)`

	result, err := config.DB.Exec(query, product, qty, entity.OrderPending)

	if err != nil {
		panic(err)
	}

	id, _ := result.LastInsertId()

	order := entity.Order{
		ID:        id,
		Product:   product,
		Quantity:  qty,
		Status:    entity.OrderPending,
		CreatedAt: time.Now(),
	}

	worker.OrderQueue <- worker.OrderJob{
		Order: order}

	return order
}

func (s *OrderService) GetOrders() ([]entity.Order, error) {

	rows, err := config.DB.Query(`SELECT id, product, quantity, status, created_at from orders ORDER BY id DESC`)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var orders []entity.Order

	for rows.Next() {
		var o entity.Order
		rows.Scan(&o.ID, &o.Product, &o.Quantity, &o.Status, &o.CreatedAt)
		orders = append(orders, o)
	}

	return orders, nil
}

func (s *OrderService) GetOrderByID(id string) (entity.Order, error) {
	var o entity.Order
	err := config.DB.QueryRow(`
	SELECT id, product, quantity, status, created_at FROM orders WHERE id = ?`, id).Scan(&o.ID, &o.Product, &o.Quantity, &o.Status, &o.CreatedAt)

	return o, err
}
