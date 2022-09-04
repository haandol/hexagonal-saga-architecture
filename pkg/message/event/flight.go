package event

import "github.com/haandol/hexagonal/pkg/message"

type FlightBooked struct {
	message.Message
	Body FlightBookedBody `json:"body" validate:"required"`
}

type FlightBookedBody struct {
	BookingID uint `json:"bookingId" validate:"required"`
}

type FlightBookingCanceled struct {
	message.Message
	Body FlightBookedBody `json:"body" validate:"required"`
}
