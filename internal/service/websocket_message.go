package service

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/gorilla/websocket"

	"otus/go-server-project/internal/models"
)

type WSClient struct {
	ID         string
	Conn       *websocket.Conn
	Send       chan []byte
	Subscribed map[string]bool
	UserID     string
}

type Hub struct {
	Clients    map[*WSClient]bool
	Register   chan *WSClient
	Unregister chan *WSClient
	Broadcast  chan models.WebSocketMessage
	mu         sync.RWMutex
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			h.Clients[client] = true
			h.mu.Unlock()
			fmt.Printf("Client %s connected", client.ID)

		case client := <-h.Unregister:
			h.mu.Lock()
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
			}
			h.mu.Unlock()
			fmt.Printf("Client %s disconnected", client.ID)

		case message := <-h.Broadcast:
			h.mu.RLock()
			bytes, err := json.Marshal(message)
			if err != nil {
				fmt.Printf("unable to marshal message %d", message)
			}
			for client := range h.Clients {
				if client.Subscribed[*message.AuthorUserID] {
					select {
					case client.Send <- bytes:
						fmt.Print("sent to client")
					default:
						close(client.Send)
						delete(h.Clients, client)
					}
				}
			}
			h.mu.RUnlock()
		}
	}
}

func (h *Hub) PublishToChannel(channel string, data interface{}) error {
	message := data.(models.WebSocketMessage)

	h.Broadcast <- message
	return nil
}

func (h *Hub) SubscribeClientToChannel(client *WSClient, channel string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	client.Subscribed[channel] = true
}

func (h *Hub) UnsubscribeClientFromChannel(client *WSClient, channel string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(client.Subscribed, channel)
}
