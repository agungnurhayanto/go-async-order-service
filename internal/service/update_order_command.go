package service

type UpdateOrderCommand struct {
	Product  string `json:"product"`
	Quantity int    `json:"quantity"`
}
