package repository

import (
	"context"
	"fmt"

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

func (r *TripRepository) Create(ctx context.Context, userID uint) (dto.Trip, error) {
	row := &entity.Trip{
		UserID: userID,
		Status: "INITIALIZED",
	}
	result := r.db.WithContext(ctx).Create(row)
	if result.Error != nil {
		return dto.Trip{}, result.Error
	}

	return row.DTO()
}

func (r *TripRepository) Update(ctx context.Context, d *dto.Trip) error {
	result := r.db.WithContext(ctx).
		Where("id = ? AND user_id", d.ID, d.UserID).
		Updates(&entity.Trip{
			ID:       d.ID,
			UserID:   d.UserID,
			HotelID:  d.HotelID,
			CarID:    d.CarID,
			FlightID: d.FlightID,
			Status:   d.Status,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("no rows affected")
	}
	return nil
}

func (r *TripRepository) List(ctx context.Context) ([]dto.Trip, error) {
	var rows entity.Trips
	result := r.db.WithContext(ctx).Find(&rows)
	if result.Error != nil {
		return []dto.Trip{}, result.Error
	}

	return rows.DTO()
}
