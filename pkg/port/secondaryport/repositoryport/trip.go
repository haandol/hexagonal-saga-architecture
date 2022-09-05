package repositoryport

import (
	"context"

	"github.com/haandol/hexagonal/pkg/dto"
	"github.com/haandol/hexagonal/pkg/message/event"
)

type TripRepository interface {
	Create(ctx context.Context, d *dto.Trip) (dto.Trip, error)
	List(ctx context.Context) ([]dto.Trip, error)
	UpdateBooking(ctx context.Context, evt *event.SagaEnded) error
}
