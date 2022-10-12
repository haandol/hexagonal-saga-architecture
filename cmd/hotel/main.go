package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/haandol/hexagonal/pkg/app"
	"github.com/haandol/hexagonal/pkg/config"
	"github.com/haandol/hexagonal/pkg/port"
	"github.com/haandol/hexagonal/pkg/util"
)

var applications []port.App

// bootstrap - register apps
func bootstrap(cfg *config.Config) {
	applications = append(applications, app.InitHotelApp(cfg))
}

func initialize() {
	for _, app := range applications {
		app.Init()
	}
}

func start() {
	for _, app := range applications {
		app.Start()
	}
}

func cleanup(ctx context.Context) {
	var wg sync.WaitGroup
	for _, app := range applications {
		wg.Add(1)
		go app.Cleanup(ctx, &wg)
	}
	wg.Wait()
}

func main() {
	logger := util.GetLogger().With(
		"module", "main",
	)

	cfg := config.Load()
	logger.Infow("\n==== Config ====\n\n", "config", cfg)

	logger.Info("Bootstraping apps...")
	bootstrap(&cfg)

	logger.Info("Initializing apps...")
	initialize()

	logger.Info("Starting apps...")
	start()

	sigs := make(chan os.Signal, 2)
	signal.Notify(sigs, os.Interrupt)
	<-sigs

	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Second*time.Duration(cfg.App.GracefulShutdownTimeout),
	)
	go func() {
		defer cancel()
		cleanup(ctx)
	}()

	select {
	case <-sigs:
		logger.Info("Received second interrupt signal; quitting without waiting for graceful close")
		os.Exit(1)
	case <-ctx.Done():
		logger.Info("Graceful close complete")
	}
}
