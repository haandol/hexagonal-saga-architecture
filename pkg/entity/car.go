package entity

import (
	"time"

	"github.com/haandol/hexagonal/pkg/dto"
)

type CarRental struct {
	ID        uint      `gorm:"type:bigint;primaryKey;autoIncrement;<-:create;"`
	TripID    uint      `gorm:"type:bigint;not null;"`
	CarID     uint      `gorm:"type:bigint;not null;"`
	Quantity  int64     `gorm:"type:int;not null;"`
	CreatedAt time.Time `gorm:"type:timestamp;<-:create;"`
	UpdatedAt time.Time `gorm:"type:timestamp;"`
}

func (m CarRental) DTO() (dto.CarRental, error) {
	return dto.CarRental{
		ID:        m.ID,
		TripID:    m.TripID,
		CarID:     m.CarID,
		Quantity:  m.Quantity,
		CreatedAt: m.CreatedAt,
	}, nil
}
