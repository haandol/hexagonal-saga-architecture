package entity

import (
	"time"

	"github.com/haandol/hexagonal/pkg/dto"
)

type CarBooking struct {
	ID        uint      `gorm:"type:bigint;primaryKey;autoIncrement;<-:create;"`
	TripID    uint      `gorm:"type:bigint;not null;"`
	CarID     uint      `gorm:"type:bigint;not null;"`
	CreatedAt time.Time `gorm:"type:timestamp;<-:create;"`
	UpdatedAt time.Time `gorm:"type:timestamp;"`
}

func (m CarBooking) DTO() (dto.CarBooking, error) {
	return dto.CarBooking{
		ID:        m.ID,
		TripID:    m.TripID,
		CarID:     m.CarID,
		CreatedAt: m.CreatedAt,
	}, nil
}
