package repository

import (
	"context"

	"github.com/haandol/hexagonal/pkg/dto"
	"github.com/haandol/hexagonal/pkg/entity"
	"gorm.io/gorm"
)

type TripRepository struct {
	db *gorm.DB
}

func NewTripRepository(db *gorm.DB) *TripRepository {
	return &TripRepository{
		db: db,
	}
}

func (r *TripRepository) Create(ctx context.Context, t *dto.Trip) (dto.Trip, error) {
	row := &entity.Trip{
		UserID:   t.UserID,
		HotelID:  t.HotelID,
		CarID:    t.CarID,
		FlightID: t.FlightID,
		Status:   t.Status,
	}
	result := r.db.WithContext(ctx).Create(row)
	if result.Error != nil {
		return dto.Trip{}, result.Error
	}

	return row.DTO()
}

func (r *TripRepository) List(ctx context.Context) ([]dto.Trip, error) {
	var rows entity.Trips
	result := r.db.WithContext(ctx).Find(&rows)
	if result.Error != nil {
		return []dto.Trip{}, result.Error
	}

	return rows.DTO()
}
