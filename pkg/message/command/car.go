package command

import "github.com/haandol/hexagonal/pkg/message"

type RentCar struct {
	message.Message
	Body RentCarBody `json:"body"`
}

type RentCarBody struct {
	TripID uint `json:"tripId"`
	CarID  uint `json:"carId"`
}
