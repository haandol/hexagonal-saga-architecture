package command

import "github.com/haandol/hexagonal/pkg/message"

type RentCar struct {
	message.Message
	Body RentCarBody `json:"body" validate:"required"`
}

type RentCarBody struct {
	TripID   uint  `json:"tripId" validate:"required"`
	CarID    uint  `json:"carId" validate:"required"`
	Quantity int64 `json:"quantity"`
}

type CancelCarRental struct {
	message.Message
	Body CancelCarRentalBody `json:"body" validate:"required"`
}

type CancelCarRentalBody struct {
	RentalID uint `json:"rentalId" validate:"required"`
}
