package kafkaEntities

import "time"

// Message is a struct to represent the message to be sent to Kafka.
type Message struct {
	RequestType string    `json:"request_type"`
	Request     string    `json:"request"`
	Timestamp   time.Time `json:"timestamp"`
}
