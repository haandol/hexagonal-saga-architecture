package command

import "github.com/haandol/hexagonal/internal/message"

type BookHotel struct {
	message.Message
	Body BookHotelBody `json:"body" validate:"required"`
}

type BookHotelBody struct {
	TripID  uint `json:"tripId" validate:"required"`
	HotelID uint `json:"hotelId" validate:"required"`
}

type CancelHotelBooking struct {
	message.Message
	Body CancelHotelBookingBody `json:"body" validate:"required"`
}

type CancelHotelBookingBody struct {
	TripID    uint `json:"tripId" validate:"required"`
	BookingID uint `json:"bookingID" validate:"required"`
}
