package dto

import "time"

type Trip struct {
	ID              uint      `json:"id"`
	UserID          uint      `json:"userId" binding:"required" validate:"required"`
	CarID           uint      `json:"carId" binding:"required" validate:"required"`
	HotelID         uint      `json:"hotelId" binding:"required" validate:"required"`
	FlightID        uint      `json:"flightId" binding:"required" validate:"required"`
	CarBookingID    uint      `json:"carBookingId"`
	HotelBookingID  uint      `json:"hotelBookingId"`
	FlightBookingID uint      `json:"flightBookingId"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}
