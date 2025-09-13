package service

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/rabbitmq/amqp091-go"

	"otus/go-server-project/internal/models"
)

type WSClient struct {
	ID         string
	Conn       *websocket.Conn
	Send       chan []byte
	Subscribed map[string]bool
	UserID     string
}

type rmq interface {
	Consume(queue string) (<-chan amqp091.Delivery, error)
	Publish(exchange, routingKey string, body []byte) error
}

type Hub struct {
	Clients    map[*WSClient]bool
	Register   chan *WSClient
	Unregister chan *WSClient
	Broadcast  chan models.WebSocketMessage

	PostsQueue string

	RMQ rmq

	mu sync.RWMutex
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

		case message := <-h.BroadcastMessage(h.PostsQueue):
			h.mu.RLock()

			bytes, err := json.Marshal(message)
			if err != nil {
				fmt.Printf("unable to marshal message %d", message)
			}

			for client := range h.Clients {
				if client.Subscribed[message.RoutingKey] {
					select {
					case client.Send <- bytes:
						fmt.Print("sent to client")
						message.Ack(false)
					default:
						close(client.Send)
						delete(h.Clients, client)
						message.Nack(false, true)
					}
				}
			}
			h.mu.RUnlock()
		}
	}
}

func (h *Hub) BroadcastMessage(queue string) <-chan amqp091.Delivery {
	chanMsg, err := h.RMQ.Consume(queue)
	if err != nil {
		fmt.Println(fmt.Errorf("unable to consume message: %v", err))
	}

	return chanMsg
}

func (h *Hub) PublishToChannel(exchange, routingKey string, data models.WebSocketMessage) error {
	msg, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	err = h.RMQ.Publish(exchange, routingKey, msg)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

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
