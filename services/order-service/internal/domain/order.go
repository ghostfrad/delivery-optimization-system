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
	ID          string
	CustomerID  string
	Address     string
	Lat         float64
	Lng         float64
	Status      OrderStatus
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
