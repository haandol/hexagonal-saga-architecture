package repositoryport

import (
	"context"

	"github.com/haandol/hexagonal/pkg/dto"
)

type HotelRepository interface {
	Book(ctx context.Context, d *dto.HotelBooking) (dto.HotelBooking, error)
	CancelBooking(ctx context.Context, id uint) (dto.HotelBooking, error)
}
