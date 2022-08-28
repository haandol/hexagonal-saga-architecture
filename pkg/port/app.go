package port

import (
	"context"
	"sync"
)

type App interface {
	Init()
	Start()
	Cleanup(context.Context, *sync.WaitGroup)
}
