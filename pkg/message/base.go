package message

type Message struct {
	Name          string `json:"name" validate:"required"`
	Version       string `json:"version" validate:"required"`
	ID            string `json:"id" validate:"required"`
	CorrelationID string `json:"correlationId" validate:"required"`
	ParentID      string `json:"parentId"`
	CreatedAt     string `json:"createdAt" validate:"required"`
}
