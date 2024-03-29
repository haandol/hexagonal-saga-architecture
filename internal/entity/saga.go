package entity

import (
	"time"

	"github.com/haandol/hexagonal/internal/dto"
	"gorm.io/datatypes"
)

type Saga struct {
	ID              uint           `gorm:"type:bigint;primaryKey;autoIncrement;<-:create;"`
	CorrelationID   string         `gorm:"type:varchar(36);<-:create;"`
	TripID          uint           `gorm:"type:bigint;<-:create;"`
	CarID           uint           `gorm:"type:bigint;"`
	HotelID         uint           `gorm:"type:bigint;"`
	FlightID        uint           `gorm:"type:bigint;"`
	CarBookingID    uint           `gorm:"type:bigint;"`
	HotelBookingID  uint           `gorm:"type:bigint;"`
	FlightBookingID uint           `gorm:"type:bigint;"`
	Status          string         `gorm:"type:varchar(32);"`
	History         datatypes.JSON `gorm:"type:json;"`
	CreatedAt       time.Time      `gorm:"type:timestamp;<-:create;"`
	UpdatedAt       time.Time      `gorm:"type:timestamp;"`
}

func (m *Saga) DTO() dto.Saga {
	return dto.Saga{
		ID:              m.ID,
		CorrelationID:   m.CorrelationID,
		TripID:          m.TripID,
		CarID:           m.CarID,
		HotelID:         m.HotelID,
		FlightID:        m.FlightID,
		CarBookingID:    m.CarBookingID,
		HotelBookingID:  m.HotelBookingID,
		FlightBookingID: m.FlightBookingID,
		Status:          m.Status,
		History:         string(m.History),
		CreatedAt:       m.CreatedAt,
	}
}
