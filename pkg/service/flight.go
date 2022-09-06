package service

import (
	"context"

	"github.com/haandol/hexagonal/pkg/dto"
	"github.com/haandol/hexagonal/pkg/message/command"
	"github.com/haandol/hexagonal/pkg/port/secondaryport/producerport"
	"github.com/haandol/hexagonal/pkg/port/secondaryport/repositoryport"
	"github.com/haandol/hexagonal/pkg/util"
)

type FlightService struct {
	flightProducer   producerport.FlightProducer
	flightRepository repositoryport.FlightRepository
}

func NewFlightService(
	flightProducer producerport.FlightProducer,
	flightRepository repositoryport.FlightRepository,
) *FlightService {
	return &FlightService{
		flightProducer:   flightProducer,
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
	booking, err := s.flightRepository.Book(ctx, req)
	if err != nil {
		logger.Errorw("failed to book flight", "req", req, "err", err.Error())
		return err
	}

	if err := s.flightProducer.PublishFlightBooked(ctx, cmd.CorrelationID, booking); err != nil {
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

	booking, err := s.flightRepository.CancelBooking(ctx, cmd.Body.BookingID)
	if err != nil {
		logger.Errorw("failed to cancel flight booking", "BookingID", cmd.Body.BookingID, "err", err.Error())

		go func(reason string) {
			if err := s.flightProducer.PublishAbortSaga(ctx, cmd.CorrelationID, cmd.Body.TripID, reason); err != nil {
				logger.Errorw("failed to publish abort saga", "command", cmd, "err", err.Error())
			}
		}(err.Error())

		return err
	}

	if err := s.flightProducer.PublishFlightBookingCancelled(ctx, cmd.CorrelationID, booking); err != nil {
		return err
	}

	return nil
}
