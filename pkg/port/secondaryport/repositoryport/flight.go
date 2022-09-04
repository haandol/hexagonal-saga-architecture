package repositoryport

import (
	"context"

	"github.com/haandol/hexagonal/pkg/dto"
)

type FlightRepository interface {
	Book(ctx context.Context, d *dto.FlightBooking) (dto.FlightBooking, error)
	CancelBooking(ctx context.Context, id uint) error
}
