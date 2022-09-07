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
	ErrNoFlightBookingFound = errors.New("no flight-booking found")
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
	booking, err := r.GetByTripID(ctx, d.TripID)
	if err != nil {
		return dto.FlightBooking{}, err
	}
	if booking.Status == status.Booked {
		return booking, nil
	}

	row := &entity.FlightBooking{
		TripID:   d.TripID,
		FlightID: d.FlightID,
		Status:   status.Booked,
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
		Model(row).
		Clauses(clause.Returning{}).
		Where("id = ?", id).
		Update("status", status.Cancelled)
	if result.Error != nil {
		return dto.FlightBooking{}, result.Error
	}
	if result.RowsAffected == 0 {
		return dto.FlightBooking{}, ErrNoFlightBookingFound
	}

	return row.DTO()
}

func (r *FlightRepository) GetByTripID(ctx context.Context, tripID uint) (dto.FlightBooking, error) {
	row := &entity.FlightBooking{}
	result := r.db.WithContext(ctx).
		Where("trip_id = ?", tripID).
		Limit(1).
		Find(&row)
	if result.Error != nil {
		return dto.FlightBooking{}, result.Error
	}
	return row.DTO()
}
