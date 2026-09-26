package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"order-service/internal/domain"
)

// OrderRepository — реализация repository.OrderRepository поверх PostgreSQL.
type OrderRepository struct {
	db *sql.DB
}

// NewOrderRepository создаёт репозиторий.
func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

// Create сохраняет заказ в БД.
func (r *OrderRepository) Create(ctx context.Context, order *domain.Order) error {
	const query = `
		INSERT INTO orders (
			id, customer_id, address, lat, lng,
			status, description, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := r.db.ExecContext(ctx, query,
		order.ID,
		order.CustomerID,
		order.Address,
		order.Lat,
		order.Lng,
		order.Status,
		order.Description,
		order.CreatedAt,
		order.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("postgres create order: %w", err)
	}
	return nil
}

// GetByID возвращает заказ по ID или ErrOrderNotFound.
func (r *OrderRepository) GetByID(ctx context.Context, id string) (*domain.Order, error) {
	const query = `
		SELECT id, customer_id, address, lat, lng,
		       status, description, created_at, updated_at
		FROM orders
		WHERE id = $1
	`

	var o domain.Order
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&o.ID,
		&o.CustomerID,
		&o.Address,
		&o.Lat,
		&o.Lng,
		&o.Status,
		&o.Description,
		&o.CreatedAt,
		&o.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrOrderNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("postgres get order %s: %w", id, err)
	}
	return &o, nil
}

// UpdateStatus меняет статус заказа.
func (r *OrderRepository) UpdateStatus(ctx context.Context, id string, status domain.OrderStatus) error {
	const query = `
		UPDATE orders
		SET status = $1, updated_at = NOW()
		WHERE id = $2
	`

	res, err := r.db.ExecContext(ctx, query, status, id)
	if err != nil {
		return fmt.Errorf("postgres update status %s: %w", id, err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("postgres rows affected: %w", err)
	}
	if rows == 0 {
		return domain.ErrOrderNotFound
	}
	return nil
}

// FindByStatus возвращает все заказы с указанным статусом.
func (r *OrderRepository) FindByStatus(ctx context.Context, status domain.OrderStatus) ([]*domain.Order, error) {
	const query = `
		SELECT id, customer_id, address, lat, lng,
		       status, description, created_at, updated_at
		FROM orders
		WHERE status = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, status)
	if err != nil {
		return nil, fmt.Errorf("postgres find by status %s: %w", status, err)
	}
	defer rows.Close()

	// Возвращаем пустой слайс, а не nil — удобнее для потребителей.
	result := make([]*domain.Order, 0)
	for rows.Next() {
		var o domain.Order
		if err := rows.Scan(
			&o.ID,
			&o.CustomerID,
			&o.Address,
			&o.Lat,
			&o.Lng,
			&o.Status,
			&o.Description,
			&o.CreatedAt,
			&o.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("postgres scan order: %w", err)
		}
		result = append(result, &o)
	}
	// Проверяем ошибки, которые могли случиться во время итерации.
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("postgres iterate rows: %w", err)
	}
	return result, nil
}

// Delete удаляет заказ по ID.
func (r *OrderRepository) Delete(ctx context.Context, id string) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM orders WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("postgres delete order %s: %w", id, err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("postgres rows affected: %w", err)
	}
	if rows == 0 {
		return domain.ErrOrderNotFound
	}
	return nil
}
