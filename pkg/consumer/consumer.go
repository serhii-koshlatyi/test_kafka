package consumer

import (
	"context"

	"github.com/serhii-koshlatyi/test_kafka/pkg/model"
)

type Consumer interface {
	Consume(ctx context.Context)
	ResultChan() <-chan model.TestMessage
}
