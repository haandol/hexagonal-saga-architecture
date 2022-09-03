package repository

import (
	"context"

	"github.com/haandol/hexagonal/pkg/dto"
	"github.com/haandol/hexagonal/pkg/entity"
	"gorm.io/gorm"
)

type CarRepository struct {
	db *gorm.DB
}

func NewCarRepository(db *gorm.DB) *CarRepository {
	return &CarRepository{
		db: db,
	}
}

func (r *CarRepository) Rent(ctx context.Context, d *dto.CarRental) (dto.CarRental, error) {
	row := &entity.CarRental{
		TripID:   d.TripID,
		CarID:    d.CarID,
		Quantity: d.Quantity,
	}
	result := r.db.WithContext(ctx).
		Where("trip_id = ? AND car_id = ? AND quantity = ?", d.TripID, d.CarID, d.Quantity).
		FirstOrCreate(row)
	if result.Error != nil {
		return dto.CarRental{}, result.Error
	}

	return row.DTO()
}

func (r *CarRepository) CancelRental(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).
		Where("id = ?", id).
		Unscoped().
		Delete(&entity.CarRental{}, id).Error
}
