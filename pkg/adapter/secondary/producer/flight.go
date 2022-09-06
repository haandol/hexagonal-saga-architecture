package producer

import (
	"context"
	"encoding/json"
	"reflect"
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
			Name:          reflect.ValueOf(event.FlightBooked{}).Type().Name(),
			Version:       "1.0.0",
			ID:            uuid.NewString(),
			CorrelationID: corrID,
			CreatedAt:     time.Now().Format(time.RFC3339),
		},
		Body: event.FlightBookedBody{
			BookingID: d.ID,
		},
	}
	if err := util.ValidateStruct(evt); err != nil {
		return err
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
			Name:          reflect.ValueOf(event.FlightBookingCancelled{}).Type().Name(),
			Version:       "1.0.0",
			ID:            uuid.NewString(),
			CorrelationID: corrID,
			CreatedAt:     time.Now().Format(time.RFC3339),
		},
		Body: event.FlightBookedBody{
			BookingID: d.ID,
		},
	}
	if err := util.ValidateStruct(evt); err != nil {
		return err
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
			Name:          reflect.ValueOf(command.AbortSaga{}).Type().Name(),
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
