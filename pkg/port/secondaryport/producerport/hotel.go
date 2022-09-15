package producerport

import (
	"context"

	"github.com/haandol/hexagonal/pkg/dto"
)

type HotelProducer interface {
	PublishHotelBooked(ctx context.Context, corrID string, parentID string, d dto.HotelBooking) error
	PublishHotelBookingCancelled(ctx context.Context, corrID string, parentID string, d dto.HotelBooking) error
	PublishAbortSaga(ctx context.Context, corrID string, parentID string, tripID uint, reason string) error
}
