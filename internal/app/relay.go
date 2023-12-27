package app

import (
	"context"
	"errors"
	"net/http"
	"os"
	"sync"

	"golang.org/x/sync/errgroup"

	"github.com/haandol/hexagonal/internal/adapter/primary/poller"
	"github.com/haandol/hexagonal/internal/port/primaryport/pollerport"
	"github.com/haandol/hexagonal/pkg/util"
)

type MessageRelayApp struct {
	server       *http.Server
	outboxPoller pollerport.Poller
}

func NewMessageRelayApp(
	server *http.Server,
	outboxPoller *poller.OutboxPoller,
) *MessageRelayApp {
	return &MessageRelayApp{
		server:       server,
		outboxPoller: outboxPoller,
	}
}

func (a *MessageRelayApp) Init() {
	logger := util.GetLogger().WithGroup("MessageRelayApp.Init")
	logger.Info("Initializing App...")

	a.outboxPoller.Init()

	logger.Info("App Initialized")
}

func (a *MessageRelayApp) Start(ctx context.Context) error {
	logger := util.GetLogger().WithGroup("MessageRelayApp.Start")
	logger.Info("Starting App...")

	g := new(errgroup.Group)
	if a.server != nil {
		g.Go(func() error {
			logger.Info("Started and serving HTTP", "addr", a.server.Addr, "pid", os.Getpid())
			if err := a.server.ListenAndServe(); err != nil {
				if errors.Is(err, http.ErrServerClosed) {
					logger.Info("server closed.")
					return err
				} else {
					logger.Error("ListenAndServe fail", "error", err)
					return err
				}
			}
			return nil
		})
	}
	g.Go(func() error {
		return a.outboxPoller.Poll(ctx)
	})

	logger.Info("App Started")

	return g.Wait()
}

func (a *MessageRelayApp) Cleanup(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	logger := util.GetLogger().WithGroup("MessageRelayApp.Cleanup")
	logger.Info("Cleaning App...")

	if a.server != nil {
		logger.Info("Shutting down server...")
		if err := a.server.Shutdown(ctx); err != nil {
			logger.Error("Error on server shutdown:", err)
		}
		logger.Info("Server shutdown.")
	}

	if err := a.outboxPoller.Close(ctx); err != nil {
		logger.Error("Failed to close poller", "error", err)
	} else {
		logger.Info("Poller stopped.")
	}

	logger.Info("App Cleaned Up")
}
