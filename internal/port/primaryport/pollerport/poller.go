package pollerport

import (
	"context"
)

type Poller interface {
	Init()
	Poll(context.Context) error
	Close(context.Context) error
}
