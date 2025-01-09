package websocket

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin:     func(r *http.Request) bool { return true },
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var clients = make(map[*websocket.Conn]bool)

func HandleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial HTTP request to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	// Register the new client
	clients[conn] = true
	log.Println("New client connected:", conn.RemoteAddr())

	// Listen for messages from the client
	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			delete(clients, conn)
			break
		}

		// Log the incoming message for debugging purposes
		log.Printf("Received message: %s", string(msg))

		// Broadcast the message to all connected clients
		for client := range clients {
			if err := client.WriteMessage(msgType, msg); err != nil {
				log.Println("Write error:", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
