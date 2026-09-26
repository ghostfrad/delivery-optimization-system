package real

import (
	"context"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

// KafkaProducer — реализация kafka.MessageProducer поверх segmentio/kafka-go.
type KafkaProducer struct {
	writer *kafka.Writer
}

// NewKafkaProducer создаёт продюсер.
func NewKafkaProducer(brokers []string, writeTimeout time.Duration) *KafkaProducer {
	return &KafkaProducer{
		writer: &kafka.Writer{
			Addr:                   kafka.TCP(brokers...),
			Balancer:               &kafka.LeastBytes{},
			RequiredAcks:           kafka.RequireAll,
			WriteTimeout:           writeTimeout,
			Async:                  false,
			BatchTimeout:           10 * time.Millisecond,
			AllowAutoTopicCreation: true,
		},
	}
}

// Publish отправляет сообщение в Kafka.
func (k *KafkaProducer) Publish(ctx context.Context, topic, key string, value []byte) error {
	msg := kafka.Message{
		Topic: topic,
		Key:   []byte(key),
		Value: value,
		Time:  time.Now(),
	}

	if err := k.writer.WriteMessages(ctx, msg); err != nil {
		return fmt.Errorf("kafka publish to %s: %w", topic, err)
	}
	return nil
}

// Close закрывает writer и сбрасывает буфер.
func (k *KafkaProducer) Close() error {
	if err := k.writer.Close(); err != nil {
		return fmt.Errorf("kafka writer close: %w", err)
	}
	return nil
}
