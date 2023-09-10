package service

import (
	"context"

	"github.com/haandol/hexagonal/pkg/adapter/secondary/producer"
	"github.com/haandol/hexagonal/pkg/adapter/secondary/repository"
	"github.com/haandol/hexagonal/pkg/constant/status"
	"github.com/haandol/hexagonal/pkg/message/command"
	"github.com/haandol/hexagonal/pkg/message/event"
	"github.com/haandol/hexagonal/pkg/port/secondaryport/producerport"
	"github.com/haandol/hexagonal/pkg/port/secondaryport/repositoryport"
	"github.com/haandol/hexagonal/pkg/util"
	"github.com/haandol/hexagonal/pkg/util/o11y"
)

type SagaService struct {
	publisher      producerport.SagaProducer
	sagaRepository repositoryport.SagaRepository
}

func NewSagaService(
	publisher *producer.SagaProducer,
	sagaRepository *repository.SagaRepository,
) *SagaService {
	return &SagaService{
		publisher:      publisher,
		sagaRepository: sagaRepository,
	}
}

func (s *SagaService) Start(ctx context.Context, cmd *command.StartSaga) error {
	logger := util.LoggerFromContext(ctx).With(
		"module", "SagaService",
		"method", "Start",
		"command", cmd,
	)

	ctx, span := o11y.BeginSubSpan(ctx, "Start")
	defer span.End()

	if err := s.sagaRepository.Start(ctx, cmd); err != nil {
		logger.Error("failed to create saga", "command", cmd, "err", err)
		span.RecordError(err)
		span.SetStatus(o11y.GetStatus(err))
		return err
	}

	return nil
}

func (s *SagaService) ProcessCarBooking(ctx context.Context, evt *event.CarBooked) error {
	logger := util.LoggerFromContext(ctx).With(
		"module", "SagaService",
		"method", "ProcessCarBooking",
		"event", evt,
	)

	ctx, span := o11y.BeginSubSpan(ctx, "ProcessCarBooking")
	defer span.End()

	if err := s.sagaRepository.ProcessCarBooking(ctx, evt); err != nil {
		logger.Error("failed to process car booked", "event", evt, "err", err)
		span.RecordError(err)
		span.SetStatus(o11y.GetStatus(err))
		return err
	}

	return nil
}

func (s *SagaService) CompensateCarBooking(ctx context.Context, evt *event.CarBookingCancelled) error {
	logger := util.LoggerFromContext(ctx).With(
		"module", "SagaService",
		"method", "CompensateCarBooking",
		"event", evt,
	)

	ctx, span := o11y.BeginSubSpan(ctx, "CompensateCarBooking")
	defer span.End()

	if err := s.sagaRepository.CompensateCarBooking(ctx, evt); err != nil {
		logger.Error("failed to process cancel car booking", "event", evt, "err", err)
		span.RecordError(err)
		span.SetStatus(o11y.GetStatus(err))
		return err
	}

	return nil
}

func (s *SagaService) ProcessHotelBooking(ctx context.Context, evt *event.HotelBooked) error {
	logger := util.LoggerFromContext(ctx).With(
		"module", "SagaService",
		"method", "ProcessHotelBooking",
		"event", evt,
	)

	ctx, span := o11y.BeginSubSpan(ctx, "ProcessHotelBooking")
	defer span.End()

	if err := s.sagaRepository.ProcessHotelBooking(ctx, evt); err != nil {
		logger.Error("failed to process Hotel booked", "event", evt, "err", err)
		span.RecordError(err)
		span.SetStatus(o11y.GetStatus(err))
		return err
	}

	return nil
}

func (s *SagaService) CompensateHotelBooking(ctx context.Context, evt *event.HotelBookingCancelled) error {
	logger := util.LoggerFromContext(ctx).With(
		"module", "SagaService",
		"method", "CompensateHotelBooking",
		"event", evt,
	)

	ctx, span := o11y.BeginSubSpan(ctx, "CompensateHotelBooking")
	defer span.End()

	_, err := s.sagaRepository.CompensateHotelBooking(ctx, evt)
	if err != nil {
		logger.Error("failed to process cancel Hotel booking", "event", evt, "err", err)
		span.RecordError(err)
		span.SetStatus(o11y.GetStatus(err))
		return err
	}

	return nil
}

func (s *SagaService) ProcessFlightBooking(ctx context.Context, evt *event.FlightBooked) error {
	logger := util.LoggerFromContext(ctx).With(
		"module", "SagaService",
		"method", "ProcessFlightBooking",
		"event", evt,
	)

	ctx, span := o11y.BeginSubSpan(ctx, "ProcessFlightBooking")
	defer span.End()

	if err := s.sagaRepository.ProcessFlightBooking(ctx, evt); err != nil {
		logger.Error("failed to process flight booked", "event", evt, "err", err)
		span.RecordError(err)
		span.SetStatus(o11y.GetStatus(err))
		return err
	}

	return nil
}

func (s *SagaService) CompensateFlightBooking(ctx context.Context, evt *event.FlightBookingCancelled) error {
	logger := util.LoggerFromContext(ctx).With(
		"module", "SagaService",
		"method", "CompensateFlightBooking",
		"event", evt,
	)

	ctx, span := o11y.BeginSubSpan(ctx, "CompensateFlightBooking")
	defer span.End()

	_, err := s.sagaRepository.CompensateFlightBooking(ctx, evt)
	if err != nil {
		logger.Error("failed to process cancel flight booking", "event", evt, "err", err)
		span.RecordError(err)
		span.SetStatus(o11y.GetStatus(err))
		return err
	}

	return nil
}

func (s *SagaService) End(ctx context.Context, cmd *command.EndSaga) error {
	logger := util.LoggerFromContext(ctx).With(
		"module", "SagaService",
		"method", "End",
		"command", cmd,
	)

	ctx, span := o11y.BeginSubSpan(ctx, "End")
	defer span.End()

	if err := s.sagaRepository.End(ctx, cmd); err != nil {
		logger.Error("failed to end saga", "command", cmd, "err", err)
		span.RecordError(err)
		span.SetStatus(o11y.GetStatus(err))
		return err
	}

	return nil
}

func (s *SagaService) Abort(ctx context.Context, cmd *command.AbortSaga) error {
	logger := util.LoggerFromContext(ctx).With(
		"module", "SagaService",
		"method", "Abort",
		"command", cmd,
	)

	ctx, span := o11y.BeginSubSpan(ctx, "Abort")
	defer span.End()

	saga, err := s.sagaRepository.Abort(ctx, cmd)
	if err != nil {
		logger.Error("failed to abort saga", "command", cmd, "err", err)
		span.RecordError(err)
		span.SetStatus(o11y.GetStatus(err))
		return err
	}

	switch cmd.Body.Source {
	case "saga", "trip":
		if err := s.publisher.PublishCancelFlightBooking(ctx, &saga); err != nil {
			logger.Error("failed to publish CancelFlightBooking", "command", cmd, "err", err)
			span.RecordError(err)
			span.SetStatus(o11y.GetStatus(err))
			return err
		}
		if err := s.publisher.PublishCancelHotelBooking(ctx, &saga); err != nil {
			logger.Error("failed to publish CancelHotelBooking", "command", cmd, "err", err)
			span.RecordError(err)
			span.SetStatus(o11y.GetStatus(err))
			return err
		}
		if err := s.publisher.PublishCancelCarBooking(ctx, &saga); err != nil {
			logger.Error("failed to publish CancelHotelBooking", "command", cmd, "err", err)
			span.RecordError(err)
			span.SetStatus(o11y.GetStatus(err))
			return err
		}
	case "flight":
		if err := s.publisher.PublishCancelHotelBooking(ctx, &saga); err != nil {
			logger.Error("failed to publish CancelHotelBooking", "command", cmd, "err", err)
			span.RecordError(err)
			span.SetStatus(o11y.GetStatus(err))
			return err
		}
		if err := s.publisher.PublishCancelCarBooking(ctx, &saga); err != nil {
			logger.Error("failed to publish CancelHotelBooking", "command", cmd, "err", err)
			span.RecordError(err)
			span.SetStatus(o11y.GetStatus(err))
			return err
		}
	case "hotel":
		if err := s.publisher.PublishCancelCarBooking(ctx, &saga); err != nil {
			logger.Error("failed to publish CancelHotelBooking", "command", cmd, "err", err)
			span.RecordError(err)
			span.SetStatus(o11y.GetStatus(err))
			return err
		}
	}

	return nil
}

func (s *SagaService) MarkAbort(ctx context.Context, tripID uint) error {
	logger := util.LoggerFromContext(ctx).With(
		"module", "SagaService",
		"method", "MarkAbort",
		"tripID", tripID,
	)

	ctx, span := o11y.BeginSubSpan(ctx, "MarkAbort")
	defer span.End()

	if err := s.sagaRepository.UpdateStatusByTripID(ctx, tripID, status.SagaAborted); err != nil {
		logger.Error("failed to update saga status", "tripID", tripID, "err", err)
		span.RecordError(err)
		span.SetStatus(o11y.GetStatus(err))
		return err
	}

	span.SetAttributes(
		o11y.AttrInt("tripID", int(tripID)),
	)

	return nil
}
