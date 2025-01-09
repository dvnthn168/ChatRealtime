package main

import (
	"log"
	"net/http"

	"github.com/dvnthn168/ChatRealtime/pkg/redisdb"
	"github.com/dvnthn168/ChatRealtime/pkg/websocket"
)

func main() {

	rdb, err := redisdb.HandleConnections()
	if err != nil {
		log.Fatalf("Could not connect to Redis %v", err)
	}
	defer rdb.Close()

	handler := &websocket.WebSocketHandler{RedisClient: rdb}

	http.HandleFunc("/ws", handler.HandleConnections)

	log.Println("Server started on :9999")
	if err := http.ListenAndServe(":9999", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
