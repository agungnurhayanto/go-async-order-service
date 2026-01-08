package worker

import "github.com/agungnurhayanto/go-async-order-service/internal/entity"

type OrderJob struct {
	Order entity.Order
}

var OrderQueue chan OrderJob

func InitQueue(buffer int) {
	OrderQueue = make(chan OrderJob, buffer)
}
