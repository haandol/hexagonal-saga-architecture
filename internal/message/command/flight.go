package command

import "github.com/haandol/hexagonal/internal/message"

type BookFlight struct {
	message.Message
	Body BookFlightBody `json:"body" validate:"required"`
}

type BookFlightBody struct {
	TripID   uint `json:"tripId" validate:"required"`
	FlightID uint `json:"flightId" validate:"required"`
}

type CancelFlightBooking struct {
	message.Message
	Body CancelFlightBookingBody `json:"body" validate:"required"`
}

type CancelFlightBookingBody struct {
	TripID    uint `json:"tripId" validate:"required"`
	BookingID uint `json:"bookingID" validate:"required"`
}
