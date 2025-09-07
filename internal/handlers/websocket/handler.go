package websocket

import (
	"context"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"otus/go-server-project/internal/models"
	"otus/go-server-project/internal/service"
)

type authenticator interface {
	Auth(ctx context.Context, userID string) error
}

type WebSocketHandler struct {
	Hub           *service.Hub
	authenticator authenticator
}

func NewWebSocketHandler(auth authenticator) *WebSocketHandler {
	return &WebSocketHandler{
		Hub: &service.Hub{
			Clients:    make(map[*service.WSClient]bool),
			Register:   make(chan *service.WSClient),
			Unregister: make(chan *service.WSClient),
			Broadcast:  make(chan models.WebSocketMessage, 100),
		},
		authenticator: auth}
}

func (h *WebSocketHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	currentUser := r.Header.Get("user_id")
	if currentUser == "" {
		http.Error(w, "user_id header required", http.StatusExpectationFailed)
		return
	}

	var upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // В production используйте proper CORS
		},
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	client := &service.WSClient{
		ID:         uuid.New().String(),
		Conn:       conn,
		Send:       make(chan []byte, 256),
		Subscribed: make(map[string]bool),
		UserID:     currentUser,
	}

	h.Hub.Register <- client

	go h.writePump(client)
	h.readPump(client)

	acceptKey := h.calculateAcceptKey(r.Header.Get("Sec-WebSocket-Key"))
	w.Header().Add("Sec-WebSocket-Accept", acceptKey)
}

// calculateAcceptKey generates the Sec-WebSocket-Accept header value as per RFC 6455.
func (h *WebSocketHandler) calculateAcceptKey(secWebSocketKey string) string {
	const magicGUID = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"
	hash := sha1Sum(secWebSocketKey + magicGUID)
	return base64.StdEncoding.EncodeToString(hash)
}

func sha1Sum(data string) []byte {
	h := sha1.New()
	h.Write([]byte(data))
	return h.Sum(nil)
}

func (h *WebSocketHandler) readPump(client *service.WSClient) {
	defer func() {
		h.Hub.Unregister <- client
		client.Conn.Close()
	}()

	client.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	client.Conn.SetPongHandler(func(string) error {
		client.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, message, err := client.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket read error: %v", err)
			}
			break
		}

		h.handleMessage(client, message)
	}
}

func (h *WebSocketHandler) handleMessage(client *service.WSClient, message []byte) {
	var wsMessage models.WebSocketMessage
	if err := json.Unmarshal(message, &wsMessage); err != nil {
		h.sendError(client, "Invalid message format", "")
		return
	}

	switch wsMessage.Type {
	case models.MessageTypeSubscribe:
		h.handleSubscribe(client, wsMessage)
	case models.MessageTypeUnsubscribe:
		h.handleUnsubscribe(client, wsMessage)
	case models.MessageTypePublish:
		h.handlePublish(client, wsMessage)
	default:
		h.sendError(client, "Unknown message type", wsMessage.PostID)
	}
}

func (h *WebSocketHandler) handleSubscribe(client *service.WSClient, msg models.WebSocketMessage) {
	if *msg.AuthorUserID == "" {
		h.sendError(client, "User to subscribe is required", msg.PostID)
		return
	}

	h.Hub.SubscribeClientToChannel(client, *msg.AuthorUserID)
	h.sendAck(client, "Subscribed successfully", msg.PostID)
	log.Printf("Client %s subscribed to %s", client.ID, *msg.AuthorUserID)
}

func (h *WebSocketHandler) handleUnsubscribe(client *service.WSClient, msg models.WebSocketMessage) {
	h.Hub.UnsubscribeClientFromChannel(client, *msg.AuthorUserID)
	h.sendAck(client, "Unsubscribed successfully", msg.PostID)
}

func (h *WebSocketHandler) handlePublish(client *service.WSClient, msg models.WebSocketMessage) {
	// Здесь может быть логика валидации и обработки публикаций
	h.sendAck(client, "Message processed", msg.PostID)
}

func (h *WebSocketHandler) writePump(client *service.WSClient) {
	ticker := time.NewTicker(30 * time.Second)
	defer func() {
		ticker.Stop()
		client.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-client.Send:
			client.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := client.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			client.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := client.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (h *WebSocketHandler) sendAck(client *service.WSClient, message string, messageID string) {
	response := models.ResponseMessage{
		Type:      models.MessageTypeAck,
		MessageID: messageID,
		Data:      message,
		Success:   true,
	}

	jsonData, _ := json.Marshal(response)
	client.Send <- jsonData
}

func (h *WebSocketHandler) sendError(client *service.WSClient, errorMsg string, messageID string) {
	response := models.ResponseMessage{
		Type:      models.MessageTypeError,
		MessageID: messageID,
		Error:     errorMsg,
		Success:   false,
	}

	jsonData, _ := json.Marshal(response)
	client.Send <- jsonData
}
