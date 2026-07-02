package telemetry

type Metrics struct {
	RequestCount int    `json:"requestCount"`
	ActiveRoute  string `json:"route"`
	CircuitState string `json:"circuitState"`
}

func NewMetrics() *Metrics {

	return &Metrics{
		RequestCount: 0,
		ActiveRoute: "",
		CircuitState: "CLOSED",
	}
}

func (m *Metrics) IncrementRequest() {
	m.RequestCount++
}

func (m *Metrics) SetRoute(route string) {
	m.ActiveRoute = route
}

func (m *Metrics) SetCircuitState(state string) {
	m.CircuitState = state
}