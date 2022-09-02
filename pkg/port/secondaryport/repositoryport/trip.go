package repositoryport

import (
	"context"

	"github.com/haandol/hexagonal/pkg/dto"
)

type TripRepository interface {
	Create(ctx context.Context, userID uint) (dto.Trip, error)
	Update(ctx context.Context, t *dto.Trip) error
	List(ctx context.Context) ([]dto.Trip, error)
}
