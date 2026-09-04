package mock

import (
	"context"

	"order-service/internal/kafka"
)

// MockMessageProducer - мок реализация для тестов
type MockMessageProducer struct {
	PublishedMessages []Message
	Err               error
}

type Message struct {
	Topic string
	Key   string
	Value []byte
}

var _ kafka.MessageProducer = (*MockMessageProducer)(nil)

func NewMockMessageProducer() *MockMessageProducer {
	return &MockMessageProducer{
		PublishedMessages: make([]Message, 0),
	}
}

func (m *MockMessageProducer) Publish(ctx context.Context, topic string, key string, value []byte) error {
	if m.Err != nil {
		return m.Err
	}
	m.PublishedMessages = append(m.PublishedMessages, Message{
		Topic: topic,
		Key:   key,
		Value: value,
	})
	return nil
}

func (m *MockMessageProducer) Close() error {
	return nil
}