package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/haandol/hexagonal/internal/app"
	"github.com/haandol/hexagonal/internal/port"
	"github.com/haandol/hexagonal/pkg/config"
	"github.com/haandol/hexagonal/pkg/connector/database"
	"github.com/haandol/hexagonal/pkg/o11y"
	"github.com/haandol/hexagonal/pkg/util"
)

var (
	BuildTag     string
	applications []port.App
)

// bootstrap - register apps
func bootstrap(cfg *config.Config) {
	applications = append(applications,
		app.InitCarApp(cfg),
	)
}

func initialize() {
	for _, a := range applications {
		a.Init()
	}

	o11y.InitOtel()
}

func start(ctx context.Context, ch chan error) {
	logger := util.GetLogger().WithGroup("start")
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
	logger := util.GetLogger().WithGroup("cleanup")

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
	logger := util.InitLogger(cfg.App.Stage).WithGroup("main").With(
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
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)

	select {
	case err := <-appErr:
		cancel()
		logger.Error("error on job", "err", err)
	case <-sigs:
		cancel()
		logger.Info("Exiting program...")
	}

	ctx, cancel = context.WithCancel(context.Background())
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
	case <-time.After(time.Second * time.Duration(cfg.App.GracefulShutdownTimeout)):
		logger.Info("Timeout on graceful close")
	}
}
