package entity

import "time"

type OrderStatus string

const (
	OrderPending    OrderStatus = "PENDING"
	OrderProcessing OrderStatus = "PROCESSING"
	OrderDone       OrderStatus = "DONE"
	OrderFailed     OrderStatus = "FAILED"
)

type Order struct {
	ID        int64
	Product   string
	Quantity  int
	Status    OrderStatus
	CreatedAt time.Time
}
