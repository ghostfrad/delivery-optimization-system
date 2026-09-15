package dto

import "time"

type CreateOrderResponse struct {
	ID      string `json:"id"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

type OrderResponse struct {
	ID         string    `json:"id"`
	CustomerID string    `json:"customer_id"`
	Address    string    `json:"address"`
	Lat        float64   `json:"lat"`
	Lng        float64   `json:"lng"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}
