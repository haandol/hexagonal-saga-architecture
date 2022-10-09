package dto

import "time"

type Outbox struct {
	ID         uint      `json:"id"`
	KafkaTopic string    `json:"kafkaTopic" binding:"required" validate:"required"`
	KafkaKey   string    `json:"kafkaKey" binding:"required" validate:"required"`
	KafkaValue string    `json:"kafkaValue" binding:"required" validate:"required,json"`
	IsSent     bool      `json:"isSent"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}
