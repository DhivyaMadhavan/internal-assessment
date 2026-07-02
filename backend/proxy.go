package main

import (
	"context"	
	"io"
	"net/http"
	"time"

	"internal-assessment/backend/circuitbreaker"
	"internal-assessment/backend/telemetry"
)

var cb = circuitbreaker.NewCircuitBreaker()
var metrics = telemetry.NewMetrics()

func proxyRequest(w http.ResponseWriter) {
	
	// Remember the circuit state before checking
	previousState := cb.GetState()

	// Check if request is allowed
	allowed := cb.AllowRequest() // checks if the timeout is expired and then can the state be moved to HALF_OPEN

	// If the circuit changed (for example OPEN -> HALF_OPEN),
	// immediately notify the frontend.
	if previousState != cb.GetState() {
		syncMetrics()
		broadcastMetrics()		
	}

	// Circuit is still OPEN
	if !allowed {		

		metrics.SetRoute("Secondary")
		syncMetrics()
		broadcastMetrics()		

		response, err := http.Get("http://secondary-api:8082/hello")

		if err != nil {
			http.Error(w, "Secondary API unavailable", http.StatusInternalServerError)
			return
		}

		defer response.Body.Close()

		body, err := io.ReadAll(response.Body)

		if err != nil {
			http.Error(w, "Unable to read Secondary API response", http.StatusInternalServerError)
			return
		}	
		w.Write(body)	
		return
	}

	// Create timeout (200 ms)
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"http://toxiproxy:8666/hello",
		nil,
	)

	if err != nil {
		http.Error(w, "Unable to create request", http.StatusInternalServerError)
		return
	}

	client := &http.Client{}

	response, err := client.Do(req) 

	// -----------------------------
	// Primary API Success
	// -----------------------------
	if err == nil {

		cb.RecordSuccess()
		metrics.SetRoute("Primary")
		syncMetrics()
		broadcastMetrics()	
		
		defer response.Body.Close()

		body, err := io.ReadAll(response.Body)

		if err != nil {
			http.Error(w, "Unable to read Primary API response", http.StatusInternalServerError)
			return
		}
		w.Write(body)
		
		return
	}

	// -----------------------------
	// Primary API Failed
	// -----------------------------
	cb.RecordFailure()
	metrics.SetRoute("Secondary")
	syncMetrics()
	broadcastMetrics()
	
	// Call Secondary API
	response, err = http.Get("http://secondary-api:8082/hello")

	if err != nil {
		http.Error(w, "Both APIs are unavailable", http.StatusInternalServerError)
		return
	}

	defer response.Body.Close()

	// FIX: Changed from '=' to ':=' because 'body' is not defined in this scope
	body, err := io.ReadAll(response.Body)

	if err != nil {
		http.Error(w, "Unable to read Secondary API response", http.StatusInternalServerError)
		return
	}
	w.Write(body)

	
}
