package producerport

import (
	"context"

	"github.com/haandol/hexagonal/pkg/dto"
)

type HotelProducer interface {
	PublishHotelBooked(ctx context.Context, corrID string, d dto.HotelBooking) error
	PublishHotelBookingCanceled(ctx context.Context, corrID string, d dto.HotelBooking) error
	PublishAbortSaga(ctx context.Context, corrID string, tripID uint, reason string) error
}
