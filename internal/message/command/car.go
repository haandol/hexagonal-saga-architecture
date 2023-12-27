package command

import "github.com/haandol/hexagonal/internal/message"

type BookCar struct {
	message.Message
	Body BookCarBody `json:"body" validate:"required"`
}

type BookCarBody struct {
	TripID uint `json:"tripId" validate:"required"`
	CarID  uint `json:"carId" validate:"required"`
}

type CancelCarBooking struct {
	message.Message
	Body CancelCarBookingBody `json:"body" validate:"required"`
}

type CancelCarBookingBody struct {
	TripID    uint `json:"tripId" validate:"required"`
	BookingID uint `json:"bookingId" validate:"required"`
}
