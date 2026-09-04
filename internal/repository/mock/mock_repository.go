package mock

import (
	"context"
	"errors"
	"sync"

	"order-service/internal/domain"
	"order-service/internal/repository"
)

// мок реализация для тестов
type MockOrderRepository struct {
	mu     sync.RWMutex
	Orders map[string]*domain.Order
	Err    error
}

// гарантируем, что реализует интерфейс
var _ repository.OrderRepository = (*MockOrderRepository)(nil)

func NewMockOrderRepository() *MockOrderRepository {
	return &MockOrderRepository{
		Orders: make(map[string]*domain.Order),
	}
}

func (m *MockOrderRepository) Create(ctx context.Context, order *domain.Order) error {
	if m.Err != nil {
		return m.Err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Orders[order.ID] = order
	return nil
}

func (m *MockOrderRepository) GetByID(ctx context.Context, id string) (*domain.Order, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	order, ok := m.Orders[id]
	if !ok {
		return nil, errors.New("order not found")
	}
	return order, nil
}

func (m *MockOrderRepository) UpdateStatus(ctx context.Context, id string, status domain.OrderStatus) error {
	if m.Err != nil {
		return m.Err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	order, ok := m.Orders[id]
	if !ok {
		return errors.New("order not found")
	}
	order.Status = status
	return nil
}

func (m *MockOrderRepository) FindByStatus(ctx context.Context, status domain.OrderStatus) ([]*domain.Order, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	var result []*domain.Order
	for _, order := range m.Orders {
		if order.Status == status {
			result = append(result, order)
		}
	}
	return result, nil
}

func (m *MockOrderRepository) Delete(ctx context.Context, id string) error {
	if m.Err != nil {
		return m.Err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.Orders, id)
	return nil
}
