package producer

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/haandol/hexagonal/pkg/dto"
	"github.com/haandol/hexagonal/pkg/message"
	"github.com/haandol/hexagonal/pkg/message/command"
)

type TripProducer struct {
	*KafkaProducer
}

func NewTripProducer(kafkaProducer *KafkaProducer) *TripProducer {
	return &TripProducer{
		KafkaProducer: kafkaProducer,
	}
}

func (p *TripProducer) PublishStartSaga(ctx context.Context, corrID string, d dto.Trip) error {
	cmd := command.StartSaga{
		Message: message.Message{
			Name:          "StartSaga",
			Version:       "1.0.0",
			ID:            uuid.NewString(),
			CorrelationID: corrID,
			CreatedAt:     time.Now().Format(time.RFC3339),
		},
		Body: command.StartSagaBody{
			TripID:   d.ID,
			CarID:    d.CarID,
			HotelID:  d.HotelID,
			FlightID: d.FlightID,
		},
	}

	v, err := json.Marshal(cmd)
	if err != nil {
		return err
	}

	if err := p.Produce(ctx, "saga-service", corrID, v); err != nil {
		return err
	}

	return nil
}
