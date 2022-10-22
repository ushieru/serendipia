package circuitbreaker

type CircuitBreakerStates string

const (
	Close    CircuitBreakerStates = "Close"
	HalfOpen CircuitBreakerStates = "HalfOpen"
	Open     CircuitBreakerStates = "Open"
)
