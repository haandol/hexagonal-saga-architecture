package publisher

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/haandol/hexagonal/pkg/dto"
	"github.com/haandol/hexagonal/pkg/message"
	"github.com/haandol/hexagonal/pkg/message/command"
	"github.com/haandol/hexagonal/pkg/port/secondaryport/producerport"
	"github.com/haandol/hexagonal/pkg/util"
)

func PublishBookCar(ctx context.Context, p producerport.Producer, d dto.Saga) error {
	cmd := &command.BookCar{
		Message: message.Message{
			Name:          "BookCar",
			Version:       "1.0.0",
			ID:            uuid.NewString(),
			CorrelationID: d.CorrelationID,
			CreatedAt:     time.Now().Format(time.RFC3339),
		},
		Body: command.BookCarBody{
			TripID: d.TripID,
			CarID:  d.CarID,
		},
	}
	if err := util.ValidateStruct(cmd); err != nil {
		return err
	}
	v, err := json.Marshal(cmd)
	if err != nil {
		return err
	}

	if err := p.Produce(ctx, "car-service", d.CorrelationID, v); err != nil {
		return err
	}

	return nil
}

func PublishBookHotel(ctx context.Context, p producerport.Producer, d dto.Saga) error {
	cmd := &command.BookHotel{
		Message: message.Message{
			Name:          "BookHotel",
			Version:       "1.0.0",
			ID:            uuid.NewString(),
			CorrelationID: d.CorrelationID,
			CreatedAt:     time.Now().Format(time.RFC3339),
		},
		Body: command.BookHotelBody{
			TripID:  d.TripID,
			HotelID: d.HotelID,
		},
	}
	if err := util.ValidateStruct(cmd); err != nil {
		return err
	}
	v, err := json.Marshal(cmd)
	if err != nil {
		return err
	}

	if err := p.Produce(ctx, "hotel-service", d.CorrelationID, v); err != nil {
		return err
	}

	return nil
}

func PublishBookFlight(ctx context.Context, p producerport.Producer, d dto.Saga) error {
	cmd := &command.BookFlight{
		Message: message.Message{
			Name:          "BookFlight",
			Version:       "1.0.0",
			ID:            uuid.NewString(),
			CorrelationID: d.CorrelationID,
			CreatedAt:     time.Now().Format(time.RFC3339),
		},
		Body: command.BookFlightBody{
			TripID:   d.TripID,
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

	if err := p.Produce(ctx, "flight-service", d.CorrelationID, v); err != nil {
		return err
	}

	return nil
}

func PublishEndSaga(ctx context.Context, p producerport.Producer, d dto.Saga) error {
	cmd := &command.EndSaga{
		Message: message.Message{
			Name:          "EndSaga",
			Version:       "1.0.0",
			ID:            uuid.NewString(),
			CorrelationID: d.CorrelationID,
			CreatedAt:     time.Now().Format(time.RFC3339),
		},
		Body: command.EndSagaBody{
			SagaID: d.ID,
		},
	}
	if err := util.ValidateStruct(cmd); err != nil {
		return err
	}
	v, err := json.Marshal(cmd)
	if err != nil {
		return err
	}

	if err := p.Produce(ctx, "saga-service", d.CorrelationID, v); err != nil {
		return err
	}

	return nil
}

func PublishAbortSaga(ctx context.Context, p producerport.Producer, d dto.Saga, reason string, source string) error {
	cmd := &command.AbortSaga{
		Message: message.Message{
			Name:          "AbortSaga",
			Version:       "1.0.0",
			ID:            uuid.NewString(),
			CorrelationID: d.CorrelationID,
			CreatedAt:     time.Now().Format(time.RFC3339),
		},
		Body: command.AbortSagaBody{
			SagaID: d.ID,
			Reason: reason,
			Source: source,
		},
	}
	if err := util.ValidateStruct(cmd); err != nil {
		return err
	}
	v, err := json.Marshal(cmd)
	if err != nil {
		return err
	}

	if err := p.Produce(ctx, "saga-service", d.CorrelationID, v); err != nil {
		return err
	}

	return nil
}
