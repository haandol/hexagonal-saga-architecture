package repositoryport

import (
	"context"

	"github.com/haandol/hexagonal/pkg/dto"
)

type OutboxRepository interface {
	QueryUnsent(ctx context.Context) ([]dto.Outbox, error)
	Delete(ctx context.Context, id uint) error
}
