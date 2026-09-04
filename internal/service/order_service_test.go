package service

import (
	"context"
	"testing"

	"order-service/internal/domain"
	kafkaMock "order-service/internal/kafka/mock"
	repoMock "order-service/internal/repository/mock"
)

func TestOrderService_CreateOrder_Success(t *testing.T) {
	mockRepo := repoMock.NewMockOrderRepository()
	mockProducer := kafkaMock.NewMockMessageProducer()
	service := NewOrderService(mockRepo, mockProducer)

	req := domain.CreateOrderRequest{
		CustomerID:  "user123",
		Address:     "ул. Пушкина, д. 10",
		Lat:         55.7558,
		Lng:         37.6173,
		Description: "Тестовый заказ",
	}

	order, err := service.CreateOrder(context.Background(), req)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if order == nil {
		t.Fatal("Expected order, got nil")
	}
	if order.CustomerID != req.CustomerID {
		t.Errorf("Expected CustomerID %s, got %s", req.CustomerID, order.CustomerID)
	}
	if order.Status != domain.StatusInPool {
		t.Errorf("Expected status %s, got %s", domain.StatusInPool, order.Status)
	}

	savedOrder, err := mockRepo.GetByID(context.Background(), order.ID)
	if err != nil {
		t.Fatalf("Expected to find order, got error: %v", err)
	}
	if savedOrder.ID != order.ID {
		t.Errorf("Expected order ID %s, got %s", order.ID, savedOrder.ID)
	}

	if len(mockProducer.PublishedMessages) != 1 {
		t.Fatalf("Expected 1 message published, got %d", len(mockProducer.PublishedMessages))
	}

	msg := mockProducer.PublishedMessages[0]
	if msg.Topic != "order-created" {
		t.Errorf("Expected topic 'order-created', got '%s'", msg.Topic)
	}
	if msg.Key != order.ID {
		t.Errorf("Expected key %s, got %s", order.ID, msg.Key)
	}
}

func TestOrderService_CreateOrder_RepositoryError(t *testing.T) {
	mockRepo := repoMock.NewMockOrderRepository()
	mockRepo.Err = context.DeadlineExceeded
	mockProducer := kafkaMock.NewMockMessageProducer()
	service := NewOrderService(mockRepo, mockProducer)

	req := domain.CreateOrderRequest{
		CustomerID: "user123",
		Address:    "ул. Пушкина, д. 10",
	}

	order, err := service.CreateOrder(context.Background(), req)

	if err == nil {
		t.Error("Expected error, got nil")
	}
	if order != nil {
		t.Error("Expected nil order, got order")
	}
}

func TestOrderService_GetOrder_Success(t *testing.T) {
	mockRepo := repoMock.NewMockOrderRepository()
	mockProducer := kafkaMock.NewMockMessageProducer()
	service := NewOrderService(mockRepo, mockProducer)

	testOrder := &domain.Order{
		ID:         "test123",
		CustomerID: "user123",
		Status:     domain.StatusPending,
	}
	mockRepo.Orders[testOrder.ID] = testOrder

	order, err := service.GetOrder(context.Background(), "test123")

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if order.ID != "test123" {
		t.Errorf("Expected ID test123, got %s", order.ID)
	}
}

func TestOrderService_GetOrder_NotFound(t *testing.T) {
	mockRepo := repoMock.NewMockOrderRepository()
	mockProducer := kafkaMock.NewMockMessageProducer()
	service := NewOrderService(mockRepo, mockProducer)

	order, err := service.GetOrder(context.Background(), "nonexistent")

	if err == nil {
		t.Error("Expected error, got nil")
	}
	if order != nil {
		t.Error("Expected nil order, got order")
	}
}

func TestOrderService_CreateOrder_KafkaError(t *testing.T) {
	mockRepo := repoMock.NewMockOrderRepository()
	mockProducer := kafkaMock.NewMockMessageProducer()
	mockProducer.Err = context.DeadlineExceeded
	service := NewOrderService(mockRepo, mockProducer)

	req := domain.CreateOrderRequest{
		CustomerID: "user123",
		Address:    "ул. Пушкина, д. 10",
	}

	order, err := service.CreateOrder(context.Background(), req)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if order == nil {
		t.Fatal("Expected order, got nil")
	}
}
