package poller

import (
	"context"
	"time"

	"github.com/haandol/hexagonal/pkg/service"
	"github.com/haandol/hexagonal/pkg/util"
)

type Poller struct {
	closing      chan bool
	closed       chan error
	relayService *service.MessageRelayService
}

func NewPoller(
	relayService *service.MessageRelayService,
) *Poller {
	return &Poller{
		closing:      make(chan bool),
		closed:       make(chan error),
		relayService: relayService,
	}
}

func (c *Poller) Init() {}

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
			num, err := c.relayService.Relay(ctx)
			if err != nil {
				logger.Errorw("Failed to relay messages", "err", err)
				return
			}
			if num > 0 {
				logger.Infow("sent messages", "num", num)
			}
			time.Sleep(1 * time.Second)
		}
	}
}

func (c *Poller) Close(ctx context.Context) error {
	c.closing <- true
	return <-c.closed
}
