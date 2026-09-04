package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"order-service/internal/domain"
	"order-service/internal/kafka"
	"order-service/internal/repository"
)

type OrderService struct {
	repo     repository.OrderRepository
	producer kafka.MessageProducer
}

func NewOrderService(repo repository.OrderRepository, producer kafka.MessageProducer) *OrderService {
	return &OrderService{
		repo:     repo,
		producer: producer,
	}
}

// создание заказа
func (s *OrderService) CreateOrder(ctx context.Context, req domain.CreateOrderRequest) (*domain.Order, error) {
	// 1. Создаем объект заказа
	order := &domain.Order{
		ID:          generateID(),
		CustomerID:  req.CustomerID,
		Address:     req.Address,
		Lat:         req.Lat,
		Lng:         req.Lng,
		Status:      domain.StatusPending,
		Description: req.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// 2. Сохраняем в БД
	if err := s.repo.Create(ctx, order); err != nil {
		return nil, fmt.Errorf("failed to save order: %w", err)
	}

	// 3. Меняем статус на "in_pool"
	order.Status = domain.StatusInPool
	if err := s.repo.UpdateStatus(ctx, order.ID, domain.StatusInPool); err != nil {
		return nil, fmt.Errorf("failed to update status: %w", err)
	}

	// 4. Отправляем событие в Kafka
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
		return order, nil
	}

	if err := s.producer.Publish(ctx, "order-created", order.ID, eventData); err != nil {
		return order, nil
	}

	return order, nil
}

// получение заказа по ID
func (s *OrderService) GetOrder(ctx context.Context, id string) (*domain.Order, error) {
	return s.repo.GetByID(ctx, id)
}

// генерация ID заказа
func generateID() string {
	return fmt.Sprintf("ORD_%s_%d",
		time.Now().Format("20060102150405"),
		time.Now().UnixNano()%10000,
	)
}
