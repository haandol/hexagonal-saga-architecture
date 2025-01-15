package service

import (
	"context"

	"github.com/haandol/hexagonal/internal/adapter/secondary/producer"
	"github.com/haandol/hexagonal/internal/adapter/secondary/repository"
	"github.com/haandol/hexagonal/internal/constant/status"
	"github.com/haandol/hexagonal/internal/instrument"
	"github.com/haandol/hexagonal/internal/message/command"
	"github.com/haandol/hexagonal/internal/message/event"
	"github.com/haandol/hexagonal/internal/port/secondaryport/producerport"
	"github.com/haandol/hexagonal/internal/port/secondaryport/repositoryport"
	"github.com/haandol/hexagonal/pkg/o11y"
	"github.com/haandol/hexagonal/pkg/util"
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
	logger := util.LoggerFromContext(ctx).WithGroup("SagaService.Start")

	ctx, span := o11y.BeginSubSpan(ctx, "Start")
	defer span.End()

	if err := s.sagaRepository.Start(ctx, cmd); err != nil {
		instrument.RecordStartSagaError(logger, span, err, cmd)
		return err
	}

	return nil
}

func (s *SagaService) ProcessCarBooking(ctx context.Context, evt *event.CarBooked) error {
	logger := util.LoggerFromContext(ctx).WithGroup("SagaService.ProcessCarBooking")

	ctx, span := o11y.BeginSubSpan(ctx, "ProcessCarBooking")
	defer span.End()

	if err := s.sagaRepository.ProcessCarBooking(ctx, evt); err != nil {
		instrument.RecordProcessSagaEventError(logger, span, err, evt)
		return err
	}

	return nil
}

func (s *SagaService) CompensateCarBooking(ctx context.Context, evt *event.CarBookingCanceled) error {
	logger := util.LoggerFromContext(ctx).WithGroup("SagaService.CompensateCarBooking")

	ctx, span := o11y.BeginSubSpan(ctx, "CompensateCarBooking")
	defer span.End()

	if err := s.sagaRepository.CompensateCarBooking(ctx, evt); err != nil {
		instrument.RecordCompensateSagaEventError(logger, span, err, evt)
		return err
	}

	return nil
}

func (s *SagaService) ProcessHotelBooking(ctx context.Context, evt *event.HotelBooked) error {
	logger := util.LoggerFromContext(ctx).WithGroup("SagaService.ProcessHotelBooking")

	ctx, span := o11y.BeginSubSpan(ctx, "ProcessHotelBooking")
	defer span.End()

	if err := s.sagaRepository.ProcessHotelBooking(ctx, evt); err != nil {
		instrument.RecordProcessSagaEventError(logger, span, err, evt)
		return err
	}

	return nil
}

func (s *SagaService) CompensateHotelBooking(ctx context.Context, evt *event.HotelBookingCanceled) error {
	logger := util.LoggerFromContext(ctx).WithGroup("SagaService.CompensateHotelBooking")

	ctx, span := o11y.BeginSubSpan(ctx, "CompensateHotelBooking")
	defer span.End()

	_, err := s.sagaRepository.CompensateHotelBooking(ctx, evt)
	if err != nil {
		instrument.RecordCompensateSagaEventError(logger, span, err, evt)
		return err
	}

	return nil
}

func (s *SagaService) ProcessFlightBooking(ctx context.Context, evt *event.FlightBooked) error {
	logger := util.LoggerFromContext(ctx).WithGroup("SagaService.ProcessFlightBooking")

	ctx, span := o11y.BeginSubSpan(ctx, "ProcessFlightBooking")
	defer span.End()

	if err := s.sagaRepository.ProcessFlightBooking(ctx, evt); err != nil {
		instrument.RecordProcessSagaEventError(logger, span, err, evt)
		return err
	}

	return nil
}

func (s *SagaService) CompensateFlightBooking(ctx context.Context, evt *event.FlightBookingCanceled) error {
	logger := util.LoggerFromContext(ctx).WithGroup("SagaService.CompensateFlightBooking")

	ctx, span := o11y.BeginSubSpan(ctx, "CompensateFlightBooking")
	defer span.End()

	_, err := s.sagaRepository.CompensateFlightBooking(ctx, evt)
	if err != nil {
		instrument.RecordCompensateSagaEventError(logger, span, err, evt)
		return err
	}

	return nil
}

func (s *SagaService) End(ctx context.Context, cmd *command.EndSaga) error {
	logger := util.LoggerFromContext(ctx).WithGroup("SagaService.End")

	ctx, span := o11y.BeginSubSpan(ctx, "End")
	defer span.End()

	if err := s.sagaRepository.End(ctx, cmd); err != nil {
		instrument.RecordEndSagaError(logger, span, err, cmd)
		return err
	}

	return nil
}

func (s *SagaService) Abort(ctx context.Context, cmd *command.AbortSaga) error {
	logger := util.LoggerFromContext(ctx).WithGroup("SagaService.Abort")

	ctx, span := o11y.BeginSubSpan(ctx, "Abort")
	defer span.End()

	saga, err := s.sagaRepository.Abort(ctx, cmd)
	if err != nil {
		instrument.RecordAbortSagaError(logger, span, err, cmd)
		return err
	}

	switch cmd.Body.Source {
	case "saga", "trip":
		if err := s.publisher.PublishCancelFlightBooking(ctx, &saga); err != nil {
			instrument.RecordPublishSagaCommandError(logger, span, err, cmd)
			return err
		}
		if err := s.publisher.PublishCancelHotelBooking(ctx, &saga); err != nil {
			instrument.RecordPublishSagaCommandError(logger, span, err, cmd)
			return err
		}
		if err := s.publisher.PublishCancelCarBooking(ctx, &saga); err != nil {
			instrument.RecordPublishSagaCommandError(logger, span, err, cmd)
			return err
		}
	case "flight":
		if err := s.publisher.PublishCancelHotelBooking(ctx, &saga); err != nil {
			instrument.RecordPublishSagaCommandError(logger, span, err, cmd)
			return err
		}
		if err := s.publisher.PublishCancelCarBooking(ctx, &saga); err != nil {
			instrument.RecordPublishSagaCommandError(logger, span, err, cmd)
			return err
		}
	case "hotel":
		if err := s.publisher.PublishCancelCarBooking(ctx, &saga); err != nil {
			instrument.RecordPublishSagaCommandError(logger, span, err, cmd)
			return err
		}
	}

	return nil
}

func (s *SagaService) MarkAbort(ctx context.Context, tripID uint) error {
	logger := util.LoggerFromContext(ctx).WithGroup("SagaService.MarkAbort").With(
		"tripID", tripID,
	)

	ctx, span := o11y.BeginSubSpan(ctx, "MarkAbort")
	defer span.End()

	if err := s.sagaRepository.UpdateStatusByTripID(ctx, tripID, status.SagaAborted); err != nil {
		instrument.RecordUpdateSagaError(logger, span, err, tripID)
		return err
	}

	span.SetAttributes(
		o11y.AttrInt("tripID", int(tripID)),
	)

	return nil
}
