package repositoryport

import (
	"context"

	"github.com/haandol/hexagonal/pkg/dto"
)

type TripRepository interface {
	Create(ctx context.Context, t *dto.Trip) (dto.Trip, error)
	List(ctx context.Context) ([]dto.Trip, error)
}
