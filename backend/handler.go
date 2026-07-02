package main

import (	
	"net/http"
	"log"
)
func syncMetrics() {	
	metrics.SetCircuitState(string(cb.GetState()))
}
func helloHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("helloHandler called")
	metrics.IncrementRequest()
	broadcastMetrics() 
	proxyRequest(w)
}