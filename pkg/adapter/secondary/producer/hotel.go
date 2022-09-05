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

type HotelProducer struct {
	*KafkaProducer
}

func NewHotelProducer(kafkaProducer *KafkaProducer) *HotelProducer {
	return &HotelProducer{
		KafkaProducer: kafkaProducer,
	}
}

func (p *HotelProducer) PublishHotelBooked(ctx context.Context, corrID string, d dto.HotelBooking) error {
	evt := &event.HotelBooked{
		Message: message.Message{
			Name:          "HotelBooked",
			Version:       "1.0.0",
			ID:            uuid.NewString(),
			CorrelationID: corrID,
			CreatedAt:     time.Now().Format(time.RFC3339),
		},
		Body: event.HotelBookedBody{
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

func (p *HotelProducer) PublishHotelBookingCanceled(ctx context.Context, corrID string, d dto.HotelBooking) error {
	evt := &event.HotelBookingCanceled{
		Message: message.Message{
			Name:          "HotelBookingCanceled",
			Version:       "1.0.0",
			ID:            uuid.NewString(),
			CorrelationID: corrID,
			CreatedAt:     time.Now().Format(time.RFC3339),
		},
		Body: event.HotelBookedBody{
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

func (p *HotelProducer) PublishAbortSaga(ctx context.Context, corrID string, tripID uint, reason string) error {
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
			Source: "hotel",
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
