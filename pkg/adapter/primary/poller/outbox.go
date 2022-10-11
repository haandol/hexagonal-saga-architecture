package poller

import (
	"context"
	"time"

	"github.com/haandol/hexagonal/pkg/config"
	"github.com/haandol/hexagonal/pkg/service"
	"github.com/haandol/hexagonal/pkg/util"
)

type Poller struct {
	closing       chan bool
	closed        chan error
	batchSize     int
	batchInterval time.Duration
	relayService  *service.MessageRelayService
}

func NewPoller(
	cfg config.Config,
	relayService *service.MessageRelayService,
) *Poller {
	return &Poller{
		closing:       make(chan bool),
		closed:        make(chan error),
		batchSize:     cfg.Relay.FetchSize,
		batchInterval: time.Duration(cfg.Relay.FetchIntervalMil) * time.Millisecond,
		relayService:  relayService,
	}
}

func (c *Poller) Init() {
}

func (c *Poller) Poll() {
	logger := util.GetLogger().With(
		"module", "Poller",
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
				logger.Errorw("Failed to relay messages", "err", err)
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

func (c *Poller) Close(ctx context.Context) error {
	c.closing <- true
	return <-c.closed
}
