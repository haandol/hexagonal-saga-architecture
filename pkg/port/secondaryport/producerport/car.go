package producerport

import (
	"context"

	"github.com/haandol/hexagonal/pkg/dto"
)

type CarProducer interface {
	PublishCarBooked(ctx context.Context, corrID string, d dto.CarBooking) error
	PublishCarBookingCanceled(ctx context.Context, corrID string, d dto.CarBooking) error
}
