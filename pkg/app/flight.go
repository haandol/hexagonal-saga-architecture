package app

import (
	"context"
	"errors"
	"net/http"
	"os"
	"sync"

	"github.com/haandol/hexagonal/pkg/adapter/primary/consumer"
	"github.com/haandol/hexagonal/pkg/connector/database"
	"github.com/haandol/hexagonal/pkg/port/primaryport/consumerport"
	"github.com/haandol/hexagonal/pkg/util"
)

type FlightApp struct {
	server    *http.Server
	consumers []consumerport.Consumer
}

func NewFlightApp(
	server *http.Server,
	flightConsumer *consumer.FlightConsumer,
) *FlightApp {
	consumers := []consumerport.Consumer{
		flightConsumer,
	}

	return &FlightApp{
		server:    server,
		consumers: consumers,
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
	logger.Info("consumers are initialized.")

	util.InitXray()
}

func (app *FlightApp) Start() {
	logger := util.GetLogger().With(
		"module", "FlightApp",
		"func", "Start",
	)
	logger.Info("Starting...")

	if app.server != nil {
		go func() {
			logger.Infow("Started and serving HTTP", "addr", app.server.Addr, "pid", os.Getpid())
			if err := app.server.ListenAndServe(); err != nil {
				if errors.Is(err, http.ErrServerClosed) {
					logger.Info("server closed.")
				} else {
					logger.Panicw("ListenAndServe fail", "error", err)
				}
			}
		}()
	}

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

	if app.server != nil {
		logger.Info("Shutting down server...")
		if err := app.server.Shutdown(ctx); err != nil {
			logger.Error("Error on server shutdown:", err)
		}
		logger.Info("Server shutdown.")
	}

	logger.Info("Closing database connection...")
	if err := database.Close(ctx); err != nil {
		logger.Error("Error on database close:", err)
	}
	logger.Info("Database connection closed.")

	logger.Info("Closing consumers...")
	for _, c := range app.consumers {
		c.Close(ctx)
	}
	logger.Info("Consumer connection closed.")

	logger.Info("Cleanup done.")
}
