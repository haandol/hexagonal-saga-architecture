package service

import (
	"context"

	"github.com/haandol/hexagonal/pkg/dto"
	"github.com/haandol/hexagonal/pkg/port/secondaryport/repositoryport"
	"github.com/haandol/hexagonal/pkg/util"
)

type TripService struct {
	tripRepository repositoryport.TripRepository
}

func NewTripService(
	tripRepository repositoryport.TripRepository,
) *TripService {
	return &TripService{
		tripRepository: tripRepository,
	}
}

func (s *TripService) Create(ctx context.Context, t *dto.Trip) (dto.Trip, error) {
	logger := util.GetLogger().With(
		"service", "TripService",
		"method", "Create",
	)

	trip, err := s.tripRepository.Create(ctx, t)
	if err != nil {
		logger.Errorw("failed to create trip", "trip", t, "err", err.Error())
		return dto.Trip{}, err
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
