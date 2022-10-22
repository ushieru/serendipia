package circuitbreaker

type CircuitBreakerState struct {
	Failures       int64
	CooldownPeriod int64
	State          CircuitBreakerStates
	NextTry        int64
}
