package dto

import "time"

type Saga struct {
	ID            uint      `json:"id"`
	CorrelationID string    `json:"correlationId" binding:"required" validate:"required"`
	TripID        uint      `json:"tripId" binding:"required" validate:"required"`
	CarID         uint      `json:"carId" binding:"required" validate:"required"`
	HotelID       uint      `json:"hotelId" binding:"required" validate:"required"`
	FlightID      uint      `json:"flightId" binding:"required" validate:"required"`
	Status        string    `json:"status" binding:"required" validate:"required"`
	History       string    `json:"history"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}
