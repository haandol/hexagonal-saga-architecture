package o11y

import "context"

type ShutdownFunc func(context.Context) error

func NoopShutdown(ctx context.Context) error {
	return nil
}
