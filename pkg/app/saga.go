package app

import (
	"context"
	"errors"
	"net/http"
	"os"
	"sync"

	"golang.org/x/sync/errgroup"

	"github.com/haandol/hexagonal/pkg/adapter/primary/consumer"
	"github.com/haandol/hexagonal/pkg/adapter/secondary/producer"
	"github.com/haandol/hexagonal/pkg/port/primaryport/consumerport"
	"github.com/haandol/hexagonal/pkg/port/secondaryport/producerport"
	"github.com/haandol/hexagonal/pkg/util"
)

type SagaApp struct {
	server   *http.Server
	consumer consumerport.Consumer
	producer producerport.Producer
}

func NewSagaApp(
	server *http.Server,
	sagaConsumer *consumer.SagaConsumer,
	sagaProducer *producer.SagaProducer,
) *SagaApp {
	return &SagaApp{
		server:   server,
		consumer: sagaConsumer,
		producer: sagaProducer,
	}
}

func (a *SagaApp) Init() {
	logger := util.GetLogger().With(
		"module", "SagaApp",
		"func", "Init",
	)
	logger.Info("Initializing...")

	a.consumer.Init()
	logger.Info("consumers are initialized.")
}

func (a *SagaApp) Start(ctx context.Context) error {
	logger := util.GetLogger().With(
		"module", "SagaApp",
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

func (a *SagaApp) Cleanup(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	logger := util.GetLogger().With(
		"module", "SagaApp",
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

	logger.Info("Closing producer...")
	if err := a.producer.Close(ctx); err != nil {
		logger.Error("Error on producer close:", err)
	}
	logger.Info("Producer connection closed.")

	logger.Info("Closing consumers...")
	if err := a.consumer.Close(ctx); err != nil {
		logger.Errorw("failed to close consumer", "err", err.Error())
	} else {
		logger.Info("Consumer closed.")
	}

	logger.Info("Cleanup done.")
}
