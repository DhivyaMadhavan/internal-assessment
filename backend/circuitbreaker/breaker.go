package circuitbreaker

import (
	
	"sync"
	"time"
)

type State string

const (
	Closed   State = "CLOSED"
	Open     State = "OPEN"
	HalfOpen State = "HALF_OPEN"
)

type CircuitBreaker struct {
	mu sync.Mutex

	State            State
	FailureCount     int
	FailureThreshold int
	OpenedAt         time.Time
	Transitions      []string
}

func NewCircuitBreaker() *CircuitBreaker {
	return &CircuitBreaker{
		State:            Closed,
		FailureThreshold: 3,
	}
}

// Called whenever Primary API fails
func (cb *CircuitBreaker) RecordFailure() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	// Half-Open test failed
	if cb.State == HalfOpen {		

		oldState := cb.State
		cb.State = Open
		cb.OpenedAt = time.Now()

		cb.logTransition(oldState, cb.State)

		return
	}

	cb.FailureCount++

	if cb.State == Closed &&
		cb.FailureCount >= cb.FailureThreshold {

		oldState := cb.State
		cb.State = Open
		cb.OpenedAt = time.Now()

		cb.logTransition(oldState, cb.State)
	}
}

// Called whenever Primary API succeeds
func (cb *CircuitBreaker) RecordSuccess() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.FailureCount = 0

	// Only transition when recovering from HALF_OPEN
	if cb.State == HalfOpen {	

		oldState := cb.State
		cb.State = Closed

		cb.logTransition(oldState, cb.State)
	}
}

// Decide whether Primary API should receive the request
func (cb *CircuitBreaker) AllowRequest() bool {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	switch cb.State {

	case Closed:
		return true

	case Open:

		// Wait before testing again
		if time.Since(cb.OpenedAt) > 5*time.Second {		

			oldState := cb.State
			cb.State = HalfOpen
			cb.logTransition(oldState, cb.State)

			return true
		}

		return false

	case HalfOpen:
		return false
	}

	return false
}

// Helper to log only real transitions
func (cb *CircuitBreaker) logTransition(oldState, newState State) {
	// Prevent duplicate entries if the state didn't actually change
	if oldState == newState {
		return
	}

	// COMPRESSED: Appends the transition string directly to the log slice
	cb.Transitions = append(cb.Transitions, string(oldState)+" -> "+string(newState))
}


func (cb *CircuitBreaker) GetState() State {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	return cb.State
}