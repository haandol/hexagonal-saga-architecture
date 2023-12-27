package consumerport

import (
	"context"
	"time"
)

type Message struct {
	Topic     string
	Key       string
	Value     []byte
	Timestamp time.Time
}

type HandlerFunc func(context.Context, *Message) error

type Consumer interface {
	Init() // should be implemented on the concrete consumer
	RegisterHandler(h HandlerFunc) error
	Consume(context.Context) error
	Close(context.Context) error
}
