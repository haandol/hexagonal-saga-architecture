package port

import (
	"context"
	"sync"
)

type App interface {
	Init()
	Start(context.Context) error
	Cleanup(context.Context, *sync.WaitGroup)
}
