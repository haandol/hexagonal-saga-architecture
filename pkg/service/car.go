package service

import (
	"context"

	"github.com/haandol/hexagonal/pkg/dto"
	"github.com/haandol/hexagonal/pkg/message/command"
	"github.com/haandol/hexagonal/pkg/port/secondaryport/repositoryport"
	"github.com/haandol/hexagonal/pkg/util"
)

type CarService struct {
	carRepository repositoryport.CarRepository
}

func NewCarService(
	carRepository repositoryport.CarRepository,
) *CarService {
	return &CarService{
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
	if err := s.carRepository.Book(ctx, req, cmd); err != nil {
		logger.Errorw("failed to book car", "req", req, "err", err.Error())

		go func(reason string) {
			if err := s.carRepository.PublishAbortSaga(ctx,
				cmd.CorrelationID, cmd.ParentID, cmd.Body.TripID, reason,
			); err != nil {
				logger.Errorw("failed to publish abort saga", "command", cmd, "err", err.Error())
			}
		}(err.Error())

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

	if err := s.carRepository.CancelBooking(ctx, cmd); err != nil {
		logger.Errorw("failed to cancel car booking", "BookingID", cmd.Body.BookingID, "err", err.Error())
		return err
	}

	return nil
}
