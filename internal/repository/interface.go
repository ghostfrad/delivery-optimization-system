package repository

import (
	"context"

	"order-service/internal/domain"
)

// контракт для работы с заказами
type OrderRepository interface {
	Create(ctx context.Context, order *domain.Order) error
	GetByID(ctx context.Context, id string) (*domain.Order, error)
	UpdateStatus(ctx context.Context, id string, status domain.OrderStatus) error
	FindByStatus(ctx context.Context, status domain.OrderStatus) ([]*domain.Order, error)
	Delete(ctx context.Context, id string) error
}
