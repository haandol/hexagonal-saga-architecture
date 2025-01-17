package app

import (
	"context"
	"net/http"
	"os"
	"sync"

	"github.com/haandol/hexagonal/internal/adapter/primary/consumer"
	"github.com/haandol/hexagonal/internal/port/primaryport/consumerport"
	"github.com/haandol/hexagonal/pkg/util"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

type FlightApp struct {
	server   *http.Server
	consumer consumerport.Consumer
}

func NewFlightApp(
	server *http.Server,
	flightConsumer *consumer.FlightConsumer,
) *FlightApp {
	return &FlightApp{
		server:   server,
		consumer: flightConsumer,
	}
}

func (a *FlightApp) Init() {
	logger := util.GetLogger().WithGroup("FlightApp.Init")
	logger.Info("Initializing App...")

	a.consumer.Init()

	logger.Info("App Initialized")
}

func (a *FlightApp) Start(ctx context.Context) error {
	logger := util.GetLogger().WithGroup("FlightApp.Start")
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

func (a *FlightApp) Cleanup(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	logger := util.GetLogger().WithGroup("FlightApp.Cleanup")
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
