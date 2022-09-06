package repository

import (
	"context"
	"errors"

	"github.com/haandol/hexagonal/pkg/constant"
	"github.com/haandol/hexagonal/pkg/dto"
	"github.com/haandol/hexagonal/pkg/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type HotelRepository struct {
	db *gorm.DB
}

func NewHotelRepository(db *gorm.DB) *HotelRepository {
	return &HotelRepository{
		db: db,
	}
}

func (r *HotelRepository) Book(ctx context.Context, d *dto.HotelBooking) (dto.HotelBooking, error) {
	booking, err := r.GetByTripID(ctx, d.TripID)
	if err != nil {
		return dto.HotelBooking{}, err
	}
	if booking.Status == constant.Booked {
		return booking, nil
	}

	row := &entity.HotelBooking{
		TripID:  d.TripID,
		HotelID: d.HotelID,
		Status:  constant.Booked,
	}
	result := r.db.WithContext(ctx).Create(row)
	if result.Error != nil {
		return dto.HotelBooking{}, result.Error
	}

	return row.DTO()
}

func (r *HotelRepository) CancelBooking(ctx context.Context, id uint) (dto.HotelBooking, error) {
	row := &entity.HotelBooking{}
	result := r.db.WithContext(ctx).
		Model(row).
		Clauses(clause.Returning{}).
		Where("id = ?", id).
		Update("status", constant.Cancelled)
	if result.Error != nil {
		return dto.HotelBooking{}, result.Error
	}
	if result.RowsAffected == 0 {
		return dto.HotelBooking{}, errors.New("booking not found")
	}

	return row.DTO()
}

func (r *HotelRepository) GetByTripID(ctx context.Context, tripID uint) (dto.HotelBooking, error) {
	row := &entity.HotelBooking{}
	result := r.db.WithContext(ctx).
		Where("trip_id = ?", tripID).
		Limit(1).
		Find(&row)
	if result.Error != nil {
		return dto.HotelBooking{}, result.Error
	}
	return row.DTO()
}
