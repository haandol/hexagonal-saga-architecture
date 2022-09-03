package dto

import "time"

type Car struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name" binding:"required" validate:"required"`
	Quantity  int64     `json:"quantity" binding:"required" validate:"required"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type CarRental struct {
	ID        uint      `json:"id"`
	TripID    uint      `json:"tripId" binding:"required" validate:"required"`
	CarID     uint      `json:"carId" binding:"required" validate:"required"`
	Quantity  int64     `json:"quantity" binding:"required" validate:"required"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
