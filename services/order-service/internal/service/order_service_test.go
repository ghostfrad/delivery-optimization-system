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
	svc := NewOrderService(mockRepo, mockProducer)

	order, err := svc.CreateOrder(
		context.Background(),
		"user123",
		"ул. Пушкина, д. 10",
		55.7558,
		37.6173,
		"Тестовый заказ",
	)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if order == nil {
		t.Fatal("Expected order, got nil")
	}
	if order.CustomerID != "user123" {
		t.Errorf("Expected CustomerID user123, got %s", order.CustomerID)
	}
	if order.Address != "ул. Пушкина, д. 10" {
		t.Errorf("Expected Address ул. Пушкина, д. 10, got %s", order.Address)
	}
	if order.Status != domain.StatusInPool {
		t.Errorf("Expected status %s, got %s", domain.StatusInPool, order.Status)
	}
	if order.Description != "Тестовый заказ" {
		t.Errorf("Expected Description 'Тестовый заказ', got %s", order.Description)
	}

	// Проверяем, что заказ сохранён в мок-репозитории
	savedOrder, err := mockRepo.GetByID(context.Background(), order.ID)
	if err != nil {
		t.Fatalf("Expected to find order, got error: %v", err)
	}
	if savedOrder.ID != order.ID {
		t.Errorf("Expected order ID %s, got %s", order.ID, savedOrder.ID)
	}

	// Проверяем, что событие опубликовано
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
	svc := NewOrderService(mockRepo, mockProducer)

	order, err := svc.CreateOrder(
		context.Background(),
		"user123",
		"ул. Пушкина, д. 10",
		0, 0, "",
	)

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
	svc := NewOrderService(mockRepo, mockProducer)

	testOrder := &domain.Order{
		ID:         "test123",
		CustomerID: "user123",
		Status:     domain.StatusPending,
	}
	mockRepo.Orders[testOrder.ID] = testOrder

	order, err := svc.GetOrder(context.Background(), "test123")

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
	svc := NewOrderService(mockRepo, mockProducer)

	order, err := svc.GetOrder(context.Background(), "nonexistent")

	if err == nil {
		t.Error("Expected error, got nil")
	}
	if order != nil {
		t.Error("Expected nil order, got order")
	}
}

// Тест на ошибку Kafka.
// ВАЖНО: поведение зависит от того, как реализован CreateOrder.
// Если ошибка Kafka возвращается — ожидаем err != nil.
// Если игнорируется — ожидаем err == nil и order != nil.
func TestOrderService_CreateOrder_KafkaError(t *testing.T) {
	mockRepo := repoMock.NewMockOrderRepository()
	mockProducer := kafkaMock.NewMockMessageProducer()
	mockProducer.Err = context.DeadlineExceeded
	svc := NewOrderService(mockRepo, mockProducer)

	order, err := svc.CreateOrder(
		context.Background(),
		"user123",
		"ул. Пушкина, д. 10",
		0, 0, "",
	)

	// Вариант 1: ошибка Kafka возвращается
	if err == nil {
		t.Error("Expected error from Kafka publish, got nil")
	}
	if order != nil {
		t.Error("Expected nil order when Kafka fails, got order")
	}

	// Вариант 2: ошибка Kafka игнорируется (раскомментируйте, если выберете этот путь)
	// if err != nil {
	//     t.Fatalf("Expected no error, got: %v", err)
	// }
	// if order == nil {
	//     t.Fatal("Expected order, got nil")
	// }
}
