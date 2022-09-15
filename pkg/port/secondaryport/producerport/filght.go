package producerport

import (
	"context"

	"github.com/haandol/hexagonal/pkg/dto"
)

type FlightProducer interface {
	PublishFlightBooked(ctx context.Context, corrID string, parentID string, d dto.FlightBooking) error
	PublishFlightBookingCancelled(ctx context.Context, corrID string, parentID string, d dto.FlightBooking) error
	PublishAbortSaga(ctx context.Context, corrID string, parentID string, tripID uint, reason string) error
}
