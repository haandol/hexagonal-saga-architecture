package event

import "github.com/haandol/hexagonal/pkg/message"

type HotelBooked struct {
	message.Message
	Body HotelBookedBody `json:"body" validate:"required"`
}

type HotelBookedBody struct {
	BookingID uint `json:"bookingId" validate:"required"`
}

type HotelBookingCancelled struct {
	message.Message
	Body HotelBookingCancelledBody `json:"body" validate:"required"`
}

type HotelBookingCancelledBody struct {
	BookingID uint `json:"bookingId" validate:"required"`
	TripID    uint `json:"tripId" validate:"required"`
}
