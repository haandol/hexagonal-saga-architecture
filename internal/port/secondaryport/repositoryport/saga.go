package repositoryport

import (
	"context"

	"github.com/haandol/hexagonal/internal/dto"
	"github.com/haandol/hexagonal/internal/message/command"
	"github.com/haandol/hexagonal/internal/message/event"
)

type SagaRepository interface {
	Start(ctx context.Context, cmd *command.StartSaga) error
	ProcessCarBooking(ctx context.Context, evt *event.CarBooked) error
	CompensateCarBooking(ctx context.Context, evt *event.CarBookingCanceled) error
	ProcessHotelBooking(ctx context.Context, evt *event.HotelBooked) error
	CompensateHotelBooking(ctx context.Context, evt *event.HotelBookingCanceled) (dto.Saga, error)
	ProcessFlightBooking(ctx context.Context, evt *event.FlightBooked) error
	CompensateFlightBooking(ctx context.Context, evt *event.FlightBookingCanceled) (dto.Saga, error)
	End(ctx context.Context, cmd *command.EndSaga) error
	Abort(ctx context.Context, cmd *command.AbortSaga) (dto.Saga, error)
	UpdateStatusByTripID(ctx context.Context, tripID uint, s string) error
}
