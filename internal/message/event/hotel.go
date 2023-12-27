package event

import "github.com/haandol/hexagonal/internal/message"

type HotelBooked struct {
	message.Message
	Body HotelBookedBody `json:"body" validate:"required"`
}

type HotelBookedBody struct {
	BookingID uint `json:"bookingId" validate:"required"`
}

type HotelBookingCanceled struct {
	message.Message
	Body HotelBookingCanceledBody `json:"body" validate:"required"`
}

type HotelBookingCanceledBody struct {
	BookingID uint `json:"bookingId" validate:"required"`
	TripID    uint `json:"tripId" validate:"required"`
}
