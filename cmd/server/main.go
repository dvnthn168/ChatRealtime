package main

import (
	"log"
	"net/http"

	"github.com/dvnthn168/ChatRealtime/pkg/websocket"
)

func main() {
	// Set up HTTP server and WebSocket route
	http.HandleFunc("/ws", websocket.HandleConnections)

	// Start the server
	log.Println("Server started on :8888")
	if err := http.ListenAndServe(":9999", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
