package kafka

import (
	"GRPC_Server/internal/kafka/kafkaEntities"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"sync"

	"github.com/IBM/sarama"
)

type Consumer struct {
	consumer sarama.Consumer
	messages io.Writer
}

func NewKafkaConsumer(brokers []string, message io.Writer) (*Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	fmt.Println(brokers)
	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("Unable to open kafka consumer")
	}
	return &Consumer{
		consumer: consumer,
		messages: message,
	}, nil
}
func (c *Consumer) ReadMessages(ctx context.Context, topic string) error {
	partitions, err := c.consumer.Partitions(topic)
	if err != nil {
		return fmt.Errorf("Failed to open partitions %v", err)
	}
	var wg sync.WaitGroup
	for _, partition := range partitions {
		wg.Add(1)
		go func(partition int32) {
			defer wg.Done()
			partitionConsumer, err := c.consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)
			if err != nil {
				return
			}
			defer partitionConsumer.Close()
			for {
				select {
				case msg := <-partitionConsumer.Messages():
					{
						var message kafkaEntities.Message
						if err := json.Unmarshal(msg.Value, &message); err != nil {
							log.Printf("Failed to unmarshal Kafka message: %v", err)
							continue
						}
						messageString := fmt.Sprintf("Received Kafka Message: RequestType: %s, Request: %s, Timestamp: %v\n", message.RequestType, message.Request, message.Timestamp)
						_, err := c.messages.Write([]byte(messageString))
						if err != nil {
							log.Printf("Failed to Write Kafka message: %v", err)
							continue
						}
					}
				case <-ctx.Done():
					{
						log.Println("Shutting down consumer for partition", partition)
						return
					}
				}
			}
		}(partition)
	}
	wg.Wait()
	return fmt.Errorf("Consumer canceled")
}
func (c *Consumer) ConsumerClose() error {
	err := c.consumer.Close()
	if err != nil {
		return fmt.Errorf("Failed to close consumer in kafka")
	}
	return nil
}
