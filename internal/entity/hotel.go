package entity

import (
	"time"

	"github.com/haandol/hexagonal/internal/dto"
)

type HotelBooking struct {
	ID        uint      `gorm:"type:bigint;primaryKey;autoIncrement;<-:create;"`
	TripID    uint      `gorm:"type:bigint;not null;"`
	HotelID   uint      `gorm:"type:bigint;not null;"`
	Status    string    `gorm:"type:varchar(32);not null;"`
	CreatedAt time.Time `gorm:"type:timestamp;<-:create;"`
	UpdatedAt time.Time `gorm:"type:timestamp;"`
}

func (m *HotelBooking) DTO() dto.HotelBooking {
	return dto.HotelBooking{
		ID:        m.ID,
		TripID:    m.TripID,
		HotelID:   m.HotelID,
		Status:    m.Status,
		CreatedAt: m.CreatedAt,
	}
}
