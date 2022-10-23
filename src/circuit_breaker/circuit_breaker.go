package circuitbreaker

import (
	"errors"
	"io"
	"net/http"

	"github.com/ushieru/serendipia/src/utils"
)

type CircuitBreaker struct {
	States           map[string]CircuitBreakerState
	FailureThreshold int64
	CooldownPeriod   int64
	RequestTimeout   int64
	client           http.Client
}

func NewCircuitBreaker(failureThreshold int64, cooldownPeriod int64, requestTimeout int64) *CircuitBreaker {
	return &CircuitBreaker{
		States:           make(map[string]CircuitBreakerState),
		FailureThreshold: failureThreshold,
		CooldownPeriod:   cooldownPeriod,
		RequestTimeout:   requestTimeout,
		client:           http.Client{},
	}
}

func (circuitBreaker *CircuitBreaker) CallService(method string, url string, body io.Reader, headers map[string]string) (*http.Response, error) {
	key := method + ":" + url
	if !circuitBreaker.canRequest(key) {
		return nil, errors.New("service open")
	}
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	for key, value := range headers {
		request.Header.Add(key, value)
	}
	return circuitBreaker.client.Do(request)
}

func (circuitBreaker *CircuitBreaker) canRequest(endpoint string) bool {
	if _, ok := circuitBreaker.States[endpoint]; !ok {
		circuitBreaker.initState(endpoint)
	}
	state := circuitBreaker.States[endpoint]
	if state.State == Close {
		return true
	}
	if utils.GetTimeStamp() >= state.NextTry {
		state.State = HalfOpen
		return true
	}
	return false
}

func (circuitBreaker *CircuitBreaker) initState(endpoint string) {
	circuitBreaker.States[endpoint] = CircuitBreakerState{
		Failures:       0,
		CooldownPeriod: circuitBreaker.CooldownPeriod,
		State:          Close,
		NextTry:        0,
	}
}
