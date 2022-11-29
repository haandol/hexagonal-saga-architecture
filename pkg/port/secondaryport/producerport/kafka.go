package producerport

import (
	"context"
)

type Producer interface {
	Produce(ctx context.Context, topic, key string, val []byte) error
}
