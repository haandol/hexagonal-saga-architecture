package producerport

import (
	"context"

	"github.com/haandol/hexagonal/pkg/dto"
)

type FlightProducer interface {
	PublishFlightBooked(ctx context.Context, corrID string, d dto.FlightBooking) error
	PublishFlightBookingCanceled(ctx context.Context, corrID string, d dto.FlightBooking) error
}
