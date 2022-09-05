package repository

import (
	"context"

	"github.com/haandol/hexagonal/pkg/dto"
	"github.com/haandol/hexagonal/pkg/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
	result := r.db.WithContext(ctx).Create(row)
	if result.Error != nil {
		return dto.FlightBooking{}, result.Error
	}

	return row.DTO()
}

func (r *FlightRepository) CancelBooking(ctx context.Context, id uint) (dto.FlightBooking, error) {
	row := &entity.FlightBooking{}
	result := r.db.WithContext(ctx).
		Clauses(clause.Returning{}).
		Where("id = ?", id).
		Unscoped().
		Delete(row, id)
	if result.Error != nil {
		return dto.FlightBooking{}, result.Error
	}

	return row.DTO()
}
