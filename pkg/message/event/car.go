package event

import (
	"github.com/haandol/hexagonal/pkg/message"
)

type CarBooked struct {
	message.Message
	Body CarBookedBody `json:"body" validate:"required"`
}

type CarBookedBody struct {
	BookingID uint `json:"bookingId" validate:"required"`
}

type CarBookingCancelled struct {
	message.Message
	Body CarBookingCancelledBody `json:"body" validate:"required"`
}

type CarBookingCancelledBody struct {
	BookingID uint `json:"bookingId" validate:"required"`
	TripID    uint `json:"tripId" validate:"required"`
}
