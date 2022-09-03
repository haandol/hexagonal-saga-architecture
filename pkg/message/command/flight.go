package command

import "github.com/haandol/hexagonal/pkg/message"

type BookFlight struct {
	message.Message
	Body BookFlightBody `json:"body"`
}

type BookFlightBody struct {
	TripID   uint `json:"tripId"`
	FlightID uint `json:"flightId"`
}
