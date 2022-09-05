package event

import "github.com/haandol/hexagonal/pkg/message"

type SagaEnded struct {
	message.Message
	Body SagaEndedBody `json:"body" validate:"required"`
}

type SagaEndedBody struct {
	SagaID uint `json:"sagaId" validate:"required"`
}

type SagaAborted struct {
	message.Message
	Body SagaAbortedBody `json:"body" validate:"required"`
}

type SagaAbortedBody struct {
	SagaID uint   `json:"sagaId" validate:"required"`
	Reason string `json:"reason" validate:"required"`
	Source string `json:"source" validate:"required"`
}
