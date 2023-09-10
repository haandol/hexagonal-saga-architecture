package producer

import (
	"context"
	"encoding/json"
	"reflect"
	"time"

	"github.com/google/uuid"
	"github.com/haandol/hexagonal/pkg/connector/producer"
	"github.com/haandol/hexagonal/pkg/dto"
	"github.com/haandol/hexagonal/pkg/message"
	"github.com/haandol/hexagonal/pkg/message/command"
	"github.com/haandol/hexagonal/pkg/message/event"
	"github.com/haandol/hexagonal/pkg/util"
	"github.com/haandol/hexagonal/pkg/util/o11y"
)

type SagaProducer struct {
	*producer.KafkaProducer
}

func NewSagaProducer(kafkaProducer *producer.KafkaProducer) *SagaProducer {
	return &SagaProducer{
		KafkaProducer: kafkaProducer,
	}
}

func (p *SagaProducer) PublishBookCar(ctx context.Context, d *dto.Saga) error {
	traceID, spanID := o11y.GetTraceSpanID(ctx)
	cmd := &command.BookCar{
		Message: message.Message{
			Name:          reflect.ValueOf(command.BookCar{}).Type().Name(),
			Version:       "1.0.0",
			ID:            uuid.NewString(),
			CorrelationID: traceID,
			ParentID:      spanID,
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

func (p *SagaProducer) PublishCancelCarBooking(ctx context.Context, d *dto.Saga) error {
	traceID, spanID := o11y.GetTraceSpanID(ctx)
	cmd := &command.CancelCarBooking{
		Message: message.Message{
			Name:          reflect.ValueOf(command.CancelCarBooking{}).Type().Name(),
			Version:       "1.0.0",
			ID:            uuid.NewString(),
			CorrelationID: traceID,
			ParentID:      spanID,
			CreatedAt:     time.Now().Format(time.RFC3339),
		},
		Body: command.CancelCarBookingBody{
			TripID:    d.TripID,
			BookingID: d.CarBookingID,
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

func (p *SagaProducer) PublishBookHotel(ctx context.Context, d *dto.Saga) error {
	traceID, spanID := o11y.GetTraceSpanID(ctx)
	cmd := &command.BookHotel{
		Message: message.Message{
			Name:          reflect.ValueOf(command.BookHotel{}).Type().Name(),
			Version:       "1.0.0",
			ID:            uuid.NewString(),
			CorrelationID: traceID,
			ParentID:      spanID,
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

func (p *SagaProducer) PublishCancelHotelBooking(ctx context.Context, d *dto.Saga) error {
	traceID, spanID := o11y.GetTraceSpanID(ctx)
	cmd := &command.CancelHotelBooking{
		Message: message.Message{
			Name:          reflect.ValueOf(command.CancelHotelBooking{}).Type().Name(),
			Version:       "1.0.0",
			ID:            uuid.NewString(),
			CorrelationID: traceID,
			ParentID:      spanID,
			CreatedAt:     time.Now().Format(time.RFC3339),
		},
		Body: command.CancelHotelBookingBody{
			TripID:    d.TripID,
			BookingID: d.HotelBookingID,
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

func (p *SagaProducer) PublishBookFlight(ctx context.Context, d *dto.Saga) error {
	traceID, spanID := o11y.GetTraceSpanID(ctx)
	cmd := &command.BookFlight{
		Message: message.Message{
			Name:          reflect.ValueOf(command.BookFlight{}).Type().Name(),
			Version:       "1.0.0",
			ID:            uuid.NewString(),
			CorrelationID: traceID,
			ParentID:      spanID,
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

func (p *SagaProducer) PublishCancelFlightBooking(ctx context.Context, d *dto.Saga) error {
	traceID, spanID := o11y.GetTraceSpanID(ctx)
	cmd := &command.CancelFlightBooking{
		Message: message.Message{
			Name:          reflect.ValueOf(command.CancelFlightBooking{}).Type().Name(),
			Version:       "1.0.0",
			ID:            uuid.NewString(),
			CorrelationID: traceID,
			ParentID:      spanID,
			CreatedAt:     time.Now().Format(time.RFC3339),
		},
		Body: command.CancelFlightBookingBody{
			TripID:    d.TripID,
			BookingID: d.FlightBookingID,
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

func (p *SagaProducer) PublishEndSaga(ctx context.Context, d *dto.Saga) error {
	traceID, spanID := o11y.GetTraceSpanID(ctx)
	cmd := &command.EndSaga{
		Message: message.Message{
			Name:          reflect.ValueOf(command.EndSaga{}).Type().Name(),
			Version:       "1.0.0",
			ID:            uuid.NewString(),
			CorrelationID: traceID,
			ParentID:      spanID,
			CreatedAt:     time.Now().Format(time.RFC3339),
		},
		Body: command.EndSagaBody{
			SagaID: d.ID,
			TripID: d.TripID,
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

func (p *SagaProducer) PublishSagaEnded(ctx context.Context,
	corrID string, parentID string, d *dto.Saga,
) error {
	logger := util.LoggerFromContext(ctx).With(
		"module", "Publisher",
		"func", "PublishSagaEnded",
	)

	evt := &event.SagaEnded{
		Message: message.Message{
			Name:          reflect.ValueOf(event.SagaEnded{}).Type().Name(),
			Version:       "1.0.0",
			ID:            uuid.NewString(),
			CorrelationID: corrID,
			ParentID:      parentID,
			CreatedAt:     time.Now().Format(time.RFC3339),
		},
		Body: event.SagaEndedBody{
			SagaID:          d.ID,
			TripID:          d.TripID,
			CarBookingID:    d.CarBookingID,
			HotelBookingID:  d.HotelBookingID,
			FlightBookingID: d.FlightBookingID,
		},
	}
	if err := util.ValidateStruct(evt); err != nil {
		return err
	}
	v, err := json.Marshal(evt)
	if err != nil {
		logger.Error("failed to marshal saga aborted event", "event", evt, "err", err)
	}

	if err := p.Produce(ctx, "trip-service", corrID, v); err != nil {
		logger.Error("failed to produce saga ended event", "event", evt, "err", err)
		return err
	}

	return nil
}

func (p *SagaProducer) PublishSagaAborted(ctx context.Context,
	corrID string, parentID string, d *dto.Saga,
) error {
	logger := util.LoggerFromContext(ctx).With(
		"module", "Publisher",
		"func", "PublishSagaAborted",
	)

	evt := &event.SagaAborted{
		Message: message.Message{
			Name:          reflect.ValueOf(event.SagaAborted{}).Type().Name(),
			Version:       "1.0.0",
			ID:            uuid.NewString(),
			CorrelationID: corrID,
			ParentID:      parentID,
			CreatedAt:     time.Now().Format(time.RFC3339),
		},
		Body: event.SagaAbortedBody{
			SagaID: d.ID,
			TripID: d.TripID,
		},
	}
	if err := util.ValidateStruct(evt); err != nil {
		return err
	}
	v, err := json.Marshal(evt)
	if err != nil {
		logger.Error("failed to marshal saga aborted event", "event", evt, "err", err)
	}

	if err := p.Produce(ctx, "trip-service", corrID, v); err != nil {
		logger.Error("failed to produce saga aborted event", "event", evt, "err", err)
	}

	return nil
}
