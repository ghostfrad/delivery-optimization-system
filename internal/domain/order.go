package domain

import "time"

type OrderStatus string

const (
	StatusPending   OrderStatus = "pending"
	StatusInPool    OrderStatus = "in_pool"
	StatusAssigned  OrderStatus = "assigned"
	StatusDelivered OrderStatus = "delivered"
	StatusCancelled OrderStatus = "cancelled"
)

type Order struct {
	ID          string      `json:"id"`
	CustomerID  string      `json:"customer_id"`
	Address     string      `json:"address"`
	Lat         float64     `json:"lat"`
	Lng         float64     `json:"lng"`
	Status      OrderStatus `json:"status"`
	Description string      `json:"description"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

type CreateOrderRequest struct {
	CustomerID  string  `json:"customer_id"`
	Address     string  `json:"address"`
	Lat         float64 `json:"lat"`
	Lng         float64 `json:"lng"`
	Description string  `json:"description"`
}

type OrderResponse struct {
	ID         string      `json:"id"`
	CustomerID string      `json:"customer_id"`
	Status     OrderStatus `json:"status"`
	CreatedAt  time.Time   `json:"created_at"`
	Message    string      `json:"message,omitempty"`
}
