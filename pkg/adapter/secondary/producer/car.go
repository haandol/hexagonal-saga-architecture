package producer

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/haandol/hexagonal/pkg/dto"
	"github.com/haandol/hexagonal/pkg/message"
	"github.com/haandol/hexagonal/pkg/message/event"
)

type CarProducer struct {
	*KafkaProducer
}

func NewCarProducer(kafkaProducer *KafkaProducer) *CarProducer {
	return &CarProducer{
		KafkaProducer: kafkaProducer,
	}
}

func (p *CarProducer) PublishCarBooked(ctx context.Context, corrID string, d dto.CarBooking) error {
	evt := &event.CarBooked{
		Message: message.Message{
			Name:          "CarBooked",
			Version:       "1.0.0",
			ID:            uuid.NewString(),
			CorrelationID: corrID,
			CreatedAt:     time.Now().Format(time.RFC3339),
		},
		Body: event.CarBookedBody{
			BookingID: d.ID,
		},
	}
	v, err := json.Marshal(evt)
	if err != nil {
		return err
	}

	if err := p.Produce(ctx, "saga-service", corrID, v); err != nil {
		return err
	}

	return nil
}

func (p *CarProducer) PublishCarBookingCanceled(ctx context.Context, corrID string, d dto.CarBooking) error {
	evt := &event.CarBookingCanceled{
		Message: message.Message{
			Name:          "CarBookingCanceled",
			Version:       "1.0.0",
			ID:            uuid.NewString(),
			CorrelationID: corrID,
			CreatedAt:     time.Now().Format(time.RFC3339),
		},
		Body: event.CarBookedBody{
			BookingID: d.ID,
		},
	}
	v, err := json.Marshal(evt)
	if err != nil {
		return err
	}

	if err := p.Produce(ctx, "saga-service", corrID, v); err != nil {
		return err
	}

	return nil
}
