package app

import (
	"context"
	"errors"
	"net/http"
	"os"
	"sync"

	"golang.org/x/sync/errgroup"

	"github.com/haandol/hexagonal/pkg/adapter/primary/consumer"
	"github.com/haandol/hexagonal/pkg/port/primaryport/consumerport"
	"github.com/haandol/hexagonal/pkg/util"
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
	logger := util.GetLogger().With(
		"module", "FlightApp",
		"func", "Init",
	)
	logger.Info("Initializing...")

	a.consumer.Init()
	logger.Info("consumers are initialized.")
}

func (a *FlightApp) Start(ctx context.Context) error {
	logger := util.GetLogger().With(
		"module", "FlightApp",
		"func", "Start",
	)
	logger.Info("Starting...")

	g := new(errgroup.Group)
	if a.server != nil {
		g.Go(func() error {
			logger.Infow("Started and serving HTTP", "addr", a.server.Addr, "pid", os.Getpid())
			if err := a.server.ListenAndServe(); err != nil {
				if errors.Is(err, http.ErrServerClosed) {
					logger.Info("server closed.")
					return err
				} else {
					logger.Errorw("ListenAndServe fail", "error", err)
					return err
				}
			}
			return nil
		})
	}
	g.Go(func() error {
		return a.consumer.Consume(ctx)
	})

	logger.Info("App Started.")

	return g.Wait()
}

func (a *FlightApp) Cleanup(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	logger := util.GetLogger().With(
		"module", "FlightApp",
		"func", "Cleanup",
	)
	logger.Info("Cleaning up...")

	if a.server != nil {
		logger.Info("Shutting down server...")
		if err := a.server.Shutdown(ctx); err != nil {
			logger.Error("Error on server shutdown:", err)
		}
		logger.Info("Server shutdown.")
	}

	logger.Info("Closing consumers...")
	if err := a.consumer.Close(ctx); err != nil {
		logger.Errorw("failed to close consumer", "err", err.Error())
	} else {
		logger.Info("Consumer closed.")
	}

	logger.Info("Cleanup done.")
}
