package dto

import "time"

type Trip struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"userId" binding:"required" validate:"required"`
	CarID     uint      `json:"carId" binding:"required" validate:"required"`
	HotelID   uint      `json:"hotelId" binding:"required" validate:"required"`
	FlightID  uint      `json:"flightId" binding:"required" validate:"required"`
	Status    string    `json:"status" binding:"required" validate:"required"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
