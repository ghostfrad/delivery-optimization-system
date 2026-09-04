package kafka

import "context"

// MessageProducer - контракт для отправки сообщений
type MessageProducer interface {
	Publish(ctx context.Context, topic string, key string, value []byte) error
	Close() error
}