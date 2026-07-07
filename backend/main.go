package main

import (
	"log"
	"net/http"
)

func main() {

	metrics.StartRPSCounter(func() {
        broadcastMetrics()
    })


	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/ws", wsHandler)

	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Fatal(err)
	}
}