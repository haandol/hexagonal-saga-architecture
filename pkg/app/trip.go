package app

import (
	"context"
	"errors"
	"net/http"
	"os"
	"sync"

	"github.com/haandol/hexagonal/pkg/adapter/primary/consumer"
	"github.com/haandol/hexagonal/pkg/adapter/primary/router"
	"github.com/haandol/hexagonal/pkg/port/primaryport/consumerport"
	"github.com/haandol/hexagonal/pkg/port/primaryport/routerport"
	"github.com/haandol/hexagonal/pkg/util"
	"golang.org/x/sync/errgroup"
)

type TripApp struct {
	server      *http.Server
	routerGroup routerport.RouterGroup
	routers     []routerport.Router
	consumer    consumerport.Consumer
}

func NewTripApp(
	server *http.Server,
	ginRouter *router.GinRouter,
	tripRouter *router.TripRouter,
	tripConsumer *consumer.TripConsumer,
) *TripApp {
	routers := []routerport.Router{
		tripRouter,
	}

	return &TripApp{
		server:      server,
		routerGroup: ginRouter,
		routers:     routers,
		consumer:    tripConsumer,
	}
}

func (a *TripApp) Init() {
	logger := util.GetLogger().WithGroup("TripApp.Init")
	logger.Info("Initializing App...")

	v1 := a.routerGroup.Group("v1")
	for _, router := range a.routers {
		router.Route(v1)
	}
	logger.Info("routers are initialized.")

	a.consumer.Init()

	logger.Info("App Initialized")
}

func (a *TripApp) Start(ctx context.Context) error {
	logger := util.GetLogger().WithGroup("TripApp.Start")
	logger.Info("Starting...")

	g := new(errgroup.Group)
	if a.server != nil {
		g.Go(func() error {
			logger.Info("Started and serving HTTP", "addr", a.server.Addr, "pid", os.Getpid())
			if err := a.server.ListenAndServe(); err != nil {
				if errors.Is(err, http.ErrServerClosed) {
					logger.Info("server closed.")
					return err
				} else {
					logger.Error("listenAndServe fail", "error", err)
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

func (a *TripApp) Cleanup(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	logger := util.GetLogger().WithGroup("TripApp.Cleanup")
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
