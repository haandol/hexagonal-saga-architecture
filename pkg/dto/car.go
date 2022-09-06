package dto

import "time"

type CarBooking struct {
	ID        uint      `json:"id"`
	TripID    uint      `json:"tripId" binding:"required" validate:"required"`
	CarID     uint      `json:"carId" binding:"required" validate:"required"`
	Status    string    `json:"status" binding:"required" validate:"required"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
