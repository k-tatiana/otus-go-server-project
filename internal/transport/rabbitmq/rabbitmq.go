package rabbitmq

import (
	"fmt"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type ConnectionConfig struct {
	URL               string
	Host              string
	Port              int
	Username          string
	Password          string
	VHost             string
	ConnectionTimeout time.Duration
	Heartbeat         time.Duration
}

// // DefaultConfig возвращает конфигурацию по умолчанию
// func DefaultConfig() ConnectionConfig {
// 	return ConnectionConfig{
// 		Host:              "localhost",
// 		Port:              5672,
// 		Username:          "guest",
// 		Password:          "guest",
// 		VHost:             "/",
// 		ConnectionTimeout: 30 * time.Second,
// 		Heartbeat:         10 * time.Second,
// 	}
// }

// Client клиент RabbitMQ
type Client struct {
	conn          *amqp.Connection
	channel       *amqp.Channel
	config        ConnectionConfig
	isConnected   bool
	mu            sync.RWMutex
	reconnectChan chan bool
	done          chan bool
}

// NewClient создает новый клиент RabbitMQ
func NewClient(config ConnectionConfig) *Client {
	if config.URL == "" {
		config.URL = fmt.Sprintf("amqp://%s:%s@%s:%d/%s",
			config.Username, config.Password, config.Host, config.Port, config.VHost)
	}

	client := &Client{
		config:        config,
		reconnectChan: make(chan bool),
		done:          make(chan bool),
	}

	if err := client.Connect(); err != nil {
		fmt.Printf("Failed to connect to RabbitMQ: %v\n", err)
	}

	return client
}

// Connect устанавливает соединение с RabbitMQ
func (c *Client) Connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.isConnected {
		return nil
	}

	conn, err := amqp.DialConfig(c.config.URL, amqp.Config{
		Dial:      amqp.DefaultDial(c.config.ConnectionTimeout),
		Heartbeat: c.config.Heartbeat,
	})
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return fmt.Errorf("failed to open channel: %w", err)
	}

	c.conn = conn
	c.channel = channel
	c.isConnected = true

	fmt.Println("Successfully connected to RabbitMQ")
	return nil
}
