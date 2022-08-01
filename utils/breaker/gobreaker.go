package breaker

import (
	"github.com/sony/gobreaker"
	"log"
	"time"
)

func NewCircuitBreaker() *gobreaker.CircuitBreaker {
	return gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:    "ORG BREAKER",
		Timeout: time.Second * 5,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.Requests >= 3 && failureRatio >= 0.6
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			if to == gobreaker.StateOpen {
				log.Println("State Open!")
			}
			if from == gobreaker.StateOpen && to == gobreaker.StateHalfOpen {
				log.Println("Going from Open to Half Open")
			}
			if from == gobreaker.StateHalfOpen && to == gobreaker.StateClosed {
				log.Println("Going from Half Open to Closed!")
			}
		},
	})
}
