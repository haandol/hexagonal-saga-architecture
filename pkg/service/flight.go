package service

import (
	"context"

	"github.com/haandol/hexagonal/pkg/dto"
	"github.com/haandol/hexagonal/pkg/message/command"
	"github.com/haandol/hexagonal/pkg/port/secondaryport/repositoryport"
	"github.com/haandol/hexagonal/pkg/util"
)

type FlightService struct {
	flightRepository repositoryport.FlightRepository
}

func NewFlightService(
	flightRepository repositoryport.FlightRepository,
) *FlightService {
	return &FlightService{
		flightRepository: flightRepository,
	}
}

func (s *FlightService) Book(ctx context.Context, cmd *command.BookFlight) error {
	logger := util.GetLogger().With(
		"service", "FlightService",
		"method", "Book",
	)

	logger.Infow("book flight", "command", cmd)

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

		return err
	}

	return nil
}

func (s *FlightService) CancelBooking(ctx context.Context, cmd *command.CancelFlightBooking) error {
	logger := util.GetLogger().With(
		"service", "FlightService",
		"method", "CancelBooking",
	)

	logger.Infow("cancel flight booking", "command", cmd)

	if err := s.flightRepository.CancelBooking(ctx, cmd); err != nil {
		logger.Errorw("failed to cancel flight booking", "BookingID", cmd.Body.BookingID, "err", err.Error())
		return err
	}

	return nil
}
