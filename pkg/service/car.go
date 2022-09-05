package service

import (
	"context"

	"github.com/haandol/hexagonal/pkg/dto"
	"github.com/haandol/hexagonal/pkg/message/command"
	"github.com/haandol/hexagonal/pkg/port/secondaryport/producerport"
	"github.com/haandol/hexagonal/pkg/port/secondaryport/repositoryport"
	"github.com/haandol/hexagonal/pkg/util"
)

type CarService struct {
	carProducer   producerport.CarProducer
	carRepository repositoryport.CarRepository
}

func NewCarService(
	carProducer producerport.CarProducer,
	carRepository repositoryport.CarRepository,
) *CarService {
	return &CarService{
		carProducer:   carProducer,
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
		logger.Errorw("failed to book car", "req", req, "err", err.Error())

		go func(reason string) {
			if err := s.carProducer.PublishAbortSaga(ctx, cmd.CorrelationID, cmd.Body.TripID, reason); err != nil {
				logger.Errorw("failed to publish abort saga", "command", cmd, "err", err.Error())
			}
		}(err.Error())

		return err
	}

	if err := s.carProducer.PublishCarBooked(ctx, cmd.CorrelationID, booking); err != nil {
		logger.Errorw("failed to publish car booked", "booking", booking, "err", err.Error())
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

	booking, err := s.carRepository.CancelBooking(ctx, cmd.Body.BookingID)
	if err != nil {
		logger.Errorw("failed to cancel car booking", "BookingID", cmd.Body.BookingID, "err", err.Error())
		return err
	}

	if err := s.carProducer.PublishCarBookingCanceled(ctx, cmd.CorrelationID, booking); err != nil {
		logger.Errorw("failed to publish car booking canceled", "booking", booking, "err", err.Error())
		return err
	}

	return nil
}
