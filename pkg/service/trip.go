package service

import (
	"context"
	"errors"

	"github.com/haandol/hexagonal/pkg/constant"
	"github.com/haandol/hexagonal/pkg/constant/status"
	"github.com/haandol/hexagonal/pkg/dto"
	"github.com/haandol/hexagonal/pkg/message/event"
	"github.com/haandol/hexagonal/pkg/port/secondaryport/producerport"
	"github.com/haandol/hexagonal/pkg/port/secondaryport/repositoryport"
	"github.com/haandol/hexagonal/pkg/util"
)

type TripService struct {
	tripProducer   producerport.TripProducer
	tripRepository repositoryport.TripRepository
}

func NewTripService(
	tripProducer producerport.TripProducer,
	tripRepository repositoryport.TripRepository,
) *TripService {
	return &TripService{
		tripProducer:   tripProducer,
		tripRepository: tripRepository,
	}
}

// create trip for the given user
func (s *TripService) Create(ctx context.Context, d *dto.Trip) (dto.Trip, error) {
	logger := util.GetLogger().With(
		"service", "TripService",
		"method", "Create",
	)

	trip, err := s.tripRepository.Create(ctx, d)
	if err != nil {
		logger.Errorw("failed to create trip", "trip", d, "err", err.Error())
		return dto.Trip{}, err
	}

	corrID := ctx.Value(constant.CtxTraceKey).(string)
	if err := s.tripProducer.PublishStartSaga(ctx, corrID, trip); err != nil {
		logger.Errorw("failed to produce start saga", "trip", trip, "err", err.Error())
	}

	return trip, nil
}

func (s *TripService) RecoverForward(ctx context.Context, tripID uint) (dto.Trip, error) {
	logger := util.GetLogger().With(
		"service", "TripService",
		"method", "RecoverForward",
	)

	corrID := ctx.Value(constant.CtxTraceKey).(string)

	trip, err := s.tripRepository.GetByID(ctx, tripID)
	if err != nil {
		logger.Errorw("failed to get a trip", "corrID", corrID, "id", tripID, "err", err.Error())
		return dto.Trip{}, err
	}

	if trip.Status == status.TripCancelled || trip.Status == status.TripBooked {
		return dto.Trip{}, errors.New("trip is already completed or aborted")
	}

	if err := s.tripProducer.PublishStartSaga(ctx, corrID, trip); err != nil {
		logger.Errorw("failed to produce start saga", "trip", trip, "err", err.Error())
	}

	return trip, nil
}

func (s *TripService) RecoverBackward(ctx context.Context, tripID uint) (dto.Trip, error) {
	logger := util.GetLogger().With(
		"service", "TripService",
		"method", "RecoverBackward",
	)

	corrID := ctx.Value(constant.CtxTraceKey).(string)

	trip, err := s.tripRepository.GetByID(ctx, tripID)
	if err != nil {
		logger.Errorw("failed to get a trip", "corrID", corrID, "id", tripID, "err", err.Error())
		return dto.Trip{}, err
	}

	if trip.Status == status.TripCancelled || trip.Status == status.TripBooked {
		return dto.Trip{}, errors.New("trip is already completed or aborted")
	}

	if err := s.tripProducer.PublishAbortSaga(ctx, corrID, trip); err != nil {
		logger.Errorw("failed to produce start saga", "trip", trip, "err", err.Error())
	}

	return trip, nil
}

func (s *TripService) List(ctx context.Context) ([]dto.Trip, error) {
	logger := util.GetLogger().With(
		"service", "TripService",
		"method", "List",
	)

	trips, err := s.tripRepository.List(ctx)
	if err != nil {
		logger.Errorw("failed to create trip", "err", err.Error())
		return []dto.Trip{}, err
	}

	return trips, nil
}

func (s *TripService) ProcessSagaEnded(ctx context.Context, evt *event.SagaEnded) error {
	logger := util.GetLogger().With(
		"service", "TripService",
		"method", "ProcessSagaEnded",
	)

	if err := s.tripRepository.Complete(ctx, evt); err != nil {
		logger.Errorw("failed to update trip booking", "event", evt, "err", err.Error())
		return err
	}

	return nil
}

func (s *TripService) ProcessSagaAborted(ctx context.Context, evt *event.SagaAborted) error {
	logger := util.GetLogger().With(
		"service", "TripService",
		"method", "ProcessSagaAborted",
	)

	if err := s.tripRepository.Abort(ctx, evt); err != nil {
		logger.Errorw("failed to abort trip booking", "event", evt, "err", err.Error())
		return err
	}

	return nil
}
