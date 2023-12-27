package app

import (
	"context"
	"errors"
	"net/http"
	"os"
	"sync"

	"golang.org/x/sync/errgroup"

	"github.com/haandol/hexagonal/internal/adapter/primary/consumer"
	"github.com/haandol/hexagonal/internal/port/primaryport/consumerport"
	"github.com/haandol/hexagonal/pkg/util"
)

type SagaApp struct {
	server   *http.Server
	consumer consumerport.Consumer
}

func NewSagaApp(
	server *http.Server,
	sagaConsumer *consumer.SagaConsumer,
) *SagaApp {
	return &SagaApp{
		server:   server,
		consumer: sagaConsumer,
	}
}

func (a *SagaApp) Init() {
	logger := util.GetLogger().WithGroup("SagaApp.Init")
	logger.Info("Initializing App...")

	a.consumer.Init()

	logger.Info("App Initialized")
}

func (a *SagaApp) Start(ctx context.Context) error {
	logger := util.GetLogger().WithGroup("SagaApp.Start")
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
		return a.consumer.Consume(ctx)
	})

	logger.Info("App Started")

	return g.Wait()
}

func (a *SagaApp) Cleanup(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	logger := util.GetLogger().WithGroup("SagaApp.Cleanup")
	logger.Info("Cleaning App...")

	if a.server != nil {
		logger.Info("Shutting down server...")
		if err := a.server.Shutdown(ctx); err != nil {
			logger.Error("Error on server shutdown:", err)
		}
		logger.Info("Server shutdown.")
	}

	if err := a.consumer.Close(ctx); err != nil {
		logger.Error("failed to close consumer", "err", err)
	} else {
		logger.Info("Consumer closed.")
	}

	logger.Info("App Cleaned Up")
}
