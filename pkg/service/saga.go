package service

import (
	"context"

	"github.com/haandol/hexagonal/pkg/constant/status"
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

	if err := s.sagaRepository.Start(ctx, cmd); err != nil {
		logger.Errorw("failed to create saga", "command", cmd, "err", err.Error())
		return err
	}

	return nil
}

func (s *SagaService) ProcessCarBooking(ctx context.Context, evt *event.CarBooked) error {
	logger := util.GetLogger().With(
		"module", "SagaService",
		"method", "ProcessCarBooking",
	)

	logger.Infow("success car booked", "event", evt)

	if err := s.sagaRepository.ProcessCarBooking(ctx, evt); err != nil {
		logger.Errorw("failed to process car booked", "event", evt, "err", err.Error())
		return err
	}

	return nil
}

func (s *SagaService) CompensateCarBooking(ctx context.Context, evt *event.CarBookingCancelled) error {
	logger := util.GetLogger().With(
		"module", "SagaService",
		"method", "CompensateCarBooking",
	)

	logger.Infow("cancel car booking", "event", evt)

	if err := s.sagaRepository.CompensateCarBooking(ctx, evt); err != nil {
		logger.Errorw("failed to process cancel car booking", "event", evt, "err", err.Error())
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

	if err := s.sagaRepository.ProcessHotelBooking(ctx, evt); err != nil {
		logger.Errorw("failed to process Hotel booked", "event", evt, "err", err.Error())
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

	if err := s.sagaRepository.ProcessFlightBooking(ctx, evt); err != nil {
		logger.Errorw("failed to process flight booked", "event", evt, "err", err.Error())
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

	if err := s.sagaRepository.End(ctx, cmd); err != nil {
		logger.Errorw("failed to end saga", "command", cmd, "err", err.Error())
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
	case "saga", "trip":
		if err := s.publisher.PublishCancelFlightBooking(ctx, saga); err != nil {
			logger.Errorw("failed to publish CancelFlightBooking", "command", cmd, "err", err.Error())
			return err
		}
		if err := s.publisher.PublishCancelHotelBooking(ctx, saga); err != nil {
			logger.Errorw("failed to publish CancelHotelBooking", "command", cmd, "err", err.Error())
			return err
		}
		if err := s.publisher.PublishCancelCarBooking(ctx, saga); err != nil {
			logger.Errorw("failed to publish CancelHotelBooking", "command", cmd, "err", err.Error())
			return err
		}
	case "flight":
		if err := s.publisher.PublishCancelHotelBooking(ctx, saga); err != nil {
			logger.Errorw("failed to publish CancelHotelBooking", "command", cmd, "err", err.Error())
			return err
		}
		if err := s.publisher.PublishCancelCarBooking(ctx, saga); err != nil {
			logger.Errorw("failed to publish CancelHotelBooking", "command", cmd, "err", err.Error())
			return err
		}
	case "hotel":
		if err := s.publisher.PublishCancelCarBooking(ctx, saga); err != nil {
			logger.Errorw("failed to publish CancelHotelBooking", "command", cmd, "err", err.Error())
			return err
		}
	}

	return nil
}

func (s *SagaService) MarkAbort(ctx context.Context, tripID uint) error {
	logger := util.GetLogger().With(
		"module", "SagaService",
		"method", "MarkAbort",
	)

	if err := s.sagaRepository.UpdateStatusByTripID(ctx, tripID, status.SagaAborted); err != nil {
		logger.Errorw("failed to update saga status", "tripID", tripID, "err", err.Error())
		return err
	}

	return nil
}
