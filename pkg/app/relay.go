package app

import (
	"context"
	"sync"

	"github.com/haandol/hexagonal/pkg/adapter/primary/poller"
	"github.com/haandol/hexagonal/pkg/adapter/secondary/producer"
	"github.com/haandol/hexagonal/pkg/port/primaryport/pollerport"
	"github.com/haandol/hexagonal/pkg/port/secondaryport/producerport"
	"github.com/haandol/hexagonal/pkg/util"
	"golang.org/x/sync/errgroup"
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
}

func (app *MessageRelayApp) Start(ctx context.Context) error {
	logger := util.GetLogger().With(
		"module", "MessageRelayApp",
		"func", "Start",
	)
	logger.Info("Starting...")

	g := new(errgroup.Group)
	g.Go(func() error {
		return app.outboxPoller.Poll(ctx)
	})

	logger.Info("App Started.")

	return g.Wait()
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

	logger.Info("Cleanup done.")
}
