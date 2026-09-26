package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"order-service/internal/domain"
)

type OrderRepository interface {
	Create(ctx context.Context, order *domain.Order) error
	GetByID(ctx context.Context, id string) (*domain.Order, error)
	UpdateStatus(ctx context.Context, id string, status domain.OrderStatus) error
	FindByStatus(ctx context.Context, status domain.OrderStatus) ([]*domain.Order, error)
	Delete(ctx context.Context, id string) error
}

type MessageProducer interface {
	Publish(ctx context.Context, topic string, key string, value []byte) error
	Close() error
}

type OrderService struct {
	repo              OrderRepository
	producer          MessageProducer
	topicOrderCreated string
}

// NewOrderService создаёт сервис.
// topicOrderCreated — имя топика Kafka для события order-created.
func NewOrderService(
	repo OrderRepository,
	producer MessageProducer,
	topicOrderCreated string,
) *OrderService {
	return &OrderService{
		repo:              repo,
		producer:          producer,
		topicOrderCreated: topicOrderCreated,
	}
}

// CreateOrder создаёт заказ и публикует событие order-created в Kafka.
func (s *OrderService) CreateOrder(
	ctx context.Context,
	customerID, address string,
	lat, lng float64,
	description string,
) (*domain.Order, error) {
	order := &domain.Order{
		ID:          uuid.NewString(),
		CustomerID:  customerID,
		Address:     address,
		Lat:         lat,
		Lng:         lng,
		Status:      domain.StatusInPool,
		Description: description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.repo.Create(ctx, order); err != nil {
		return nil, fmt.Errorf("failed to save order: %w", err)
	}

	event := map[string]interface{}{
		"order_id":    order.ID,
		"address":     order.Address,
		"lat":         order.Lat,
		"lng":         order.Lng,
		"description": order.Description,
		"timestamp":   time.Now(),
	}

	eventData, err := json.Marshal(event)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal event: %w", err)
	}

	if err := s.producer.Publish(ctx, s.topicOrderCreated, order.ID, eventData); err != nil {
		return nil, errors.Join(ErrKafkaUnavailable, err)
	}

	return order, nil
}

// GetOrder возвращает заказ по ID.
func (s *OrderService) GetOrder(ctx context.Context, id string) (*domain.Order, error) {
	return s.repo.GetByID(ctx, id)
}

// generateID генерирует уникальный ID заказа.
func generateID() string {
	return fmt.Sprintf("ORD_%s_%d",
		time.Now().Format("20060102150405"),
		time.Now().UnixNano()%10000,
	)
}
