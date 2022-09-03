package command

import "github.com/haandol/hexagonal/pkg/message"

type StartSaga struct {
	message.Message
	Body StartSagaBody `json:"body"`
}

type StartSagaBody struct {
	TripID   uint `json:"tripId"`
	CarID    uint `json:"carId"`
	HotelID  uint `json:"hotelId"`
	FlightID uint `json:"flightId"`
}

type EndSaga struct {
	message.Message
	Body EndSagaBody `json:"body"`
}

type EndSagaBody struct {
	SagaID uint `json:"sagaId"`
}

type AbortSaga struct {
	message.Message
	Body AbortSagaBody `json:"body"`
}

type AbortSagaBody struct {
	SagaID uint `json:"sagaId"`
}
