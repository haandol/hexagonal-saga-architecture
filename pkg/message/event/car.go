package event

import "github.com/haandol/hexagonal/pkg/message"

type CarRented struct {
	message.Message
	Body CarRentedBody `json:"body"`
}

type CarRentedBody struct {
	CarRentalID uint `json:"carRentalId"`
}

type CarRentalCanceled struct {
	message.Message
	Body CarRentedBody `json:"body"`
}
