package service

import (
	"context"
	"fmt"

	"github.com/haandol/hexagonal/pkg/adapter/secondary/repository"
	"github.com/haandol/hexagonal/pkg/dto"
	"github.com/haandol/hexagonal/pkg/message/command"
	"github.com/haandol/hexagonal/pkg/port/secondaryport/repositoryport"
	"github.com/haandol/hexagonal/pkg/util"
	"github.com/haandol/hexagonal/pkg/util/o11y"
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

	ctx, span := o11y.BeginSubSpan(ctx, "Book")
	defer span.End()

	// Transaction Begins
	panicked := true
	txCtx, err := s.carRepository.BeginTx(ctx)
	if err != nil {
		logger.Errorw("Failed to begin transaction", "err", err.Error())
		span.RecordError(err)
		span.SetStatus(o11y.GetStatus(err))
		return err
	}
	defer func() {
		if r := recover(); r != nil || panicked {
			if err := s.carRepository.RollbackTx(txCtx); err != nil {
				span.RecordError(err)
				span.SetStatus(o11y.GetStatus(err))
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
			span.RecordError(err)
			span.SetStatus(o11y.GetStatus(err))
			return err
		}
	} else {
		if err := s.carRepository.PublishCarBooked(ctx, cmd.CorrelationID, cmd.ParentID, &booking); err != nil {
			logger.Errorw("Failed to publish CarBooked", "booking", booking, "err", err.Error())
			span.RecordError(err)
			span.SetStatus(o11y.GetStatus(err))
			return err
		}
	}

	if err := s.carRepository.CommitTx(txCtx); err == nil {
		panicked = false
	} else {
		logger.Errorw("Failed to commit transaction", "err", err.Error())
		span.RecordError(err)
		span.SetStatus(o11y.GetStatus(err))
		return err
	}
	// Transaction Ends

	span.SetAttributes(
		o11y.AttrString("booking", fmt.Sprintf("%v", booking)),
		o11y.AttrString("panicked", fmt.Sprintf("%v", panicked)),
	)

	return nil
}

func (s *CarService) CancelBooking(ctx context.Context, cmd *command.CancelCarBooking) error {
	logger := util.GetLogger().WithContext(ctx).With(
		"service", "CarService",
		"method", "CancelBooking",
		"command", cmd,
	)

	ctx, span := o11y.BeginSubSpan(ctx, "CancelBooking")
	defer span.End()

	if err := s.carRepository.CancelBooking(ctx, cmd); err != nil {
		logger.Errorw("failed to cancel car booking", "BookingID", cmd.Body.BookingID, "err", err.Error())
		span.RecordError(err)
		span.SetStatus(o11y.GetStatus(err))
		return err
	}

	return nil
}
