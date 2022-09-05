package service

import (
	"context"

	"github.com/haandol/hexagonal/pkg/dto"
	"github.com/haandol/hexagonal/pkg/message/command"
	"github.com/haandol/hexagonal/pkg/port/secondaryport/producerport"
	"github.com/haandol/hexagonal/pkg/port/secondaryport/repositoryport"
	"github.com/haandol/hexagonal/pkg/util"
)

type HotelService struct {
	hotelProducer   producerport.HotelProducer
	hotelRepository repositoryport.HotelRepository
}

func NewHotelService(
	hotelProducer producerport.HotelProducer,
	hotelRepository repositoryport.HotelRepository,
) *HotelService {
	return &HotelService{
		hotelProducer:   hotelProducer,
		hotelRepository: hotelRepository,
	}
}

func (s *HotelService) Book(ctx context.Context, cmd *command.BookHotel) error {
	logger := util.GetLogger().With(
		"service", "HotelService",
		"method", "Book",
	)

	logger.Infow("book hotel", "command", cmd)

	req := &dto.HotelBooking{
		TripID:  cmd.Body.TripID,
		HotelID: cmd.Body.HotelID,
	}
	booking, err := s.hotelRepository.Book(ctx, req)
	if err != nil {
		logger.Errorf("failed to book hotel", "req", req, "err", err.Error())
		return err
	}

	if err := s.hotelProducer.PublishHotelBooked(ctx, cmd.CorrelationID, booking); err != nil {
		return err
	}

	return nil
}

func (s *HotelService) CancelBooking(ctx context.Context, cmd *command.CancelHotelBooking) error {
	logger := util.GetLogger().With(
		"service", "HotelService",
		"method", "CancelBooking",
	)

	logger.Infow("cancel hotel booking", "command", cmd)

	booking, err := s.hotelRepository.CancelBooking(ctx, cmd.Body.BookingID)
	if err != nil {
		logger.Errorf("failed to cancel hotel booking", "BookingID", cmd.Body.BookingID, "err", err.Error())
		return err
	}

	if err := s.hotelProducer.PublishHotelBookingCanceled(ctx, cmd.CorrelationID, booking); err != nil {
		return err
	}

	return nil
}
