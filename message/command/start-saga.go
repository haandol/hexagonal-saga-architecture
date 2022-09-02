package command

type StartSaga struct {
	Command
	Body StartSagaBody `json:"body"`
}

type StartSagaBody struct {
	TripID uint `json:"tripId"`
}
