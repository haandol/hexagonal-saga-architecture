package event

import "github.com/haandol/hexagonal/internal/message"

type SagaEnded struct {
	message.Message
	Body SagaEndedBody `json:"body" validate:"required"`
}

type SagaEndedBody struct {
	SagaID          uint `json:"sagaId" validate:"required"`
	TripID          uint `json:"tripId" validate:"required"`
	CarBookingID    uint `json:"carBookingId" validate:"required"`
	HotelBookingID  uint `json:"hotelBookingId" validate:"required"`
	FlightBookingID uint `json:"flightBookingId" validate:"required"`
}

type SagaAborted struct {
	message.Message
	Body SagaAbortedBody `json:"body" validate:"required"`
}

type SagaAbortedBody struct {
	SagaID uint `json:"sagaId" validate:"required"`
	TripID uint `json:"tripId" validate:"required"`
}
