package event

import (
	"github.com/haandol/hexagonal/internal/message"
)

type CarBooked struct {
	message.Message
	Body CarBookedBody `json:"body" validate:"required"`
}

type CarBookedBody struct {
	BookingID uint `json:"bookingId" validate:"required"`
}

type CarBookingCanceled struct {
	message.Message
	Body CarBookingCanceledBody `json:"body" validate:"required"`
}

type CarBookingCanceledBody struct {
	BookingID uint `json:"bookingId" validate:"required"`
	TripID    uint `json:"tripId" validate:"required"`
}
