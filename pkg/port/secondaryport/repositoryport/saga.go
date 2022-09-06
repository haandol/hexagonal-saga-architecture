package repositoryport

import (
	"context"

	"github.com/haandol/hexagonal/pkg/dto"
	"github.com/haandol/hexagonal/pkg/message/command"
	"github.com/haandol/hexagonal/pkg/message/event"
)

type SagaRepository interface {
	Start(ctx context.Context, cmd *command.StartSaga) (dto.Saga, error)
	ProcessCarBooking(ctx context.Context, evt *event.CarBooked) (dto.Saga, error)
	CompensateCarBooking(ctx context.Context, evt *event.CarBookingCancelled) (dto.Saga, error)
	ProcessHotelBooking(ctx context.Context, evt *event.HotelBooked) (dto.Saga, error)
	CompensateHotelBooking(ctx context.Context, evt *event.HotelBookingCancelled) (dto.Saga, error)
	ProcessFlightBooking(ctx context.Context, evt *event.FlightBooked) (dto.Saga, error)
	CompensateFlightBooking(ctx context.Context, evt *event.FlightBookingCancelled) (dto.Saga, error)
	End(ctx context.Context, cmd *command.EndSaga) (dto.Saga, error)
	Abort(ctx context.Context, cmd *command.AbortSaga) (dto.Saga, error)
}
