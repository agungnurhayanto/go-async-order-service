package worker

import (
	"log"
	"time"
)

type Worker struct {
	ID int
}

func (w Worker) Start() {
	go func() {
		log.Printf("Pekerjaan %d started", w.ID)

		for job := range OrderQueue {
			log.Printf(
				"Pekerjaan %d memproses order: %s",
				w.ID,
				job.Order.Product,
			)

			time.Sleep(2 * time.Second)

			log.Printf(
				"Pekerjaan %d selesai memproses order: %s",
				w.ID,
				job.Order.Product,
			)
		}
	}()
}
