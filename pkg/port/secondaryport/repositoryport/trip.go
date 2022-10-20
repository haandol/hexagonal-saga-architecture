package repositoryport

import (
	"context"

	"github.com/haandol/hexagonal/pkg/dto"
	"github.com/haandol/hexagonal/pkg/message/event"
)

type TripRepository interface {
	BaseRepository
	Create(ctx context.Context, corrID, parentID string, d *dto.Trip) (dto.Trip, error)
	List(ctx context.Context) ([]dto.Trip, error)
	Complete(ctx context.Context, evt *event.SagaEnded) error
	Abort(ctx context.Context, evt *event.SagaAborted) error
	GetByID(ctx context.Context, id uint) (dto.Trip, error)
	PublishStartSaga(ctx context.Context, corrID, parentID string, d *dto.Trip) error
	PublishAbortSaga(ctx context.Context, corrID, parentID string, tripID uint, reason string) error
}
