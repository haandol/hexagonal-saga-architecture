package producerport

import (
	"context"

	"github.com/haandol/hexagonal/pkg/dto"
)

type SagaProducer interface {
	PublishBookCar(ctx context.Context, d dto.Saga) error
	PublishCancelCarBooking(ctx context.Context, d dto.Saga) error
	PublishBookHotel(ctx context.Context, d dto.Saga) error
	PublishCancelHotelBooking(ctx context.Context, d dto.Saga) error
	PublishBookFlight(ctx context.Context, d dto.Saga) error
	PublishCancelFlightBooking(ctx context.Context, d dto.Saga) error
	PublishEndSaga(ctx context.Context, d dto.Saga) error
	PublishSagaEnded(ctx context.Context, corrID string, d dto.Saga) error
	PublishSagaAborted(ctx context.Context, corrID string, d dto.Saga) error
}
