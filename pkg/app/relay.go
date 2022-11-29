package app

import (
	"context"
	"sync"

	"github.com/haandol/hexagonal/pkg/adapter/primary/poller"
	"github.com/haandol/hexagonal/pkg/port/primaryport/pollerport"
	"github.com/haandol/hexagonal/pkg/util"
	"golang.org/x/sync/errgroup"
)

type MessageRelayApp struct {
	outboxPoller pollerport.Poller
}

func NewMessageRelayApp(
	outboxPoller *poller.OutboxPoller,
) *MessageRelayApp {
	return &MessageRelayApp{
		outboxPoller: outboxPoller,
	}
}

func (a *MessageRelayApp) Init() {
	logger := util.GetLogger().With(
		"module", "MessageRelayApp",
		"func", "Init",
	)
	logger.Info("Initializing App...")

	a.outboxPoller.Init()

	logger.Info("App Initialized")
}

func (a *MessageRelayApp) Start(ctx context.Context) error {
	logger := util.GetLogger().With(
		"module", "MessageRelayApp",
		"func", "Start",
	)
	logger.Info("Starting App...")

	g := new(errgroup.Group)
	g.Go(func() error {
		return a.outboxPoller.Poll(ctx)
	})

	logger.Info("App Started")

	return nil
}

func (a *MessageRelayApp) Cleanup(ctx context.Context, wg *sync.WaitGroup) {
	logger := util.GetLogger().With(
		"module", "MessageRelayApp",
		"func", "Cleanup",
	)
	logger.Info("Cleaning App...")

	defer wg.Done()

	if err := a.outboxPoller.Close(ctx); err != nil {
		logger.Error("Failed to close poller", "error", err.Error())
	} else {
		logger.Info("Poller stopped.")
	}

	logger.Info("App Cleaned Up")
}
