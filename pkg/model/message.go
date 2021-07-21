package model

import "time"

// TestMessage model for Testing Producer and Consumer
type TestMessage struct {
	Broker        Broker    `json:"broker" db:"broker"`
	Message       string    `json:"message" db:"message"`
	ActionCreated time.Time `json:"action_created" db:"action_created"`
}

// Broker represent brokers
type Broker int64

// Brokers
const (
	Kafka Broker = iota
	RabbitMq
)
