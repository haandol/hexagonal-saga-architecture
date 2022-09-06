package repository

import (
	"context"
	"errors"

	"github.com/haandol/hexagonal/pkg/constant/status"
	"github.com/haandol/hexagonal/pkg/dto"
	"github.com/haandol/hexagonal/pkg/entity"
	"github.com/haandol/hexagonal/pkg/message/event"
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

func (r *TripRepository) Create(ctx context.Context, d *dto.Trip) (dto.Trip, error) {
	row := &entity.Trip{
		UserID:   d.UserID,
		CarID:    d.CarID,
		HotelID:  d.HotelID,
		FlightID: d.FlightID,
		Status:   status.TripInitialized,
	}
	result := r.db.WithContext(ctx).Create(row)
	if result.Error != nil {
		return dto.Trip{}, result.Error
	}

	return row.DTO()
}

func (r *TripRepository) Update(ctx context.Context, d *dto.Trip) error {
	result := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", d.ID, d.UserID).
		Updates(&entity.Trip{
			ID:       d.ID,
			UserID:   d.UserID,
			CarID:    d.CarID,
			HotelID:  d.HotelID,
			FlightID: d.FlightID,
			Status:   d.Status,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no rows affected")
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

func (r *TripRepository) Complete(ctx context.Context, evt *event.SagaEnded) error {
	return r.db.WithContext(ctx).
		Where("id = ?", evt.Body.TripID).
		Updates(&entity.Trip{
			CarBookingID:    evt.Body.CarBookingID,
			HotelBookingID:  evt.Body.HotelBookingID,
			FlightBookingID: evt.Body.FlightBookingID,
			Status:          status.TripCompleted,
		}).Error
}

func (r *TripRepository) Abort(ctx context.Context, evt *event.SagaAborted) error {
	return r.db.WithContext(ctx).
		Where("id = ?", evt.Body.TripID).
		Updates(&entity.Trip{
			Status: status.TripAborted,
		}).Error
}
