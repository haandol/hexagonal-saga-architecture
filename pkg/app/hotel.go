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

type HotelApp struct {
	consumers []consumerport.Consumer
	producers []producerport.Producer
}

func NewHotelApp(
	hotelConsumer *consumer.HotelConsumer,
	hotelProducer *producer.HotelProducer,
) *HotelApp {
	consumers := []consumerport.Consumer{
		hotelConsumer,
	}
	producers := []producerport.Producer{
		hotelProducer,
	}

	return &HotelApp{
		consumers: consumers,
		producers: producers,
	}
}

func (app *HotelApp) Init() {
	logger := util.GetLogger().With(
		"module", "HotelApp",
		"func", "Init",
	)
	logger.Info("Initializing...")

	for _, c := range app.consumers {
		c.Init()
	}

	util.InitXray()
}

func (app *HotelApp) Start() {
	logger := util.GetLogger().With(
		"module", "HotelApp",
		"func", "Start",
	)
	logger.Info("Starting...")

	for _, c := range app.consumers {
		go c.Consume()
	}
}

func (app *HotelApp) Cleanup(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	logger := util.GetLogger().With(
		"module", "HotelApp",
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
