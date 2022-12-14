package service

import (
	"context"

	"github.com/haandol/hexagonal/pkg/adapter/secondary/repository"
	"github.com/haandol/hexagonal/pkg/dto"
	"github.com/haandol/hexagonal/pkg/message/command"
	"github.com/haandol/hexagonal/pkg/port/secondaryport/repositoryport"
	"github.com/haandol/hexagonal/pkg/util"
	"github.com/haandol/hexagonal/pkg/util/o11y"
)

type FlightService struct {
	flightRepository repositoryport.FlightRepository
}

func NewFlightService(
	flightRepository *repository.FlightRepository,
) *FlightService {
	return &FlightService{
		flightRepository: flightRepository,
	}
}

func (s *FlightService) Book(ctx context.Context, cmd *command.BookFlight) error {
	logger := util.GetLogger().WithContext(ctx).With(
		"service", "FlightService",
		"method", "Book",
		"command", cmd,
	)

	ctx, span := o11y.BeginSubSpan(ctx, "Book")
	defer span.End()

	req := &dto.FlightBooking{
		TripID:   cmd.Body.TripID,
		FlightID: cmd.Body.FlightID,
	}
	if err := s.flightRepository.Book(ctx, req, cmd); err != nil {
		logger.Errorw("failed to book flight", "req", req, "err", err.Error())

		go func(reason string) {
			if err := s.flightRepository.PublishAbortSaga(ctx,
				cmd.CorrelationID, cmd.ParentID, cmd.Body.TripID, reason,
			); err != nil {
				logger.Errorw("failed to publish abort saga", "command", cmd, "err", err.Error())
			}
		}(err.Error())

		span.RecordError(err)
		span.SetStatus(o11y.GetStatus(err))
		return err
	}

	return nil
}

func (s *FlightService) CancelBooking(ctx context.Context, cmd *command.CancelFlightBooking) error {
	logger := util.GetLogger().WithContext(ctx).With(
		"service", "FlightService",
		"method", "CancelBooking",
		"command", cmd,
	)

	ctx, span := o11y.BeginSubSpan(ctx, "CancelBooking")
	defer span.End()

	if err := s.flightRepository.CancelBooking(ctx, cmd); err != nil {
		logger.Errorw("failed to cancel flight booking", "BookingID", cmd.Body.BookingID, "err", err.Error())
		return err
	}

	return nil
}
