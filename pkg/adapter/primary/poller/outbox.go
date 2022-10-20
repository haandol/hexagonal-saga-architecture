package poller

import (
	"context"
	"time"

	"github.com/haandol/hexagonal/pkg/config"
	"github.com/haandol/hexagonal/pkg/service"
	"github.com/haandol/hexagonal/pkg/util"
)

type OutboxPoller struct {
	closing       chan bool
	closed        chan error
	batchSize     int
	batchInterval time.Duration
	relayService  *service.MessageRelayService
}

func NewOutboxPoller(
	cfg *config.Config,
	relayService *service.MessageRelayService,
) *OutboxPoller {
	return &OutboxPoller{
		closing:       make(chan bool),
		closed:        make(chan error),
		batchSize:     cfg.Relay.FetchSize,
		batchInterval: time.Duration(cfg.Relay.FetchIntervalMil) * time.Millisecond,
		relayService:  relayService,
	}
}

func (c *OutboxPoller) Init() {
}

func (c *OutboxPoller) Poll() {
	logger := util.GetLogger().With(
		"module", "OutboxPoller",
		"func", "Poll",
	)

	logger.Infow("Polling outbox...", "time", time.Now().Format(time.RFC3339))
	for {
		select {
		case <-c.closing:
			c.closed <- nil
			return
		default:
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
			defer cancel()

			messages, err := c.relayService.Fetch(ctx, c.batchSize)
			if err != nil {
				logger.Errorw("Failed to fetch messages", "err", err)
				return
			}

			if err := c.relayService.Relay(ctx, messages); err != nil {
				logger.Errorw("Failed to relay messages", "err", err)
				return
			}

			time.Sleep(c.batchInterval)
		}
	}
}

func (c *OutboxPoller) Close(ctx context.Context) error {
	c.closing <- true
	return <-c.closed
}
