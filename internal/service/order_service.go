package service

import (
	"errors"
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

// Method UpdateOrder adalah method milik OrderService
// Method ini bertugas untuk meng-update data order (product, quantity)
// sekaligus mengubah status order menjadi PROCESSING
func (s *OrderService) UpdateOrder(

	// id adalah ID order yang akan di-update
	// Nilai ini biasanya berasal dari URL: /orders/:id
	id int64,

	// cmd adalah data input update dari client
	// Berisi product dan quantity
	// Disebut command karena ini adalah "niat perubahan", bukan data DB
	cmd UpdateOrderCommand,

) (*entity.Order, error) {

	// order adalah variable struct untuk menampung data order dari database
	// Data lama akan di-load ke sini sebelum dilakukan update
	var order entity.Order

	// QueryRow digunakan untuk mengambil satu baris data order berdasarkan id
	// Data ini diperlukan untuk:
	// 1. Memastikan order ada
	// 2. Mengecek status order saat ini
	err := config.DB.QueryRow(
		`SELECT id, product, quantity, status, created_at 
		 FROM orders 
		 WHERE id = ?`,
		id,
	).Scan(
		// Scan mengisi field-field struct order
		&order.ID,
		&order.Product,
		&order.Quantity,
		&order.Status,
		&order.CreatedAt,
	)

	// Jika query gagal (order tidak ditemukan atau error DB)
	// maka proses dihentikan dan error dikembalikan
	if err != nil {
		return nil, err
	}

	// Validasi bisnis:
	// Jika status order sudah DONE
	// maka order tidak boleh diubah lagi
	if order.Status == entity.OrderDone {
		return nil, errors.New("Order Done")
	}

	// Menentukan status baru untuk order
	// Pada desain ini:
	// - CREATE  -> PENDING
	// - UPDATE  -> PROCESSING
	// - WORKER  -> DONE
	newStatus := entity.OrderProcessing

	// Melakukan update ke database
	// Field yang di-update:
	// - product
	// - quantity
	// - status
	_, err = config.DB.Exec(
		`UPDATE orders 
		 SET product = ?, quantity = ?, status = ? 
		 WHERE id = ?`,
		// Data baru diambil dari cmd (input client)
		cmd.Product,
		cmd.Quantity,
		// Status baru hasil aturan bisnis
		newStatus,
		// ID order yang di-update
		id,
	)

	// Jika proses update gagal di database
	// maka error dikembalikan
	if err != nil {
		return nil, err
	}

	// Update struct order di memory
	// Tujuannya agar response API sesuai dengan data terbaru di database
	order.Product = cmd.Product
	order.Quantity = cmd.Quantity
	order.Status = newStatus

	// Mengembalikan pointer ke order yang sudah di-update
	// Nil error menandakan proses berhasil
	return &order, nil
}
