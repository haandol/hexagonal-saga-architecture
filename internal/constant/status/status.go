package status

// for each service
const (
	Initialized = "Initialized"
	Booked      = "Booked"
	Canceled    = "Canceled"
)

// for trip
const (
	TripInitialized = "Initialized"
	TripBooked      = "Booked"
	TripCanceled    = "Canceled"
)

// for saga
const (
	SagaStarted = "Started"
	SagaEnded   = "Ended"
	SagaAborted = "Aborted"
)
