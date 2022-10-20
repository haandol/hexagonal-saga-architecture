package app

import (
	"context"
	"sync"

	"github.com/haandol/hexagonal/pkg/adapter/primary/poller"
	"github.com/haandol/hexagonal/pkg/adapter/secondary/producer"
	"github.com/haandol/hexagonal/pkg/connector/database"
	"github.com/haandol/hexagonal/pkg/port/primaryport/pollerport"
	"github.com/haandol/hexagonal/pkg/port/secondaryport/producerport"
	"github.com/haandol/hexagonal/pkg/util"
)

type MessageRelayApp struct {
	outboxPoller pollerport.Poller
	producer     producerport.Producer
}

func NewMessageRelayApp(
	outboxPoller *poller.OutboxPoller,
	kafkaProducer *producer.KafkaProducer,
) *MessageRelayApp {
	return &MessageRelayApp{
		outboxPoller: outboxPoller,
		producer:     kafkaProducer,
	}
}

func (app *MessageRelayApp) Init() {
	logger := util.GetLogger().With(
		"module", "MessageRelayApp",
		"func", "Init",
	)
	logger.Info("Initializing...")

	util.InitXray()
}

func (app *MessageRelayApp) Start() {
	logger := util.GetLogger().With(
		"module", "MessageRelayApp",
		"func", "Start",
	)
	logger.Info("Starting...")

	go app.outboxPoller.Poll()
}

func (app *MessageRelayApp) Cleanup(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	logger := util.GetLogger().With(
		"module", "MessageRelayApp",
		"func", "Cleanup",
	)
	logger.Info("Cleaning up...")

	logger.Info("Closing poller...")
	if err := app.outboxPoller.Close(ctx); err != nil {
		logger.Error("Failed to close poller", "error", err)
	}
	logger.Info("Poller stopped.")

	logger.Info("Closing database connection...")
	if err := database.Close(ctx); err != nil {
		logger.Error("Error on database close:", err)
	}
	logger.Info("Database connection closed.")

	logger.Info("Closing producer...")
	if err := app.producer.Close(ctx); err != nil {
		logger.Error("Error on producer close:", err)
	}
	logger.Info("Producer connection closed.")

	logger.Info("Cleanup done.")
}
