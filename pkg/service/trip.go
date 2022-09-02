package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/haandol/hexagonal/message/command"
	"github.com/haandol/hexagonal/pkg/dto"
	"github.com/haandol/hexagonal/pkg/port/secondaryport/producerport"
	"github.com/haandol/hexagonal/pkg/port/secondaryport/repositoryport"
	"github.com/haandol/hexagonal/pkg/util"
)

type TripService struct {
	producer       producerport.Producer
	tripRepository repositoryport.TripRepository
}

func NewTripService(
	producer producerport.Producer,
	tripRepository repositoryport.TripRepository,
) *TripService {
	return &TripService{
		producer:       producer,
		tripRepository: tripRepository,
	}
}

func (s *TripService) Create(ctx context.Context, userID uint) (dto.Trip, error) {
	logger := util.GetLogger().With(
		"service", "TripService",
		"method", "Create",
	)

	trip, err := s.tripRepository.Create(ctx, userID)
	if err != nil {
		logger.Errorw("failed to create trip", "userID", userID, "err", err.Error())
		return dto.Trip{}, err
	}

	cmd := command.StartSaga{
		Command: command.Command{
			Name:          "StartTripSaga",
			Version:       "1.0",
			ID:            uuid.NewString(),
			CorrelationID: uuid.NewString(), // TODO: use client provided value
			CreatedAt:     time.Now().Format(time.RFC3339),
		},
		Body: command.StartSagaBody{
			TripID: trip.ID,
		},
	}

	v, err := json.Marshal(cmd)
	if err != nil {
		logger.Errorw("failed to marshal trip", "trip", trip, "err", err.Error())
	}

	if err := s.producer.Produce(ctx, "trip-service", cmd.CorrelationID, v); err != nil {
		logger.Errorw("failed to produce trip.created", "trip", trip, "err", err.Error())
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
