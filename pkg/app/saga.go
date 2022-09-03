package app

import (
	"context"
	"sync"

	"github.com/haandol/hexagonal/pkg/adapter/primary/consumer"
	"github.com/haandol/hexagonal/pkg/adapter/secondary/producer"
	"github.com/haandol/hexagonal/pkg/port/primaryport/consumerport"
	"github.com/haandol/hexagonal/pkg/port/secondaryport/producerport"
	"github.com/haandol/hexagonal/pkg/util"
)

type SagaApp struct {
	consumers []consumerport.Consumer
	producers []producerport.Producer
}

func NewSagaApp(
	sagaConsumer *consumer.SagaConsumer,
	kafkaProducer *producer.KafkaProducer,
) *SagaApp {
	consumers := []consumerport.Consumer{
		sagaConsumer,
	}
	producers := []producerport.Producer{
		kafkaProducer,
	}

	return &SagaApp{
		consumers: consumers,
		producers: producers,
	}
}

func (app *SagaApp) Init() {
	logger := util.GetLogger().With(
		"module", "SagaApp",
		"func", "Init",
	)
	logger.Info("Initializing...")

	for _, c := range app.consumers {
		c.Init()
	}
}

func (app *SagaApp) Start() {
	logger := util.GetLogger().With(
		"module", "SagaApp",
		"func", "Start",
	)
	logger.Info("Starting...")

	for _, c := range app.consumers {
		go c.Consume()
	}
}

func (app *SagaApp) Cleanup(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	logger := util.GetLogger().With(
		"module", "SagaApp",
		"func", "Cleanup",
	)
	logger.Info("Cleaning up...")

	for _, c := range app.consumers {
		c.Close(ctx)
	}

	logger.Info("Closing producer...")
	for _, producer := range app.producers {
		if err := producer.Close(ctx); err != nil {
			logger.Error("Error on producer close:", err)
		}
	}
	logger.Info("Producer connection closed.")

	logger.Info("Cleanup done.")
}
