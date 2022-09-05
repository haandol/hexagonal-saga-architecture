package repository

import (
	"context"

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
	row := &entity.HotelBooking{
		TripID:  d.TripID,
		HotelID: d.HotelID,
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
		Clauses(clause.Returning{}).
		Where("id = ?", id).
		Unscoped().
		Delete(row, id)
	if result.Error != nil {
		return dto.HotelBooking{}, result.Error
	}

	return row.DTO()
}
