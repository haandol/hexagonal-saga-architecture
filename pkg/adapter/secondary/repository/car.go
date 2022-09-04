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

func (r *CarRepository) Book(ctx context.Context, d *dto.CarBooking) (dto.CarBooking, error) {
	row := &entity.CarBooking{
		TripID: d.TripID,
		CarID:  d.CarID,
	}
	result := r.db.WithContext(ctx).Create(row)
	if result.Error != nil {
		return dto.CarBooking{}, result.Error
	}

	return row.DTO()
}

func (r *CarRepository) CancelBooking(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).
		Where("id = ?", id).
		Unscoped().
		Delete(&entity.CarBooking{}, id).Error
}
