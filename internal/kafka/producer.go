package kafka

import (
	"GRPC_Server/internal/kafka/kafkaEntities"
	"encoding/json"
	"fmt"

	"github.com/IBM/sarama"
)

type Producer struct {
	syncProducer sarama.SyncProducer
	topic        string
}

func NewProducer(brokers []string, topic string) (*Producer, error) {
	config := sarama.NewConfig()
	// wants a message from broker leader replica acks = 1
	config.Producer.RequiredAcks = sarama.WaitForLocal
	// retying send message max 2 times
	config.Producer.Retry.Max = 2
	// waits message delivered
	config.Producer.Return.Successes = true
	// Create the Kafka Producer
	syncProducer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	return &Producer{
		syncProducer: syncProducer,
		topic:        topic,
	}, nil
}

func (p *Producer) SendMessage(topic string, kakfaRequestMessage kafkaEntities.Message) error {
	messageBytes, err := json.Marshal(kakfaRequestMessage)
	if err != nil {
		return fmt.Errorf("Failed to Marshal message")
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(messageBytes),
	}

	partition, offset, err := p.syncProducer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("Failed to send message")
	}

	fmt.Println("Partition: ", partition, " Offset: ", offset, " AnswerID:", kakfaRequestMessage.RequestType)

	return nil
}

func (p *Producer) Close() error {
	err := p.syncProducer.Close()
	if err != nil {
		return fmt.Errorf("Failed to close producer in kafka")
	}

	return nil
}
func (p *Producer) Topic() string {
	return p.topic
}
