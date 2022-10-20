package app

import (
	"context"
	"errors"
	"net/http"
	"os"
	"sync"

	"github.com/haandol/hexagonal/pkg/adapter/primary/consumer"
	"github.com/haandol/hexagonal/pkg/adapter/secondary/producer"
	"github.com/haandol/hexagonal/pkg/connector/database"
	"github.com/haandol/hexagonal/pkg/port/primaryport/consumerport"
	"github.com/haandol/hexagonal/pkg/port/secondaryport/producerport"
	"github.com/haandol/hexagonal/pkg/util"
)

type SagaApp struct {
	server       *http.Server
	sagaConsumer consumerport.Consumer
	sagaProducer producerport.Producer
}

func NewSagaApp(
	server *http.Server,
	sagaConsumer *consumer.SagaConsumer,
	sagaProducer *producer.SagaProducer,
) *SagaApp {
	return &SagaApp{
		server:       server,
		sagaConsumer: sagaConsumer,
		sagaProducer: sagaProducer,
	}
}

func (app *SagaApp) Init() {
	logger := util.GetLogger().With(
		"module", "SagaApp",
		"func", "Init",
	)
	logger.Info("Initializing...")

	logger.Info("Initializing Consumer...")
	app.sagaConsumer.Init()
	logger.Info("Consumer initialized.")

	util.InitXray()
}

func (app *SagaApp) Start() {
	logger := util.GetLogger().With(
		"module", "SagaApp",
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

	go app.sagaConsumer.Consume()
}

func (app *SagaApp) Cleanup(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	logger := util.GetLogger().With(
		"module", "SagaApp",
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

	logger.Info("Closing producer...")
	if err := app.sagaProducer.Close(ctx); err != nil {
		logger.Error("Error on producer close:", err)
	}
	logger.Info("Producer connection closed.")

	logger.Info("Closing consumer...")
	if err := app.sagaConsumer.Close(ctx); err != nil {
		logger.Error("Error on consumer close:", err)
	}
	logger.Info("Consumer connection closed.")

	logger.Info("Cleanup done.")
}
