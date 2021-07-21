package kafkaproducer

import (
	"encoding/json"

	"github.com/serhii-koshlatyi/test_kafka/pkg/config"
	"github.com/serhii-koshlatyi/test_kafka/pkg/logger"
	"github.com/serhii-koshlatyi/test_kafka/pkg/model"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type KafkaProducer struct {
	Producer *kafka.Producer
	Topic    string
}

func NewKafkaProducer(cfg *config.KafkaConfig) (*KafkaProducer, error) {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": cfg.Host})
	if err != nil {
		return nil, err
	}

	kafkaProducer := &KafkaProducer{
		Producer: producer,
		Topic:    cfg.Topic,
	}

	go func() { // Drain statuses of sent messages to prevent mem leak
		for range kafkaProducer.Producer.Events() {
		}
	}()

	return kafkaProducer, nil
}

func (p *KafkaProducer) Run(messages <-chan *model.TestMessage) {
	for msg := range messages {
		err := p.produceTestMessage("test_key", msg)
		if err != nil {
			logger.Errorf("Processing states message: %s", err)
		}
	}

	logger.Infof("Stopping producer...")
	p.Producer.Flush(2000)
	p.Producer.Close()
	logger.Infof("Producer stopped")
}

func (p *KafkaProducer) produceTestMessage(key string, msg *model.TestMessage) error {
	jsonState, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return p.Producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &p.Topic, Partition: kafka.PartitionAny},
		Value:          jsonState,
		Key:            []byte(key),
	}, nil)
}

func (p *KafkaProducer) produceAnyMessage(key string, value interface{}) error {
	jsonState, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return p.Producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &p.Topic, Partition: kafka.PartitionAny},
		Value:          jsonState,
		Key:            []byte(key),
	}, nil)
}
