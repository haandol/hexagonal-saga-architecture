package command

import "github.com/haandol/hexagonal/pkg/message"

type BookHotel struct {
	message.Message
	Body BookHotelBody `json:"body"`
}

type BookHotelBody struct {
	TripID  uint `json:"tripId"`
	HotelID uint `json:"hotelId"`
}
