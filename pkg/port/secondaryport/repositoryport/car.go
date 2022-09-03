package repositoryport

import (
	"context"

	"github.com/haandol/hexagonal/pkg/dto"
)

type CarRepository interface {
	Rent(ctx context.Context, d *dto.CarRental) (dto.CarRental, error)
	CancelRental(ctx context.Context, id uint) error
}
