//go:generate mockgen -source ./kafkaInterface.go -destination=./mocks/kafkaInterface.go -package=mock_kafka
package kafka

import (
	"GRPC_Server/internal/kafka/kafkaEntities"
	"context"
)

type ProducerInterface interface {
	SendMessage(topic string, kafkaRequestMessage kafkaEntities.Message) error
	Close() error
	Topic() string
}

type ConsumerInterface interface {
	ReadMessages(ctx context.Context, topic string) error
	ConsumerClose() error
}
