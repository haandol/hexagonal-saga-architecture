package service

import (
	"context"
	"fmt"

	"github.com/haandol/hexagonal/internal/adapter/secondary/repository"
	"github.com/haandol/hexagonal/internal/dto"
	"github.com/haandol/hexagonal/internal/instrument"
	"github.com/haandol/hexagonal/internal/message/command"
	"github.com/haandol/hexagonal/internal/port/secondaryport/repositoryport"
	"github.com/haandol/hexagonal/pkg/o11y"
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
	logger := util.LoggerFromContext(ctx).WithGroup("CarService.Book")
	logger.Debug("Book car", "command", cmd)

	ctx, span := o11y.BeginSubSpan(ctx, "Book")
	defer span.End()

	// Transaction Begins
	panicked := true
	txCtx, err := s.carRepository.BeginTx(ctx)
	if err != nil {
		instrument.RecordBeginTxError(logger, span, err)
		return err
	}
	defer func() {
		if r := recover(); r != nil || panicked {
			if err := s.carRepository.RollbackTx(txCtx); err != nil {
				instrument.RecordRollbackTxError(logger, span, err)
			}
		}
	}()

	req := &dto.CarBooking{
		TripID: cmd.Body.TripID,
		CarID:  cmd.Body.CarID,
	}
	booking, err := s.carRepository.Book(txCtx, req)
	if err != nil {
		instrument.RecordBookCarError(logger, span, err, req)

		go func(reason string) {
			if err := s.carRepository.PublishAbortSaga(txCtx, cmd, reason); err != nil {
				instrument.RecordPublishAbortSagaError(logger, span, err, cmd)
			}
		}(err.Error())

		return err
	}

	if err := s.carRepository.PublishCarBooked(ctx, cmd.CorrelationID, cmd.ParentID, &booking); err != nil {
		instrument.RecordPublishCarBookedError(logger, span, err, cmd)
		return err
	}

	if err := s.carRepository.CommitTx(txCtx); err == nil {
		panicked = false
	} else {
		instrument.RecordCommitTxError(logger, span, err)
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
	logger := util.LoggerFromContext(ctx).WithGroup("CarService.CancelBooking")

	ctx, span := o11y.BeginSubSpan(ctx, "CancelBooking")
	defer span.End()

	if err := s.carRepository.CancelBooking(ctx, cmd); err != nil {
		instrument.RecordCancelCarBookingError(logger, span, err, cmd)
		return err
	}

	return nil
}
