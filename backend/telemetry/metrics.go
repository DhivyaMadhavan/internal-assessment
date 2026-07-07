package telemetry

import (
    "sync"
    "time"
)

type Metrics struct {
	mu sync.Mutex

	RequestCount int    `json:"requestCount"`
	RPS          int    `json:"rps"`
	ActiveRoute  string `json:"route"`
	CircuitState string `json:"circuitState"`
	currentSecond int
}

func NewMetrics() *Metrics {

	return &Metrics{
		RequestCount: 0,
		ActiveRoute: "N/A",
		CircuitState: "CLOSED",
	}
}

func (m *Metrics) IncrementRequest() {
    m.mu.Lock()
    defer m.mu.Unlock()

    m.RequestCount++
    m.currentSecond++
}

func (m *Metrics) SetRoute(route string) {
	m.ActiveRoute = route
}

func (m *Metrics) SetCircuitState(state string) {
	m.CircuitState = state
}

func (m *Metrics) StartRPSCounter(callback func()) {

    ticker := time.NewTicker(1 * time.Second)

    go func() {
        for range ticker.C {

            m.mu.Lock()

            m.RPS = m.currentSecond
            m.currentSecond = 0

            m.mu.Unlock()

            callback()
        }
    }()
}

func (m *Metrics) GetRequestCount() int {
    m.mu.Lock()
    defer m.mu.Unlock()
    return m.RequestCount
}

func (m *Metrics) GetRPS() int {
    m.mu.Lock()
    defer m.mu.Unlock()
    return m.RPS
}

func (m *Metrics) GetRoute() string {
    m.mu.Lock()
    defer m.mu.Unlock()
    return m.ActiveRoute
}

func (m *Metrics) GetCircuitState() string {
    m.mu.Lock()
    defer m.mu.Unlock()
    return m.CircuitState
}