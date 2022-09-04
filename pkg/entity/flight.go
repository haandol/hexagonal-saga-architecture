package entity

import (
	"time"

	"github.com/haandol/hexagonal/pkg/dto"
)

type FlightBooking struct {
	ID        uint      `gorm:"type:bigint;primaryKey;autoIncrement;<-:create;"`
	TripID    uint      `gorm:"type:bigint;not null;"`
	FlightID  uint      `gorm:"type:bigint;not null;"`
	CreatedAt time.Time `gorm:"type:timestamp;<-:create;"`
	UpdatedAt time.Time `gorm:"type:timestamp;"`
}

func (m FlightBooking) DTO() (dto.FlightBooking, error) {
	return dto.FlightBooking{
		ID:        m.ID,
		TripID:    m.TripID,
		FlightID:  m.FlightID,
		CreatedAt: m.CreatedAt,
	}, nil
}
