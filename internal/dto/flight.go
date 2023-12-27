package dto

import "time"

type FlightBooking struct {
	ID        uint      `json:"id"`
	TripID    uint      `json:"tripId" binding:"required" validate:"required"`
	FlightID  uint      `json:"flightId" binding:"required" validate:"required"`
	Status    string    `json:"status" binding:"required" validate:"required"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
