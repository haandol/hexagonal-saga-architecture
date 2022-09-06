package entity

import (
	"time"

	"github.com/haandol/hexagonal/pkg/dto"
)

type CarBooking struct {
	ID        uint      `gorm:"type:bigint;primaryKey;autoIncrement;<-:create;"`
	TripID    uint      `gorm:"type:bigint;not null;"`
	CarID     uint      `gorm:"type:bigint;not null;"`
	Status    string    `gorm:"type:varchar(16);not null;"`
	CreatedAt time.Time `gorm:"type:timestamp;<-:create;"`
	UpdatedAt time.Time `gorm:"type:timestamp;"`
}

func (m CarBooking) DTO() (dto.CarBooking, error) {
	return dto.CarBooking{
		ID:        m.ID,
		TripID:    m.TripID,
		CarID:     m.CarID,
		Status:    m.Status,
		CreatedAt: m.CreatedAt,
	}, nil
}
