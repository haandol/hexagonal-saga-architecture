package repository

import (
	"context"

	"github.com/haandol/hexagonal/pkg/dto"
	"github.com/haandol/hexagonal/pkg/entity"
	"gorm.io/gorm"
)

type FlightRepository struct {
	db *gorm.DB
}

func NewFlightRepository(db *gorm.DB) *FlightRepository {
	return &FlightRepository{
		db: db,
	}
}

func (r *FlightRepository) Book(ctx context.Context, d *dto.FlightBooking) (dto.FlightBooking, error) {
	row := &entity.FlightBooking{
		TripID:   d.TripID,
		FlightID: d.FlightID,
	}
	result := r.db.WithContext(ctx).
		Where("trip_id = ? AND flight_id = ?", d.TripID, d.FlightID).
		FirstOrCreate(row)
	if result.Error != nil {
		return dto.FlightBooking{}, result.Error
	}

	return row.DTO()
}

func (r *FlightRepository) CancelBooking(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).
		Where("id = ?", id).
		Unscoped().
		Delete(&entity.FlightBooking{}, id).Error
}
