package status

// for each service
const (
	Initialized = "Initialized"
	Booked      = "Booked"
	Cancelled   = "Cancelled"
)

// for trip
const (
	TripInitialized = "Initialized"
	TripBooked      = "Booked"
	TripCancelled   = "Cancelled"
)

// for saga
const (
	SagaStarted = "Started"
	SagaEnded   = "Ended"
	SagaAborted = "Aborted"
)
