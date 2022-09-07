package repository

import (
	"context"
	"errors"

	"github.com/haandol/hexagonal/pkg/constant/status"
	"github.com/haandol/hexagonal/pkg/dto"
	"github.com/haandol/hexagonal/pkg/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrNoCarBookingFound = errors.New("no car-booking found")
)

type CarRepository struct {
	db *gorm.DB
}

func NewCarRepository(db *gorm.DB) *CarRepository {
	return &CarRepository{
		db: db,
	}
}

func (r *CarRepository) Book(ctx context.Context, d *dto.CarBooking) (dto.CarBooking, error) {
	booking, err := r.GetByTripID(ctx, d.TripID)
	if err != nil {
		return dto.CarBooking{}, err
	}
	if booking.Status == status.Booked {
		return booking, nil
	}

	row := &entity.CarBooking{
		TripID: d.TripID,
		CarID:  d.CarID,
		Status: status.Booked,
	}
	result := r.db.WithContext(ctx).Create(row)
	if result.Error != nil {
		return dto.CarBooking{}, result.Error
	}

	return row.DTO()
}

func (r *CarRepository) CancelBooking(ctx context.Context, id uint) (dto.CarBooking, error) {
	row := &entity.CarBooking{}
	result := r.db.WithContext(ctx).
		Model(row).
		Clauses(clause.Returning{}).
		Where("id = ?", id).
		Update("status", status.Cancelled)
	if result.Error != nil {
		return dto.CarBooking{}, result.Error
	}
	if result.RowsAffected == 0 {
		return dto.CarBooking{}, ErrNoCarBookingFound
	}

	return row.DTO()
}

func (r *CarRepository) GetByTripID(ctx context.Context, tripID uint) (dto.CarBooking, error) {
	row := &entity.CarBooking{}
	result := r.db.WithContext(ctx).
		Where("trip_id = ?", tripID).
		Limit(1).
		Find(&row)
	if result.Error != nil {
		return dto.CarBooking{}, result.Error
	}
	return row.DTO()
}
