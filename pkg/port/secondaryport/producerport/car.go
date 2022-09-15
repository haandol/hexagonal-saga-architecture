package producerport

import (
	"context"

	"github.com/haandol/hexagonal/pkg/dto"
)

type CarProducer interface {
	PublishCarBooked(ctx context.Context, corrID string, parentID string, d dto.CarBooking) error
	PublishCarBookingCancelled(ctx context.Context, corrID string, parentID string, d dto.CarBooking) error
	PublishAbortSaga(ctx context.Context, corrID string, parentID string, tripID uint, reason string) error
}
