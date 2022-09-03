package event

import "github.com/haandol/hexagonal/pkg/message"

type SagaEnded struct {
	message.Message
	Body SagaEndedBody `json:"body"`
}

type SagaEndedBody struct {
	SagaID uint `json:"sagaId"`
}

type SagaAborted struct {
	message.Message
	Body SagaAbortedBody `json:"body"`
}

type SagaAbortedBody struct {
	SagaID uint `json:"sagaId"`
}
