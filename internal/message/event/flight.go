package event

import "github.com/haandol/hexagonal/internal/message"

type FlightBooked struct {
	message.Message
	Body FlightBookedBody `json:"body" validate:"required"`
}

type FlightBookedBody struct {
	BookingID uint `json:"bookingId" validate:"required"`
}

type FlightBookingCanceledBody struct {
	BookingID uint `json:"bookingId" validate:"required"`
	TripID    uint `json:"tripId" validate:"required"`
}

type FlightBookingCanceled struct {
	message.Message
	Body FlightBookingCanceledBody `json:"body" validate:"required"`
}
