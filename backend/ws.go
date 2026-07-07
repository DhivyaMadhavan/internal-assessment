package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // dev only
	},
}

type WSMessage struct {
	RequestCount int    `json:"requestCount"`
	RPS          int    `json:"rps"`
	Route        string `json:"route"`
	CircuitState string `json:"circuitState"`
}

var (
	clients = make(map[*websocket.Conn]bool)
	mu      sync.Mutex
)

// 🔥 WebSocket handler
func wsHandler(w http.ResponseWriter, r *http.Request) {

	log.Println("WebSocket connection received")

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	mu.Lock()
	clients[conn] = true
	mu.Unlock()

	// Keep connection alive (read loop)
	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			log.Println("Client disconnected")

			mu.Lock()
			delete(clients, conn)
			mu.Unlock()

			conn.Close()
			break
		}
	}
}

// 🔥 SAFE BROADCAST (THIS IS THE IMPORTANT FIX)
func broadcastMetrics() {

	msg := WSMessage{
		RequestCount: metrics.RequestCount,
		RPS:          metrics.GetRPS(),
		Route:        metrics.ActiveRoute,
		CircuitState: metrics.CircuitState,
	}

	data, err := json.Marshal(msg)
	if err != nil {
		log.Println("JSON marshal error:", err)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	// if no clients → do nothing
	if len(clients) == 0 {
		return
	}
	log.Printf("Broadcasting: %+v\n", msg)
	for conn := range clients {

		// 🔥 IMPORTANT: each write should not crash broadcast loop
		err := conn.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			log.Println("Write error, removing client:", err)

			conn.Close()
			delete(clients, conn)
		}
	}
}