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

type HotelService struct {
	hotelRepository repositoryport.HotelRepository
}

func NewHotelService(
	hotelRepository *repository.HotelRepository,
) *HotelService {
	return &HotelService{
		hotelRepository: hotelRepository,
	}
}

func (s *HotelService) Book(ctx context.Context, cmd *command.BookHotel) error {
	logger := util.GetLogger().WithContext(ctx).With(
		"service", "HotelService",
		"method", "Book",
		"command", cmd,
	)

	ctx, span := o11y.BeginSubSpan(ctx, "Book")
	defer span.End()

	req := &dto.HotelBooking{
		TripID:  cmd.Body.TripID,
		HotelID: cmd.Body.HotelID,
	}
	if err := s.hotelRepository.Book(ctx, req, cmd); err != nil {
		logger.Errorw("failed to book hotel", "req", req, "err", err)

		go func(reason string) {
			if err := s.hotelRepository.PublishAbortSaga(ctx,
				cmd.CorrelationID, cmd.ParentID, cmd.Body.TripID, reason,
			); err != nil {
				logger.Errorw("failed to publish abort saga", "command", cmd, "err", err)
			}
		}(err.Error())

		span.RecordError(err)
		span.SetStatus(o11y.GetStatus(err))
		return err
	}

	return nil
}

func (s *HotelService) CancelBooking(ctx context.Context, cmd *command.CancelHotelBooking) error {
	logger := util.GetLogger().WithContext(ctx).With(
		"service", "HotelService",
		"method", "CancelBooking",
		"command", cmd,
	)

	ctx, span := o11y.BeginSubSpan(ctx, "CancelBooking")
	defer span.End()

	if err := s.hotelRepository.CancelBooking(ctx, cmd); err != nil {
		logger.Errorw("failed to cancel hotel booking", "BookingID", cmd.Body.BookingID, "err", err)
		return err
	}

	return nil
}
