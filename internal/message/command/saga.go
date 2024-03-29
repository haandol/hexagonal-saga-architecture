package command

import "github.com/haandol/hexagonal/internal/message"

type StartSaga struct {
	message.Message
	Body StartSagaBody `json:"body" validate:"required"`
}

type StartSagaBody struct {
	TripID   uint `json:"tripId" validate:"required"`
	CarID    uint `json:"carId" validate:"required"`
	HotelID  uint `json:"hotelId" validate:"required"`
	FlightID uint `json:"flightId" validate:"required"`
}

type EndSaga struct {
	message.Message
	Body EndSagaBody `json:"body" validate:"required"`
}

type EndSagaBody struct {
	SagaID uint `json:"sagaId" validate:"required"`
	TripID uint `json:"tripId" validate:"required"`
}

type AbortSaga struct {
	message.Message
	Body AbortSagaBody `json:"body" validate:"required"`
}

type AbortSagaBody struct {
	TripID uint   `json:"tripId" validate:"required"`
	Reason string `json:"reason" validate:"required"`
	Source string `json:"source" validate:"required"`
}
