package event

import "github.com/haandol/hexagonal/pkg/message"

type HotelBooked struct {
	message.Message
	Body HotelBookedBody `json:"body" validate:"required"`
}

type HotelBookedBody struct {
	BookingID uint `json:"bookingId" validate:"required"`
}

type HotelBookingCanceled struct {
	message.Message
	Body HotelBookedBody `json:"body" validate:"required"`
}
