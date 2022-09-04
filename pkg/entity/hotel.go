package entity

import (
	"time"

	"github.com/haandol/hexagonal/pkg/dto"
)

type HotelBooking struct {
	ID        uint      `gorm:"type:bigint;primaryKey;autoIncrement;<-:create;"`
	TripID    uint      `gorm:"type:bigint;not null;"`
	HotelID   uint      `gorm:"type:bigint;not null;"`
	CreatedAt time.Time `gorm:"type:timestamp;<-:create;"`
	UpdatedAt time.Time `gorm:"type:timestamp;"`
}

func (m HotelBooking) DTO() (dto.HotelBooking, error) {
	return dto.HotelBooking{
		ID:        m.ID,
		TripID:    m.TripID,
		HotelID:   m.HotelID,
		CreatedAt: m.CreatedAt,
	}, nil
}
