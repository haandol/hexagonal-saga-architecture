package event

import "github.com/haandol/hexagonal/pkg/message"

type FlightBooked struct {
	message.Message
	Body FlightBookedBody `json:"body" validate:"required"`
}

type FlightBookedBody struct {
	BookingID uint `json:"bookingId" validate:"required"`
}

type FlightBookingCancelledBody struct {
	BookingID uint `json:"bookingId" validate:"required"`
	TripID    uint `json:"tripId" validate:"required"`
}

type FlightBookingCancelled struct {
	message.Message
	Body FlightBookingCancelledBody `json:"body" validate:"required"`
}
