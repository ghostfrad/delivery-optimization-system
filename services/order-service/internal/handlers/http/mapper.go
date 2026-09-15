package http

import (
	"order-service/internal/domain"
	"order-service/internal/handlers/http/dto"
)

func toDTOOrder(order *domain.Order) dto.OrderResponse {
	return dto.OrderResponse{
		ID:         order.ID,
		CustomerID: order.CustomerID,
		Address:    order.Address,
		Lat:        order.Lat,
		Lng:        order.Lng,
		Status:     string(order.Status),
		CreatedAt:  order.CreatedAt,
	}
}

func toDTOCreateOrder(order *domain.Order) dto.CreateOrderResponse {
	return dto.CreateOrderResponse{
		ID:      order.ID,
		Status:  string(order.Status),
		Message: "Заказ принят",
	}
}
