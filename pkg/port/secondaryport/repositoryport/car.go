package repositoryport

import (
	"context"

	"github.com/haandol/hexagonal/pkg/dto"
)

type CarRepository interface {
	Book(ctx context.Context, d *dto.CarBooking) (dto.CarBooking, error)
	CancelBooking(ctx context.Context, id uint) (dto.CarBooking, error)
}
