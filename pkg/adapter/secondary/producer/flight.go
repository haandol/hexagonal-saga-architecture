package producer

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/haandol/hexagonal/pkg/dto"
	"github.com/haandol/hexagonal/pkg/message"
	"github.com/haandol/hexagonal/pkg/message/command"
	"github.com/haandol/hexagonal/pkg/message/event"
	"github.com/haandol/hexagonal/pkg/util"
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

func (p *FlightProducer) PublishFlightBookingCancelled(ctx context.Context, corrID string, d dto.FlightBooking) error {
	evt := &event.FlightBookingCancelled{
		Message: message.Message{
			Name:          "FlightBookingCancelled",
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

func (p *FlightProducer) PublishAbortSaga(ctx context.Context, corrID string, tripID uint, reason string) error {
	cmd := &command.AbortSaga{
		Message: message.Message{
			Name:          "AbortSaga",
			Version:       "1.0.0",
			ID:            uuid.NewString(),
			CorrelationID: corrID,
			CreatedAt:     time.Now().Format(time.RFC3339),
		},
		Body: command.AbortSagaBody{
			TripID: tripID,
			Reason: reason,
			Source: "flight",
		},
	}
	if err := util.ValidateStruct(cmd); err != nil {
		return err
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
