package dto

type CreateOrderRequest struct {
	CustomerID  string  `json:"customer_id" binding:"required"`
	Address     string  `json:"address" binding:"required"`
	Lat         float64 `json:"lat"`
	Lng         float64 `json:"lng"`
	Description string  `json:"description"`
}