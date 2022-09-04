package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/haandol/hexagonal/pkg/dto"
	"github.com/haandol/hexagonal/pkg/message"
	"github.com/haandol/hexagonal/pkg/message/command"
	"github.com/haandol/hexagonal/pkg/message/event"
	"github.com/haandol/hexagonal/pkg/port/secondaryport/producerport"
	"github.com/haandol/hexagonal/pkg/port/secondaryport/repositoryport"
	"github.com/haandol/hexagonal/pkg/util"
)

type HotelService struct {
	producer        producerport.Producer
	hotelRepository repositoryport.HotelRepository
}

func NewHotelService(
	producer producerport.Producer,
	hotelRepository repositoryport.HotelRepository,
) *HotelService {
	return &HotelService{
		producer:        producer,
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
	boooking, err := s.hotelRepository.Book(ctx, req)
	if err != nil {
		logger.Errorf("failed to book hotel", "req", req, "err", err.Error())
		return err
	}

	evt := &event.HotelBooked{
		Message: message.Message{
			Name:          "HotelBooked",
			Version:       "1.0.0",
			ID:            uuid.NewString(),
			CorrelationID: cmd.CorrelationID,
			CreatedAt:     time.Now().Format(time.RFC3339),
		},
		Body: event.HotelBookedBody{
			BookingID: boooking.ID,
		},
	}
	v, err := json.Marshal(evt)
	if err != nil {
		logger.Errorf("failed to marshal event", "event", evt, "err", err.Error())
		return err
	}

	if err := s.producer.Produce(ctx, "saga-service", cmd.CorrelationID, v); err != nil {
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

	if err := s.hotelRepository.CancelBooking(ctx, cmd.Body.BookingID); err != nil {
		logger.Errorf("failed to cancel hotel booking", "BookingID", cmd.Body.BookingID, "err", err.Error())
		return err
	}

	evt := &event.HotelBookingCanceled{
		Message: message.Message{
			Name:          "HotelBookingCanceled",
			Version:       "1.0.0",
			ID:            uuid.NewString(),
			CorrelationID: cmd.CorrelationID,
			CreatedAt:     time.Now().Format(time.RFC3339),
		},
		Body: event.HotelBookedBody{
			BookingID: cmd.Body.BookingID,
		},
	}
	v, err := json.Marshal(evt)
	if err != nil {
		logger.Errorf("failed to marshal event", "event", evt, "err", err.Error())
		return err
	}

	if err := s.producer.Produce(ctx, "saga-service", cmd.CorrelationID, v); err != nil {
		return err
	}

	return nil
}
