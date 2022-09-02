package command

type Command struct {
	Name          string `json:"name" validate:"required"`
	Version       string `json:"version" validate:"required"`
	ID            string `json:"id" validate:"required"`
	CorrelationID string `json:"correlationId" validate:"required"`
	CreatedAt     string `json:"createdAt" validate:"required"`
}
