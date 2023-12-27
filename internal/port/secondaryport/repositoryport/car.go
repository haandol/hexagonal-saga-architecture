package repositoryport

import (
	"context"

	"github.com/haandol/hexagonal/internal/dto"
	"github.com/haandol/hexagonal/internal/message/command"
)

type CarRepository interface {
	BaseRepository
	Book(ctx context.Context, d *dto.CarBooking) (dto.CarBooking, error)
	CancelBooking(ctx context.Context, cmd *command.CancelCarBooking) error
	PublishAbortSaga(ctx context.Context, cmd *command.BookCar, reason string) error
	PublishCarBooked(ctx context.Context, corrID, parentID string, d *dto.CarBooking) error
}
