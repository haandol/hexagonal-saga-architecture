package entity

import (
	"time"

	"github.com/haandol/hexagonal/pkg/dto"
	"github.com/haandol/hexagonal/pkg/util"
)

type Trip struct {
	ID        uint      `gorm:"type:bigint;primaryKey;autoIncrement;<-:create;"`
	UserID    uint      `gorm:"type:bigint;<-:create;"`
	HotelID   uint      `gorm:"type:bigint;<-:create;"`
	CarID     uint      `gorm:"type:bigint;<-:create;"`
	FlightID  uint      `gorm:"type:bigint;<-:create;"`
	Status    string    `gorm:"type:varchar(16);"`
	CreatedAt time.Time `gorm:"type:timestamp;<-:create;"`
	UpdatedAt time.Time `gorm:"type:timestamp;<-:create;"`
}

type Trips []Trip

func (m Trip) DTO() (dto.Trip, error) {
	return dto.Trip{
		ID:        m.ID,
		UserID:    m.UserID,
		HotelID:   m.HotelID,
		CarID:     m.CarID,
		FlightID:  m.FlightID,
		Status:    m.Status,
		CreatedAt: m.CreatedAt,
	}, nil
}

func (m Trips) DTO() ([]dto.Trip, error) {
	logger := util.GetLogger()

	trips := make([]dto.Trip, 0)
	for _, trip := range m {
		t, err := trip.DTO()
		if err != nil {
			logger.Error(err.Error())
			continue
		}
		trips = append(trips, t)
	}
	return trips, nil
}
