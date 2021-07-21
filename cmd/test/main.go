package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/serhii-koshlatyi/test_kafka/pkg/config"
	"github.com/serhii-koshlatyi/test_kafka/pkg/consumer/kafkaconsumer"
	"github.com/serhii-koshlatyi/test_kafka/pkg/logger"
	"github.com/serhii-koshlatyi/test_kafka/pkg/model"
	"github.com/serhii-koshlatyi/test_kafka/pkg/producer/kafkaproducer"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "./config/config.toml", "")
	flag.Parse()

	conf, err := config.New(configPath)
	if err != nil {
		log.Fatalf("Error read config file: %s", err.Error())
	}

	logger.NewLogger(logger.LogLevel(conf.Log.Level), conf.Log.Mode == "dev")

	ctx, shutdown := context.WithCancel(context.Background())

	messageConsumer, err := kafkaconsumer.NewConsumer(&conf.KafkaConfig)
	if err != nil {
		logger.Fatalf("Creating Kafka consumer: %s", err)
	}

	messageProducer, err := kafkaproducer.NewKafkaProducer(&conf.KafkaConfig)
	if err != nil {
		logger.Fatalf("Creating Kafka consumer: %s", err)
	}

	prodChan := make(chan *model.TestMessage)

	quitCh := make(chan os.Signal, 1)
	signal.Notify(quitCh, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		sign := <-quitCh
		logger.Infof("Received signal: %v. Running graceful shutdown...", sign)
		shutdown()
		close(prodChan)
	}()

	go messageConsumer.Consume(ctx)

	go func() {
		for i := 0; i < 100; i++ {
			test := &model.TestMessage{
				Broker:        model.Kafka,
				Message:       fmt.Sprintf("Message number %d", i),
				ActionCreated: time.Now(),
			}

			prodChan <- test
			time.Sleep(time.Millisecond * 100)
		}
	}()

	go messageProducer.Run(prodChan)

	for msg := range messageConsumer.ResultChan() {
		jsonModel, err := json.Marshal(msg)
		if err != nil {
			logger.Errorf("Error Marshal: %w", err)
		}
		logger.Infof("Consume Message: %s", string(jsonModel))
	}
}
