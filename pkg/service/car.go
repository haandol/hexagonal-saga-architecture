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

type CarService struct {
	producer      producerport.Producer
	carRepository repositoryport.CarRepository
}

func NewCarService(
	producer producerport.Producer,
	carRepository repositoryport.CarRepository,
) *CarService {
	return &CarService{
		producer:      producer,
		carRepository: carRepository,
	}
}

func (s *CarService) Book(ctx context.Context, cmd *command.BookCar) error {
	logger := util.GetLogger().With(
		"service", "CarService",
		"method", "Book",
	)

	logger.Infow("Book car", "command", cmd)

	req := &dto.CarBooking{
		TripID: cmd.Body.TripID,
		CarID:  cmd.Body.CarID,
	}
	booking, err := s.carRepository.Book(ctx, req)
	if err != nil {
		logger.Errorf("failed to book car", "req", req, "err", err.Error())
		return err
	}

	evt := &event.CarBooked{
		Message: message.Message{
			Name:          "CarBooked",
			Version:       "1.0.0",
			ID:            uuid.NewString(),
			CorrelationID: cmd.CorrelationID,
			CreatedAt:     time.Now().Format(time.RFC3339),
		},
		Body: event.CarBookedBody{
			BookingID: booking.ID,
		},
	}
	v, err := json.Marshal(evt)
	if err != nil {
		logger.Errorf("failed to marshal event", "event", evt, "err", err.Error())
		return err
	}

	if err := s.producer.Produce(ctx, "saga-service", cmd.CorrelationID, v); err != nil {
		return err
	}

	return nil
}

func (s *CarService) CancelBooking(ctx context.Context, cmd *command.CancelCarBooking) error {
	logger := util.GetLogger().With(
		"service", "CarService",
		"method", "CancelBooking",
	)

	logger.Infow("cancel car booking", "command", cmd)

	if err := s.carRepository.CancelBooking(ctx, cmd.Body.BookingID); err != nil {
		logger.Errorf("failed to cancel car booking", "BookingID", cmd.Body.BookingID, "err", err.Error())
		return err
	}

	evt := &event.CarBookingCanceled{
		Message: message.Message{
			Name:          "CarBookingCanceled",
			Version:       "1.0.0",
			ID:            uuid.NewString(),
			CorrelationID: cmd.CorrelationID,
			CreatedAt:     time.Now().Format(time.RFC3339),
		},
		Body: event.CarBookedBody{
			BookingID: cmd.Body.BookingID,
		},
	}
	v, err := json.Marshal(evt)
	if err != nil {
		logger.Errorf("failed to marshal event", "event", evt, "err", err.Error())
		return err
	}

	if err := s.producer.Produce(ctx, "saga-service", cmd.CorrelationID, v); err != nil {
		return err
	}

	return nil
}
