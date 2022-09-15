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
	"github.com/haandol/hexagonal/pkg/util"
)

type TripProducer struct {
	*KafkaProducer
}

func NewTripProducer(kafkaProducer *KafkaProducer) *TripProducer {
	return &TripProducer{
		KafkaProducer: kafkaProducer,
	}
}

func (p *TripProducer) PublishStartSaga(ctx context.Context,
	corrID string, parentID string, d dto.Trip,
) error {
	cmd := command.StartSaga{
		Message: message.Message{
			Name:          reflect.ValueOf(command.StartSaga{}).Type().Name(),
			Version:       "1.0.0",
			ID:            uuid.NewString(),
			CorrelationID: corrID,
			ParentID:      parentID,
			CreatedAt:     time.Now().Format(time.RFC3339),
		},
		Body: command.StartSagaBody{
			TripID:   d.ID,
			CarID:    d.CarID,
			HotelID:  d.HotelID,
			FlightID: d.FlightID,
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

func (p *TripProducer) PublishAbortSaga(ctx context.Context,
	corrID string, parentID string, d dto.Trip,
) error {
	cmd := command.AbortSaga{
		Message: message.Message{
			Name:          reflect.ValueOf(command.AbortSaga{}).Type().Name(),
			Version:       "1.0.0",
			ID:            uuid.NewString(),
			CorrelationID: corrID,
			ParentID:      parentID,
			CreatedAt:     time.Now().Format(time.RFC3339),
		},
		Body: command.AbortSagaBody{
			TripID: d.ID,
			Reason: "user rollback",
			Source: "trip",
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
