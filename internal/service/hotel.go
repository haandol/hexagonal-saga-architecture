package service

import (
	"context"

	"github.com/haandol/hexagonal/internal/adapter/secondary/repository"
	"github.com/haandol/hexagonal/internal/dto"
	"github.com/haandol/hexagonal/internal/instrument"
	"github.com/haandol/hexagonal/internal/message/command"
	"github.com/haandol/hexagonal/internal/port/secondaryport/repositoryport"
	"github.com/haandol/hexagonal/pkg/o11y"
	"github.com/haandol/hexagonal/pkg/util"
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
	logger := util.LoggerFromContext(ctx).WithGroup("HotelService.Book")

	ctx, span := o11y.BeginSubSpan(ctx, "Book")
	defer span.End()

	req := &dto.HotelBooking{
		TripID:  cmd.Body.TripID,
		HotelID: cmd.Body.HotelID,
	}
	if err := s.hotelRepository.Book(ctx, req, cmd); err != nil {
		instrument.RecordBookHotelError(logger, span, err, cmd)

		go func(reason string) {
			if err := s.hotelRepository.PublishAbortSaga(ctx,
				cmd.CorrelationID, cmd.ParentID, cmd.Body.TripID, reason,
			); err != nil {
				instrument.RecordPublishAbortSagaError(logger, span, err, cmd)
			}
		}(err.Error())

		return err
	}

	return nil
}

func (s *HotelService) CancelBooking(ctx context.Context, cmd *command.CancelHotelBooking) error {
	logger := util.LoggerFromContext(ctx).WithGroup("HotelService.CancelBooking")

	ctx, span := o11y.BeginSubSpan(ctx, "CancelBooking")
	defer span.End()

	if err := s.hotelRepository.CancelBooking(ctx, cmd); err != nil {
		instrument.RecordCancelHotelBookingError(logger, span, err, cmd)
		return err
	}

	return nil
}
