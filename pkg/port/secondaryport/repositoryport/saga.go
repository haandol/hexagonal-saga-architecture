package repositoryport

import (
	"context"

	"github.com/haandol/hexagonal/pkg/dto"
	"github.com/haandol/hexagonal/pkg/message/command"
	"github.com/haandol/hexagonal/pkg/message/event"
)

type SagaRepository interface {
	Start(ctx context.Context, cmd *command.StartSaga) (dto.Saga, error)
	ProcessCarRental(ctx context.Context, evt *event.CarRented) error
	CompensateCarRental(ctx context.Context, evt *event.CarRentalCanceled) error
	End(ctx context.Context, cmd *command.EndSaga) error
	Abort(ctx context.Context, cmd *command.AbortSaga) error
}
