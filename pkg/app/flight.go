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

type FlightApp struct {
	consumers []consumerport.Consumer
	producers []producerport.Producer
}

func NewFlightApp(
	flightConsumer *consumer.FlightConsumer,
	flightProducer *producer.FlightProducer,
) *FlightApp {
	consumers := []consumerport.Consumer{
		flightConsumer,
	}
	producers := []producerport.Producer{
		flightProducer,
	}

	return &FlightApp{
		consumers: consumers,
		producers: producers,
	}
}

func (app *FlightApp) Init() {
	logger := util.GetLogger().With(
		"module", "FlightApp",
		"func", "Init",
	)
	logger.Info("Initializing...")

	for _, c := range app.consumers {
		c.Init()
	}

	util.InitXray()
}

func (app *FlightApp) Start() {
	logger := util.GetLogger().With(
		"module", "FlightApp",
		"func", "Start",
	)
	logger.Info("Starting...")

	for _, c := range app.consumers {
		go c.Consume()
	}
}

func (app *FlightApp) Cleanup(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	logger := util.GetLogger().With(
		"module", "FlightApp",
		"func", "Cleanup",
	)
	logger.Info("Cleaning up...")

	logger.Info("Closing producer...")
	for _, producer := range app.producers {
		if err := producer.Close(ctx); err != nil {
			logger.Error("Error on producer close:", err)
		}
	}
	logger.Info("Producer connection closed.")

	logger.Info("Closing consumers...")
	for _, c := range app.consumers {
		c.Close(ctx)
	}
	logger.Info("Consumer connection closed.")

	logger.Info("Cleanup done.")
}
