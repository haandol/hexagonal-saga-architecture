package service

import (
	"context"

	"github.com/haandol/hexagonal/pkg/adapter/secondary/repository"
	"github.com/haandol/hexagonal/pkg/dto"
	"github.com/haandol/hexagonal/pkg/message/command"
	"github.com/haandol/hexagonal/pkg/port/secondaryport/repositoryport"
	"github.com/haandol/hexagonal/pkg/util"
)

type CarService struct {
	carRepository repositoryport.CarRepository
}

func NewCarService(
	carRepository *repository.CarRepository,
) *CarService {
	return &CarService{
		carRepository: carRepository,
	}
}

func (s *CarService) Book(ctx context.Context, cmd *command.BookCar) error {
	logger := util.GetLogger().WithContext(ctx).With(
		"service", "CarService",
		"method", "Book",
	)
	logger.Debugw("Book car", "command", cmd)

	// Transaction Begins
	panicked := true
	txCtx, err := s.carRepository.BeginTx(ctx)
	if err != nil {
		logger.Errorw("Failed to begin transaction", "err", err.Error())
		return err
	}
	defer func() {
		if r := recover(); r != nil || panicked {
			if err := s.carRepository.RollbackTx(txCtx); err != nil {
				logger.Errorw("Failed to rollback transaction", "err", err.Error())
			}
		}
	}()

	req := &dto.CarBooking{
		TripID: cmd.Body.TripID,
		CarID:  cmd.Body.CarID,
	}
	booking, err := s.carRepository.Book(txCtx, req)
	if err != nil {
		logger.Errorw("Failed to book car", "req", req, "err", err.Error())

		if err := s.carRepository.PublishAbortSaga(txCtx, cmd, err.Error()); err != nil {
			logger.Errorw("Failed to publish AbortSaga", "command", cmd, "err", err.Error())
			return err
		}
	} else {
		if err := s.carRepository.PublishCarBooked(ctx, cmd.CorrelationID, cmd.ParentID, booking); err != nil {
			logger.Errorw("Failed to publish CarBooked", "booking", booking, "err", err.Error())
			return err
		}
	}

	if err := s.carRepository.CommitTx(txCtx); err == nil {
		panicked = false
	} else {
		logger.Errorw("Failed to commit transaction", "err", err.Error())
		return err
	}
	// Transaction Ends

	return nil
}

func (s *CarService) CancelBooking(ctx context.Context, cmd *command.CancelCarBooking) error {
	logger := util.GetLogger().WithContext(ctx).With(
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
