package entity

import (
	"time"

	"github.com/haandol/hexagonal/pkg/dto"
)

type Trip struct {
	ID              uint      `gorm:"type:bigint;primaryKey;autoIncrement;<-:create;"`
	UserID          uint      `gorm:"type:bigint;<-:create;"`
	CarID           uint      `gorm:"type:bigint;"`
	HotelID         uint      `gorm:"type:bigint;"`
	FlightID        uint      `gorm:"type:bigint;"`
	CarBookingID    uint      `gorm:"type:bigint;"`
	HotelBookingID  uint      `gorm:"type:bigint;"`
	FlightBookingID uint      `gorm:"type:bigint;"`
	Status          string    `gorm:"type:varchar(32);"`
	CreatedAt       time.Time `gorm:"type:timestamp;<-:create;"`
	UpdatedAt       time.Time `gorm:"type:timestamp;"`
}

type Trips []Trip

func (m Trip) DTO() dto.Trip {
	return dto.Trip{
		ID:              m.ID,
		UserID:          m.UserID,
		CarID:           m.CarID,
		HotelID:         m.HotelID,
		FlightID:        m.FlightID,
		CarBookingID:    m.CarBookingID,
		HotelBookingID:  m.HotelBookingID,
		FlightBookingID: m.FlightBookingID,
		Status:          m.Status,
		CreatedAt:       m.CreatedAt,
	}
}

func (m Trips) DTO() []dto.Trip {
	trips := make([]dto.Trip, 0)
	for _, trip := range m {
		trips = append(trips, trip.DTO())
	}
	return trips
}
