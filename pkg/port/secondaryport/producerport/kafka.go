package producerport

import (
	"context"
)

type Producer interface {
	Produce(ctx context.Context, topic string, key string, val []byte) error
	Close(context.Context) error
}
