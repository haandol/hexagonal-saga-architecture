package repositoryport

import "context"

type BaseRepository interface {
	BeginTx(context.Context) (context.Context, error)
	CommitTx(context.Context) error
	RollbackTx(context.Context) error
}
