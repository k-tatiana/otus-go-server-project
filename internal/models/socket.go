package models

type MessageType string

const (
	MessageTypeSubscribe   MessageType = "subscribe"
	MessageTypeUnsubscribe MessageType = "unsubscribe"
	MessageTypePublish     MessageType = "publish"
	MessageTypeAck         MessageType = "ack"
	MessageTypeError       MessageType = "error"
)

type ResponseMessage struct {
	Type      MessageType `json:"type"`
	MessageID string      `json:"message_id,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	Error     string      `json:"error,omitempty"`
	Success   bool        `json:"success"`
}
