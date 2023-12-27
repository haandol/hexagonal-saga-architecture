package repositoryport

import (
	"context"

	"github.com/haandol/hexagonal/internal/dto"
	"github.com/haandol/hexagonal/internal/message/command"
)

type HotelRepository interface {
	BaseRepository
	Book(ctx context.Context, d *dto.HotelBooking, cmd *command.BookHotel) error
	CancelBooking(ctx context.Context, cmd *command.CancelHotelBooking) error
	PublishAbortSaga(ctx context.Context, corrID, parentID string, tripID uint, reason string) error
}
