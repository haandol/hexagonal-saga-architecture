package producerport

import (
	"context"

	"github.com/haandol/hexagonal/pkg/dto"
	"github.com/haandol/hexagonal/pkg/message/command"
)

type SagaProducer interface {
	PublishBookCar(ctx context.Context, d dto.Saga) error
	PublishBookHotel(ctx context.Context, d dto.Saga) error
	PublishBookFlight(ctx context.Context, d dto.Saga) error
	PublishEndSaga(ctx context.Context, d dto.Saga) error
	PublishAbortSaga(ctx context.Context, d dto.Saga, reason string, source string) error
	PublishSagaEnded(ctx context.Context, corrID string, d dto.Saga) error
	PublishSagaAborted(ctx context.Context, cmd *command.AbortSaga) error
}
