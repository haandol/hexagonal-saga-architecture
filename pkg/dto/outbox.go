package dto

import "time"

type Outbox struct {
	ID         uint      `json:"id"`
	KafkaTopic string    `json:"kafka_topic" binding:"required" validate:"required"`
	KafkaKey   string    `json:"kafka_key" binding:"required" validate:"required"`
	KafkaValue string    `json:"kafka_value" binding:"required" validate:"required,json"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}
