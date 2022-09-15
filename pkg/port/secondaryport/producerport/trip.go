package producerport

import (
	"context"

	"github.com/haandol/hexagonal/pkg/dto"
)

type TripProducer interface {
	PublishStartSaga(ctx context.Context, corrID string, parentID string, d dto.Trip) error
	PublishAbortSaga(ctx context.Context, corrID string, parentID string, d dto.Trip) error
}
