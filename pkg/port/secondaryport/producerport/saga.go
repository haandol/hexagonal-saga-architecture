package producerport

import (
	"context"

	"github.com/haandol/hexagonal/pkg/dto"
)

type SagaProducer interface {
	PublishCancelCarBooking(ctx context.Context, d dto.Saga) error
	PublishCancelHotelBooking(ctx context.Context, d dto.Saga) error
	PublishCancelFlightBooking(ctx context.Context, d dto.Saga) error
}
