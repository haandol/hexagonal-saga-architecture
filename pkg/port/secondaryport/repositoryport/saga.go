package repositoryport

import (
	"context"

	"github.com/haandol/hexagonal/pkg/dto"
	"github.com/haandol/hexagonal/pkg/message/command"
)

type SagaRepository interface {
	Start(ctx context.Context, cmd *command.StartSaga) (dto.Saga, error)
	End(ctx context.Context, cmd *command.EndSaga) error
	Abort(ctx context.Context, cmd *command.AbortSaga) error
}
