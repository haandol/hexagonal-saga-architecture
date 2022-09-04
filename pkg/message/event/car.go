package event

import "github.com/haandol/hexagonal/pkg/message"

type CarBooked struct {
	message.Message
	Body CarBookedBody `json:"body" validate:"required"`
}

type CarBookedBody struct {
	BookingID uint `json:"bookingId" validate:"required"`
}

type CarBookingCanceled struct {
	message.Message
	Body CarBookedBody `json:"body" validate:"required"`
}