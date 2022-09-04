package dto

import "time"

type HotelBooking struct {
	ID        uint      `json:"id"`
	TripID    uint      `json:"tripId" binding:"required" validate:"required"`
	HotelID   uint      `json:"hotelId" binding:"required" validate:"required"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
