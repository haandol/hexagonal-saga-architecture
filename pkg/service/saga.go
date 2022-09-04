package service

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/haandol/hexagonal/pkg/dto"
	"github.com/haandol/hexagonal/pkg/message"
	"github.com/haandol/hexagonal/pkg/message/command"
	"github.com/haandol/hexagonal/pkg/message/event"
	"github.com/haandol/hexagonal/pkg/port/secondaryport/producerport"
	"github.com/haandol/hexagonal/pkg/port/secondaryport/repositoryport"
	"github.com/haandol/hexagonal/pkg/util"
)

type SagaService struct {
	producer       producerport.Producer
	sagaRepository repositoryport.SagaRepository
}

func NewSagaService(
	producer producerport.Producer,
	sagaRepository repositoryport.SagaRepository,
) *SagaService {
	return &SagaService{
		producer:       producer,
		sagaRepository: sagaRepository,
	}
}

func (s *SagaService) Start(ctx context.Context, cmd *command.StartSaga) error {
	logger := util.GetLogger().With(
		"module", "SagaService",
		"method", "Start",
	)

	logger.Infow("start saga", "command", cmd)

	saga, err := s.sagaRepository.Start(ctx, cmd)
	if err != nil {
		logger.Errorw("failed to create saga", "command", cmd, "err", err.Error())
	}

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		if err := s.publishBookCar(ctx, saga); err != nil {
			logger.Errorw("failed to publish book car", "saga", saga, "error", err.Error())
		}
		wg.Done()
	}()

	go func() {
		if err := s.publishBookHotel(ctx, saga); err != nil {
			logger.Errorw("failed to publish book hotel", "saga", saga, "error", err.Error())
		}
		wg.Done()
	}()

	go func() {
		if err := s.publishBookFlight(ctx, saga); err != nil {
			logger.Errorw("failed to publish book flight", "saga", saga, "error", err.Error())
		}
		wg.Done()
	}()

	wg.Wait()

	return nil
}

func (s *SagaService) publishBookCar(ctx context.Context, d dto.Saga) error {
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
	v, err := json.Marshal(cmd)
	if err != nil {
		return err
	}

	if err := s.producer.Produce(ctx, "car-service", d.CorrelationID, v); err != nil {
		return err
	}

	return nil
}

func (s *SagaService) publishBookHotel(ctx context.Context, d dto.Saga) error {
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
	v, err := json.Marshal(cmd)
	if err != nil {
		return err
	}

	if err := s.producer.Produce(ctx, "hotel-service", d.CorrelationID, v); err != nil {
		return err
	}

	return nil
}

func (s *SagaService) publishBookFlight(ctx context.Context, d dto.Saga) error {
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
	v, err := json.Marshal(cmd)
	if err != nil {
		return err
	}

	if err := s.producer.Produce(ctx, "flight-service", d.CorrelationID, v); err != nil {
		return err
	}

	return nil
}

func (s *SagaService) ProcessCarBooking(ctx context.Context, evt *event.CarBooked) error {
	logger := util.GetLogger().With(
		"module", "SagaService",
		"method", "ProcessCarBooking",
	)

	logger.Infow("success car booked", "event", evt)

	if err := s.sagaRepository.ProcessCarBooking(ctx, evt); err != nil {
		logger.Errorf("failed to process car booked", "event", evt, "err", err.Error())
		return err
	}

	return nil
}

func (s *SagaService) CompensateCarBooking(ctx context.Context, evt *event.CarBookingCanceled) error {
	logger := util.GetLogger().With(
		"module", "SagaService",
		"method", "CompensateCarBooking",
	)

	logger.Infow("cancel car booking", "event", evt)

	if err := s.sagaRepository.CompensateCarBooking(ctx, evt); err != nil {
		logger.Errorf("failed to process cancel car booking", "event", evt, "err", err.Error())
		return err
	}

	return nil
}

func (s *SagaService) ProcessHotelBooking(ctx context.Context, evt *event.HotelBooked) error {
	logger := util.GetLogger().With(
		"module", "SagaService",
		"method", "ProcessHotelBooking",
	)

	logger.Infow("success hotel booked", "event", evt)

	if err := s.sagaRepository.ProcessHotelBooking(ctx, evt); err != nil {
		logger.Errorf("failed to process Hotel booked", "event", evt, "err", err.Error())
		return err
	}

	return nil
}

func (s *SagaService) CompensateHotelBooking(ctx context.Context, evt *event.HotelBookingCanceled) error {
	logger := util.GetLogger().With(
		"module", "SagaService",
		"method", "CompensateHotelBooking",
	)

	logger.Infow("cancel hotel booking", "event", evt)

	if err := s.sagaRepository.CompensateHotelBooking(ctx, evt); err != nil {
		logger.Errorf("failed to process cancel Hotel booking", "event", evt, "err", err.Error())
		return err
	}

	return nil
}

func (s *SagaService) ProcessFlightBooking(ctx context.Context, evt *event.FlightBooked) error {
	logger := util.GetLogger().With(
		"module", "SagaService",
		"method", "ProcessFlightBooking",
	)

	logger.Infow("success flight booked", "event", evt)

	if err := s.sagaRepository.ProcessFlightBooking(ctx, evt); err != nil {
		logger.Errorf("failed to process flight booked", "event", evt, "err", err.Error())
		return err
	}

	return nil
}

func (s *SagaService) CompensateFlightBooking(ctx context.Context, evt *event.FlightBookingCanceled) error {
	logger := util.GetLogger().With(
		"module", "SagaService",
		"method", "CompensateFlightBooking",
	)

	logger.Infow("cancel flight booking", "event", evt)

	if err := s.sagaRepository.CompensateFlightBooking(ctx, evt); err != nil {
		logger.Errorf("failed to process cancel flight booking", "event", evt, "err", err.Error())
		return err
	}

	return nil
}

func (s *SagaService) End(ctx context.Context, cmd *command.EndSaga) error {
	logger := util.GetLogger().With(
		"module", "SagaService",
		"method", "End",
	)

	logger.Infow("end saga", "command", cmd)

	if err := s.sagaRepository.End(ctx, cmd); err != nil {
		logger.Errorw("failed to end saga", "command", cmd, "err", err.Error())
	}

	event := &event.SagaEnded{
		Message: message.Message{
			Name: "SagaEnded",
		},
		Body: event.SagaEndedBody{
			SagaID: cmd.Body.SagaID,
		},
	}
	v, err := json.Marshal(event)
	if err != nil {
		logger.Errorw("failed to marshal saga ended event", "event", event, "err", err.Error())
	}

	if err := s.producer.Produce(ctx, "saga-service", cmd.CorrelationID, v); err != nil {
		return err
	}

	return nil
}

func (s *SagaService) Abort(ctx context.Context, cmd *command.AbortSaga) error {
	logger := util.GetLogger().With(
		"module", "SagaService",
		"method", "Abort",
	)

	logger.Infow("abort saga", "command", cmd)

	if err := s.sagaRepository.Abort(ctx, cmd); err != nil {
		logger.Errorw("failed to abort saga", "command", cmd, "err", err.Error())
	}

	return nil
}
