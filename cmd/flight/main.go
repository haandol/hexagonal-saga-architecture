package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/haandol/hexagonal/pkg/app"
	"github.com/haandol/hexagonal/pkg/config"
	"github.com/haandol/hexagonal/pkg/connector/database"
	"github.com/haandol/hexagonal/pkg/port"
	"github.com/haandol/hexagonal/pkg/util"
	"github.com/haandol/hexagonal/pkg/util/o11y"
)

var (
	BuildTag     string
	applications []port.App
)

// bootstrap - register apps
func bootstrap(cfg *config.Config) {
	applications = append(applications,
		app.InitFlightApp(cfg),
	)
}

func initialize() {
	for _, a := range applications {
		a.Init()
	}

	o11y.InitOtel()
}

func start(ctx context.Context, ch chan error) {
	logger := util.GetLogger().With(
		"module", "start",
	)
	logger.Info("Starting apps...")

	for _, a := range applications {
		a := a
		go func() {
			if err := a.Start(ctx); err != nil {
				ch <- err
			}
		}()
	}

	logger.Info("Apps started")
}

func cleanup(ctx context.Context) {
	logger := util.GetLogger().With(
		"module", "cleanup",
	)

	var wg sync.WaitGroup
	for _, a := range applications {
		wg.Add(1)
		go a.Cleanup(ctx, &wg)
	}
	wg.Wait()

	logger.Info("Closing database connection...")
	if err := database.Close(ctx); err != nil {
		logger.Error("error on database close", "err", err)
	} else {
		logger.Info("Database connection closed.")
	}

	logger.Info("Closing o11y connection...")
	if err := o11y.Close(ctx); err != nil {
		logger.Error("error on o11y close:", err)
	} else {
		logger.Info("o11y connection closed.")
	}
}

func main() {
	cfg := config.Load()
	logger := util.InitLogger(cfg.App.Stage).With(
		"module", "main",
		"build_tag", BuildTag,
	)
	logger.Info("\n==== Config ====\n\n", "config", fmt.Sprintf("%v", cfg))

	logger.Info("Bootstraping apps...")
	bootstrap(&cfg)

	logger.Info("Initializing apps...")
	initialize()

	logger.Info("Starting apps...")
	appErr := make(chan error)
	ctx, cancel := context.WithCancel(context.Background())
	start(ctx, appErr)

	sigs := make(chan os.Signal, 2)
	signal.Notify(sigs, os.Interrupt)

	select {
	case err := <-appErr:
		cancel()
		logger.Error("error on job", "err", err)
	case <-sigs:
		cancel()
		logger.Info("User interrupt for quitting...")
	}

	ctx, cancel = context.WithTimeout(
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
