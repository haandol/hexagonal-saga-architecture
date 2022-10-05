package service

import (
	"context"

	"github.com/haandol/hexagonal/pkg/dto"
	"github.com/haandol/hexagonal/pkg/message/command"
	"github.com/haandol/hexagonal/pkg/port/secondaryport/repositoryport"
	"github.com/haandol/hexagonal/pkg/util"
)

type HotelService struct {
	hotelRepository repositoryport.HotelRepository
}

func NewHotelService(
	hotelRepository repositoryport.HotelRepository,
) *HotelService {
	return &HotelService{
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
	if err := s.hotelRepository.Book(ctx, req, cmd); err != nil {
		logger.Errorw("failed to book hotel", "req", req, "err", err.Error())

		go func(reason string) {
			if err := s.hotelRepository.PublishAbortSaga(ctx,
				cmd.CorrelationID, cmd.ParentID, cmd.Body.TripID, reason,
			); err != nil {
				logger.Errorw("failed to publish abort saga", "command", cmd, "err", err.Error())
			}
		}(err.Error())

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

	if err := s.hotelRepository.CancelBooking(ctx, cmd); err != nil {
		logger.Errorw("failed to cancel hotel booking", "BookingID", cmd.Body.BookingID, "err", err.Error())
		return err
	}

	return nil
}
