package repositoryport

import (
	"context"

	"github.com/haandol/hexagonal/pkg/dto"
	"github.com/haandol/hexagonal/pkg/message/command"
)

type FlightRepository interface {
	BaseRepository
	Book(ctx context.Context, d *dto.FlightBooking, cmd *command.BookFlight) error
	CancelBooking(ctx context.Context, cmd *command.CancelFlightBooking) error
	PublishAbortSaga(ctx context.Context, corrID, parentID string, tripID uint, reason string) error
}
