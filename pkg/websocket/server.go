package websocket

import (
	"log"
	"net/http"
)

func StartServer() {
	http.HandleFunc("/ws", HandleConnections)
	log.Println("WebSocket server is running at ws://localhost:8080/ws")
	http.ListenAndServe(":8080", nil)
}
