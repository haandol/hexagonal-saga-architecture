package service

import (
	"context"
	"errors"

	"github.com/haandol/hexagonal/internal/adapter/secondary/repository"
	"github.com/haandol/hexagonal/internal/constant/status"
	"github.com/haandol/hexagonal/internal/dto"
	"github.com/haandol/hexagonal/internal/message/event"
	"github.com/haandol/hexagonal/internal/port/secondaryport/repositoryport"
	"github.com/haandol/hexagonal/pkg/util"
	"github.com/haandol/hexagonal/pkg/util/o11y"
)

type TripService struct {
	tripRepository repositoryport.TripRepository
}

func NewTripService(
	tripRepository *repository.TripRepository,
) *TripService {
	return &TripService{
		tripRepository: tripRepository,
	}
}

// create trip for the given user
func (s *TripService) Create(ctx context.Context, d *dto.Trip) (dto.Trip, error) {
	logger := util.LoggerFromContext(ctx).WithGroup("TripService.Create")

	traceID, spanID := o11y.GetTraceSpanID(ctx)
	trip, err := s.tripRepository.Create(ctx, traceID, spanID, d)
	if err != nil {
		logger.Error("failed to create trip", "trip", d, "err", err)
		return dto.Trip{}, err
	}

	return trip, nil
}

func (s *TripService) RecoverForward(ctx context.Context, tripID uint) (dto.Trip, error) {
	logger := util.LoggerFromContext(ctx).WithGroup("TripService.RecoverForward")

	traceID, spanID := o11y.GetTraceSpanID(ctx)

	trip, err := s.tripRepository.GetByID(ctx, tripID)
	if err != nil {
		logger.Error("failed to get a trip", "traceID", traceID, "id", tripID, "err", err)
		return dto.Trip{}, err
	}

	if trip.Status == status.TripCanceled || trip.Status == status.TripBooked {
		return dto.Trip{}, errors.New("trip is already completed or aborted")
	}

	if err := s.tripRepository.PublishStartSaga(ctx, traceID, spanID, &trip); err != nil {
		logger.Error("failed to produce start saga", "trip", trip, "err", err)
	}

	return trip, nil
}

func (s *TripService) RecoverBackward(ctx context.Context, tripID uint) (dto.Trip, error) {
	logger := util.LoggerFromContext(ctx).WithGroup("TripService.RecoverBackward")
	traceID, spanID := o11y.GetTraceSpanID(ctx)

	trip, err := s.tripRepository.GetByID(ctx, tripID)
	if err != nil {
		logger.Error("failed to get a trip", "traceID", traceID, "id", tripID, "err", err)
		return dto.Trip{}, err
	}

	if trip.Status == status.TripCanceled || trip.Status == status.TripBooked {
		return dto.Trip{}, errors.New("trip is already completed or aborted")
	}

	if err := s.tripRepository.PublishAbortSaga(ctx,
		traceID, spanID, tripID, "force revert",
	); err != nil {
		logger.Error("failed to produce start saga", "trip", trip, "err", err)
	}

	return trip, nil
}

func (s *TripService) List(ctx context.Context) ([]dto.Trip, error) {
	logger := util.LoggerFromContext(ctx).WithGroup("TripService.List")

	trips, err := s.tripRepository.List(ctx)
	if err != nil {
		logger.Error("failed to create trip", "err", err)
		return []dto.Trip{}, err
	}

	return trips, nil
}

func (s *TripService) ProcessSagaEnded(ctx context.Context, evt *event.SagaEnded) error {
	logger := util.LoggerFromContext(ctx).WithGroup("TripService.ProcessSagaEnded")

	if err := s.tripRepository.Complete(ctx, evt); err != nil {
		logger.Error("failed to update trip booking", "event", evt, "err", err)
		return err
	}

	return nil
}

func (s *TripService) ProcessSagaAborted(ctx context.Context, evt *event.SagaAborted) error {
	logger := util.LoggerFromContext(ctx).WithGroup("TripService.ProcessSagaAborted")

	if err := s.tripRepository.Abort(ctx, evt); err != nil {
		logger.Error("failed to abort trip booking", "event", evt, "err", err)
		return err
	}

	return nil
}
