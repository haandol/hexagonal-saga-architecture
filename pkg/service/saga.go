package service

import (
	"context"
	"encoding/json"
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

	if err := s.publishRentCar(ctx, saga); err != nil {
		logger.Errorw("failed to publish rent car", "saga", saga, "error", err.Error())
		return err
	}

	if err := s.publishBookHotel(ctx, saga); err != nil {
		logger.Errorw("failed to publish book hotel", "saga", saga, "error", err.Error())
		return err
	}

	if err := s.publishBookFlight(ctx, saga); err != nil {
		logger.Errorw("failed to publish book flight", "saga", saga, "error", err.Error())
		return err
	}

	return nil
}

func (s *SagaService) publishRentCar(ctx context.Context, d dto.Saga) error {
	cmd := &command.RentCar{
		Message: message.Message{
			Name:          "RentCar",
			Version:       "1.0.0",
			ID:            uuid.NewString(),
			CorrelationID: d.CorrelationID,
			CreatedAt:     time.Now().Format(time.RFC3339),
		},
		Body: command.RentCarBody{
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

func (s *SagaService) ProcessCarRental(ctx context.Context, evt *event.CarRented) error {
	logger := util.GetLogger().With(
		"module", "SagaService",
		"method", "ProcessCarRental",
	)

	logger.Infow("success car rented", "event", evt)

	if err := s.sagaRepository.ProcessCarRental(ctx, evt); err != nil {
		logger.Errorf("failed to process car rented", "event", evt, "err", err.Error())
		return err
	}

	return nil
}

func (s *SagaService) CompensateCarRental(ctx context.Context, evt *event.CarRentalCanceled) error {
	logger := util.GetLogger().With(
		"module", "SagaService",
		"method", "CompensateCarRental",
	)

	logger.Infow("cancel car rental", "event", evt)

	if err := s.sagaRepository.CompensateCarRental(ctx, evt); err != nil {
		logger.Errorf("failed to process cancel car rental", "event", evt, "err", err.Error())
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
