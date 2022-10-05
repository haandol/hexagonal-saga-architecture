package pollerport

import (
	"context"
)

type Poller interface {
	Init()
	Poll()
	Close(context.Context) error
}
