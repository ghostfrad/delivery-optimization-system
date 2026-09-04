package real

import (
	"context"

	"order-service/internal/kafka"
)

// TODO: Реализовать реальный Kafka продюсер

type KafkaProducer struct {
	brokers []string
}

func NewKafkaProducer(brokers []string) *KafkaProducer {
	return &KafkaProducer{brokers: brokers}
}

var _ kafka.MessageProducer = (*KafkaProducer)(nil)

func (k *KafkaProducer) Publish(ctx context.Context, topic string, key string, value []byte) error {
	// TODO: Реализовать
	return nil
}

func (k *KafkaProducer) Close() error {
	// TODO: Реализовать
	return nil
}