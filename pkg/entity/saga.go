package entity

import (
	"time"

	"github.com/haandol/hexagonal/pkg/dto"
	"gorm.io/datatypes"
)

type Saga struct {
	ID            uint           `gorm:"type:bigint;primaryKey;autoIncrement;<-:create;"`
	CorrelationID string         `gorm:"type:varchar(36);<-:create;"`
	TripID        uint           `gorm:"type:bigint;<-:create;"`
	CarID         uint           `gorm:"type:bigint;"`
	HotelID       uint           `gorm:"type:bigint;"`
	FlightID      uint           `gorm:"type:bigint;"`
	Status        string         `gorm:"type:varchar(16);"`
	History       datatypes.JSON `gorm:"type:json;"`
	CreatedAt     time.Time      `gorm:"type:timestamp;<-:create;"`
	UpdatedAt     time.Time      `gorm:"type:timestamp;"`
}

func (m Saga) DTO() (dto.Saga, error) {
	return dto.Saga{
		ID:            m.ID,
		CorrelationID: m.CorrelationID,
		TripID:        m.TripID,
		CarID:         m.CarID,
		HotelID:       m.HotelID,
		FlightID:      m.FlightID,
		Status:        m.Status,
		History:       string(m.History),
		CreatedAt:     m.CreatedAt,
	}, nil
}
