package producer

import "github.com/serhii-koshlatyi/test_kafka/pkg/model"

// Producer represent interface for Message Broker producer
type Producer interface {
	Run(messages <-chan model.TestMessage)
}
