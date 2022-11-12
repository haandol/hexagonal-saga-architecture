package poller

import (
	"context"
	"time"

	"github.com/haandol/hexagonal/pkg/config"
	"github.com/haandol/hexagonal/pkg/service"
	"github.com/haandol/hexagonal/pkg/util"
)

type OutboxPoller struct {
	batchSize     int
	batchInterval time.Duration
	relayService  *service.MessageRelayService
}

func NewOutboxPoller(
	cfg *config.Config,
	relayService *service.MessageRelayService,
) *OutboxPoller {
	return &OutboxPoller{
		batchSize:     cfg.Relay.FetchSize,
		batchInterval: time.Duration(cfg.Relay.FetchIntervalMil) * time.Millisecond,
		relayService:  relayService,
	}
}

func (c *OutboxPoller) Init() {}

func (c *OutboxPoller) Poll(ctx context.Context) error {
	logger := util.GetLogger().With(
		"module", "OutboxPoller",
		"func", "Poll",
	)

	logger.Infow("Polling outbox...", "time", time.Now().Format(time.RFC3339))
	jobDone := make(chan error)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
			go func() {
				defer cancel()
				messages, err := c.relayService.Fetch(ctx, c.batchSize)
				if err != nil {
					logger.Errorw("Failed to fetch messages", "err", err)
					jobDone <- err
				}
				logger.Debugw("Fetched messages", "messages", messages)

				if len(messages) > 0 {
					if err := c.relayService.Relay(ctx, messages); err != nil {
						logger.Errorw("Failed to relay messages", "err", err)
						jobDone <- err
					}
				}

				jobDone <- nil
			}()
			select {
			case err := <-jobDone:
				logger.Infow("jobDone", "err", err)
				if err != nil {
					return err
				}
			case <-ctx.Done():
				logger.Infow("ctx.Done", "err", ctx.Err())
			}

			logger.Debugw("Sleep...", "interval", c.batchInterval)
			time.Sleep(c.batchInterval)
		}
	}
}

func (c *OutboxPoller) Close(ctx context.Context) error {
	logger := util.GetLogger().With(
		"module", "OutboxPoller",
		"func", "Close",
	)
	logger.Info("Closing poller...")
	return nil
}
