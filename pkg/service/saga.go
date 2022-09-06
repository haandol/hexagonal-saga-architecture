package service

import (
	"context"

	"github.com/haandol/hexagonal/pkg/message/command"
	"github.com/haandol/hexagonal/pkg/message/event"
	"github.com/haandol/hexagonal/pkg/port/secondaryport/producerport"
	"github.com/haandol/hexagonal/pkg/port/secondaryport/repositoryport"
	"github.com/haandol/hexagonal/pkg/util"
)

type SagaService struct {
	publisher      producerport.SagaProducer
	sagaRepository repositoryport.SagaRepository
}

func NewSagaService(
	publisher producerport.SagaProducer,
	sagaRepository repositoryport.SagaRepository,
) *SagaService {
	return &SagaService{
		publisher:      publisher,
		sagaRepository: sagaRepository,
	}
}

func (s *SagaService) Start(ctx context.Context, cmd *command.StartSaga) error {
	logger := util.GetLogger().With(
		"module", "SagaService",
		"method", "Start",
	)

	logger.Infow("start saga", "command", cmd)

	saga, err := s.sagaRepository.Start(ctx, cmd)
	if err != nil {
		logger.Errorw("failed to create saga", "command", cmd, "err", err.Error())
	}

	if err := s.publisher.PublishBookCar(ctx, saga); err != nil {
		logger.Errorw("failed to publish book car", "saga", saga, "error", err.Error())
	}

	return nil
}

func (s *SagaService) ProcessCarBooking(ctx context.Context, evt *event.CarBooked) error {
	logger := util.GetLogger().With(
		"module", "SagaService",
		"method", "ProcessCarBooking",
	)

	logger.Infow("success car booked", "event", evt)

	saga, err := s.sagaRepository.ProcessCarBooking(ctx, evt)
	if err != nil {
		logger.Errorw("failed to process car booked", "event", evt, "err", err.Error())
		return err
	}

	if err := s.publisher.PublishBookHotel(ctx, saga); err != nil {
		logger.Errorw("failed to publish book hotel", "saga", saga, "error", err.Error())
	}

	return nil
}

func (s *SagaService) CompensateCarBooking(ctx context.Context, evt *event.CarBookingCancelled) error {
	logger := util.GetLogger().With(
		"module", "SagaService",
		"method", "CompensateCarBooking",
	)

	logger.Infow("cancel car booking", "event", evt)

	saga, err := s.sagaRepository.CompensateCarBooking(ctx, evt)
	if err != nil {
		logger.Errorw("failed to process cancel car booking", "event", evt, "err", err.Error())
		return err
	}

	if err := s.publisher.PublishSagaAborted(ctx, evt.CorrelationID, saga); err != nil {
		logger.Errorw("failed to publish SagaAborted", "event", evt, "err", err.Error())
		return err
	}

	return nil
}

func (s *SagaService) ProcessHotelBooking(ctx context.Context, evt *event.HotelBooked) error {
	logger := util.GetLogger().With(
		"module", "SagaService",
		"method", "ProcessHotelBooking",
	)

	logger.Infow("success hotel booked", "event", evt)

	saga, err := s.sagaRepository.ProcessHotelBooking(ctx, evt)
	if err != nil {
		logger.Errorw("failed to process Hotel booked", "event", evt, "err", err.Error())
		return err
	}

	if err := s.publisher.PublishBookFlight(ctx, saga); err != nil {
		logger.Errorw("failed to publish book flight", "saga", saga, "error", err.Error())
		return err
	}

	return nil
}

func (s *SagaService) CompensateHotelBooking(ctx context.Context, evt *event.HotelBookingCancelled) error {
	logger := util.GetLogger().With(
		"module", "SagaService",
		"method", "CompensateHotelBooking",
	)

	logger.Infow("cancel hotel booking", "event", evt)

	_, err := s.sagaRepository.CompensateHotelBooking(ctx, evt)
	if err != nil {
		logger.Errorw("failed to process cancel Hotel booking", "event", evt, "err", err.Error())
		return err
	}

	return nil
}

func (s *SagaService) ProcessFlightBooking(ctx context.Context, evt *event.FlightBooked) error {
	logger := util.GetLogger().With(
		"module", "SagaService",
		"method", "ProcessFlightBooking",
	)

	logger.Infow("success flight booked", "event", evt)

	saga, err := s.sagaRepository.ProcessFlightBooking(ctx, evt)
	if err != nil {
		logger.Errorw("failed to process flight booked", "event", evt, "err", err.Error())
		return err
	}

	if err := s.publisher.PublishEndSaga(ctx, saga); err != nil {
		logger.Errorw("failed to publish end saga", "saga", saga, "error", err.Error())
		return err
	}

	return nil
}

func (s *SagaService) CompensateFlightBooking(ctx context.Context, evt *event.FlightBookingCancelled) error {
	logger := util.GetLogger().With(
		"module", "SagaService",
		"method", "CompensateFlightBooking",
	)

	logger.Infow("cancel flight booking", "event", evt)

	_, err := s.sagaRepository.CompensateFlightBooking(ctx, evt)
	if err != nil {
		logger.Errorw("failed to process cancel flight booking", "event", evt, "err", err.Error())
		return err
	}

	return nil
}

func (s *SagaService) End(ctx context.Context, cmd *command.EndSaga) error {
	logger := util.GetLogger().With(
		"module", "SagaService",
		"method", "End",
	)

	logger.Infow("end saga", "command", cmd)

	saga, err := s.sagaRepository.End(ctx, cmd)
	if err != nil {
		logger.Errorw("failed to end saga", "command", cmd, "err", err.Error())
	}

	if err := s.publisher.PublishSagaEnded(ctx, cmd.CorrelationID, saga); err != nil {
		logger.Errorw("failed to publish SagaEnded", "command", cmd, "err", err.Error())
		return err
	}

	return nil
}

func (s *SagaService) Abort(ctx context.Context, cmd *command.AbortSaga) error {
	logger := util.GetLogger().With(
		"module", "SagaService",
		"method", "Abort",
	)

	logger.Infow("abort saga", "command", cmd)

	saga, err := s.sagaRepository.Abort(ctx, cmd)
	if err != nil {
		logger.Errorw("failed to abort saga", "command", cmd, "err", err.Error())
		return err
	}

	switch cmd.Body.Source {
	case "saga":
		if err := s.publisher.PublishCancelFlightBooking(ctx, saga); err != nil {
			logger.Errorw("failed to publish CancelFlightBooking", "command", cmd, "err", err.Error())
		}
		if err := s.publisher.PublishCancelHotelBooking(ctx, saga); err != nil {
			logger.Errorw("failed to publish CancelHotelBooking", "command", cmd, "err", err.Error())
		}
		if err := s.publisher.PublishCancelCarBooking(ctx, saga); err != nil {
			logger.Errorw("failed to publish CancelHotelBooking", "command", cmd, "err", err.Error())
		}
	case "flight":
		if err := s.publisher.PublishCancelHotelBooking(ctx, saga); err != nil {
			logger.Errorw("failed to publish CancelHotelBooking", "command", cmd, "err", err.Error())
		}
		if err := s.publisher.PublishCancelCarBooking(ctx, saga); err != nil {
			logger.Errorw("failed to publish CancelHotelBooking", "command", cmd, "err", err.Error())
		}
	case "hotel":
		if err := s.publisher.PublishCancelCarBooking(ctx, saga); err != nil {
			logger.Errorw("failed to publish CancelHotelBooking", "command", cmd, "err", err.Error())
		}
	}

	return nil
}
