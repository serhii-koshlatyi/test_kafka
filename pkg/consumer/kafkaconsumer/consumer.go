package kafkaconsumer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/serhii-koshlatyi/test_kafka/pkg/config"
	"github.com/serhii-koshlatyi/test_kafka/pkg/logger"
	"github.com/serhii-koshlatyi/test_kafka/pkg/model"
)

// Consumer is a Kafka client for message consumption
type Consumer struct {
	consumer *kafka.Consumer
	results  chan model.TestMessage
}

func NewConsumer(conf *config.KafkaConfig) (*Consumer, error) {
	cons, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": conf.Host,
		"group.id":          conf.GroupID,
		"auto.offset.reset": "smallest"})

	if err != nil {
		return nil, fmt.Errorf("creating consumer: %w", err)
	}

	if err := cons.SubscribeTopics([]string{conf.Topic}, nil); err != nil {
		return nil, fmt.Errorf("subscribing to topics: %w", err)
	}

	return &Consumer{
		consumer: cons,
		results:  make(chan model.TestMessage),
	}, nil
}

func (c *Consumer) ResultChan() <-chan model.TestMessage {
	return c.results
}

// Consume starts listening to topic messages and writes new messages to returned channel
func (c *Consumer) Consume(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			logger.Infof("Stopped consumer loop...")

			err := c.consumer.Close()
			if err != nil {
				logger.Errorf("Closing Kafka client: %s", err)
			}

			close(c.results)
			logger.Infof("Consumer stopped")
			return
		default:
			msg, err := c.consumer.ReadMessage(100 * time.Millisecond)

			switch e := err.(type) {
			case kafka.Error:
				if e.Code() == kafka.ErrTimedOut {
					//logger.Debugf("Kafka consumer read timeout: %v", err)
					continue
				}
			}

			if err != nil {
				logger.Errorf("Kafka consumer read error: %v", err)
				continue
			}

			err = c.processKafkaMessage(ctx, msg)
			if err != nil {
				logger.Errorf("Kafka processing message error: %v", err)
				continue
			}
		}
	}
}

// processKafkaMessage provides processing message from kafka, sync state and adding to the results chanel
func (c *Consumer) processKafkaMessage(ctx context.Context, msg *kafka.Message) error {
	var message model.TestMessage

	err := json.Unmarshal(msg.Value, &message)
	if err != nil {
		logger.Errorf("Can`t unmarshal state model: %v", err)
		return err
	}

	c.results <- message
	return nil
}
