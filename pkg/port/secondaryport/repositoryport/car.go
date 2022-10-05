package repositoryport

import (
	"context"

	"github.com/haandol/hexagonal/pkg/dto"
	"github.com/haandol/hexagonal/pkg/message/command"
)

type CarRepository interface {
	Book(ctx context.Context, d *dto.CarBooking, cmd *command.BookCar) error
	CancelBooking(ctx context.Context, cmd *command.CancelCarBooking) error
	PublishAbortSaga(ctx context.Context, corrID, parentID string, tripID uint, reason string) error
}
