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

type FlightProducer struct {
	*KafkaProducer
}

func NewFlightProducer(kafkaProducer *KafkaProducer) *FlightProducer {
	return &FlightProducer{
		KafkaProducer: kafkaProducer,
	}
}

func (p *FlightProducer) PublishFlightBooked(ctx context.Context, corrID string, d dto.FlightBooking) error {
	evt := &event.FlightBooked{
		Message: message.Message{
			Name:          "FlightBooked",
			Version:       "1.0.0",
			ID:            uuid.NewString(),
			CorrelationID: corrID,
			CreatedAt:     time.Now().Format(time.RFC3339),
		},
		Body: event.FlightBookedBody{
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

func (p *FlightProducer) PublishFlightBookingCanceled(ctx context.Context, corrID string, d dto.FlightBooking) error {
	evt := &event.FlightBookingCanceled{
		Message: message.Message{
			Name:          "FlightBookingCanceled",
			Version:       "1.0.0",
			ID:            uuid.NewString(),
			CorrelationID: corrID,
			CreatedAt:     time.Now().Format(time.RFC3339),
		},
		Body: event.FlightBookedBody{
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
