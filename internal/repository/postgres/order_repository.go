package postgres

import (
	"context"
	"database/sql"

	"order-service/internal/domain"
	"order-service/internal/repository"
)

// TODO: Реализовать PostgreSQL репозиторий

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

var _ repository.OrderRepository = (*OrderRepository)(nil)

func (r *OrderRepository) Create(ctx context.Context, order *domain.Order) error {
	// TODO: Реализовать
	return nil
}

func (r *OrderRepository) GetByID(ctx context.Context, id string) (*domain.Order, error) {
	// TODO: Реализовать
	return nil, nil
}

func (r *OrderRepository) UpdateStatus(ctx context.Context, id string, status domain.OrderStatus) error {
	// TODO: Реализовать
	return nil
}

func (r *OrderRepository) FindByStatus(ctx context.Context, status domain.OrderStatus) ([]*domain.Order, error) {
	// TODO: Реализовать
	return nil, nil
}

func (r *OrderRepository) Delete(ctx context.Context, id string) error {
	// TODO: Реализовать
	return nil
}