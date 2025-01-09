package websocket

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/dvnthn168/ChatRealtime/pkg/models"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

var upgrader = websocket.Upgrader{
	CheckOrigin:     func(r *http.Request) bool { return true },
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type WebSocketHandler struct {
	RedisClient *redis.Client
}

var (
	ctx       = context.Background()
	clientMap sync.Map
)

func (h *WebSocketHandler) HandleConnections(w http.ResponseWriter, r *http.Request) {

	if h.RedisClient == nil {
		log.Println("Redis client is not available")
		return
	}
	if err := h.RedisClient.Ping(ctx).Err(); err != nil {
		log.Printf("Redis connection error: %v", err)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	userID := r.URL.Query().Get("userID")
	if userID == "" {
		log.Println("userID is missing")
		return
	}

	clientMap.Store(userID, conn)
	defer clientMap.Delete(userID)

	log.Printf("Client connected: %s\n", userID)

	pubsub := h.RedisClient.Subscribe(ctx, userID) 
	defer pubsub.Close()

	go func() {
		for msg := range pubsub.Channel() {
			log.Printf("Received message: %s", msg.Payload)
			if err := conn.WriteMessage(websocket.TextMessage, []byte(msg.Payload)); err != nil {
				log.Println("Write error:", err)
				return
			}
		}
	}()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}

		log.Printf("Received message: %s", string(msg))
		var m models.Message
		if err := json.Unmarshal(msg, &m); err != nil {
			log.Println("Unmarshal error:", err)
			continue
		}

		if h.RedisClient != nil {
			log.Printf("Recipient missing in the message %v", m.Recipient)
			if err := h.RedisClient.Publish(ctx, m.Recipient, m.Message).Err(); err != nil {
				log.Printf("Publish error: %v", err)
				panic(err)
			}
		} else {
			log.Println("Redis client is not available")
		}
	}
}
