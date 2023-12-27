package entity

import (
	"time"

	"github.com/haandol/hexagonal/internal/dto"
)

type FlightBooking struct {
	ID        uint      `gorm:"type:bigint;primaryKey;autoIncrement;<-:create;"`
	TripID    uint      `gorm:"type:bigint;not null;"`
	FlightID  uint      `gorm:"type:bigint;not null;"`
	Status    string    `gorm:"type:varchar(32);not null;"`
	CreatedAt time.Time `gorm:"type:timestamp;<-:create;"`
	UpdatedAt time.Time `gorm:"type:timestamp;"`
}

func (m *FlightBooking) DTO() dto.FlightBooking {
	return dto.FlightBooking{
		ID:        m.ID,
		TripID:    m.TripID,
		FlightID:  m.FlightID,
		Status:    m.Status,
		CreatedAt: m.CreatedAt,
	}
}
