package main

import (
	"encoding/json"
	"net/http"
	"sync"
	"log"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // allow all origins for dev
	},
}

type WSMessage struct {
	RequestCount int    `json:"requestCount"`
	Route        string `json:"route"`
	CircuitState string `json:"circuitState"`
	
}

var clients = make(map[*websocket.Conn]bool)
var mu sync.Mutex

func wsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("WebSocket connection received")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	mu.Lock()
	clients[conn] = true
	mu.Unlock()

	for {
		// keep connection alive
		_, _, err := conn.ReadMessage()
		if err != nil {
			mu.Lock()
			delete(clients, conn)
			mu.Unlock()
			conn.Close()
			break
		}
	}
}
func broadcastMetrics() {
	msg := WSMessage{
		RequestCount: metrics.RequestCount,
		Route:        metrics.ActiveRoute,
		CircuitState: metrics.CircuitState,
		
	}

	data, _ := json.Marshal(msg)

	mu.Lock()
	defer mu.Unlock()

	for conn := range clients {
		err := conn.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			conn.Close()
			delete(clients, conn)
		}
	}
}