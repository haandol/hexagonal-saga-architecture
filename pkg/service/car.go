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

func (s *CarService) Rent(ctx context.Context, cmd *command.RentCar) error {
	logger := util.GetLogger().With(
		"service", "CarService",
		"method", "Rent",
	)

	logger.Infow("rent car", "command", cmd)

	req := &dto.CarRental{
		TripID:   cmd.Body.TripID,
		CarID:    cmd.Body.CarID,
		Quantity: 1,
	}
	rental, err := s.carRepository.Rent(ctx, req)
	if err != nil {
		logger.Errorf("failed to rent car", "req", req, "err", err.Error())
		return err
	}

	evt := &event.CarRented{
		Message: message.Message{
			Name:          "CarRented",
			Version:       "1.0.0",
			ID:            uuid.NewString(),
			CorrelationID: cmd.CorrelationID,
			CreatedAt:     time.Now().Format(time.RFC3339),
		},
		Body: event.CarRentedBody{
			CarRentalID: rental.ID,
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

func (s *CarService) CancelRental(ctx context.Context, cmd *command.CancelCarRental) error {
	logger := util.GetLogger().With(
		"service", "CarService",
		"method", "CancelRental",
	)

	logger.Infow("cancel car rental", "command", cmd)

	if err := s.carRepository.CancelRental(ctx, cmd.Body.RentalID); err != nil {
		logger.Errorf("failed to cancel car rental", "rentalID", cmd.Body.RentalID, "err", err.Error())
		return err
	}

	evt := &event.CarRentalCanceled{
		Message: message.Message{
			Name:          "CarRentalCanceled",
			Version:       "1.0.0",
			ID:            uuid.NewString(),
			CorrelationID: cmd.CorrelationID,
			CreatedAt:     time.Now().Format(time.RFC3339),
		},
		Body: event.CarRentedBody{
			CarRentalID: cmd.Body.RentalID,
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
