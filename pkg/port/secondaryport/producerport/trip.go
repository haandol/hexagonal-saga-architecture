package producerport

import (
	"context"

	"github.com/haandol/hexagonal/pkg/dto"
)

type TripProducer interface {
	PublishStartSaga(ctx context.Context, corrID string, d dto.Trip) error
}
