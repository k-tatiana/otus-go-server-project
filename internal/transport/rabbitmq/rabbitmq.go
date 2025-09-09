package rabbitmq

import (
	"fmt"
	"log"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQConfig struct {
	URL               string
	Host              string
	Port              int
	Username          string
	Password          string
	VHost             string
	ConnectionTimeout time.Duration
	Heartbeat         time.Duration
}

// Client клиент RabbitMQ
type Client struct {
	conn          *amqp.Connection
	channel       *amqp.Channel
	config        RabbitMQConfig
	isConnected   bool
	mu            sync.RWMutex
	reconnectChan chan bool
	done          chan bool
}

// NewClient создает новый клиент RabbitMQ
func NewRabbitMQClient(config RabbitMQConfig) *Client {
	if config.URL == "" {
		config.URL = fmt.Sprintf("amqp://%s:%s@%s:%d/%s",
			config.Username, config.Password, config.Host, config.Port, config.VHost)
	}

	client := &Client{
		config:        config,
		reconnectChan: make(chan bool),
		done:          make(chan bool),
	}

	go client.connectionManager()

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

	log.Println("Successfully connected to RabbitMQ")
	return nil
}

// Close закрывает соединение
func (c *Client) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	close(c.done)

	if c.channel != nil {
		if err := c.channel.Close(); err != nil {
			log.Printf("Error closing channel: %v", err)
		}
	}

	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			log.Printf("Error closing connection: %v", err)
		}
	}

	c.isConnected = false
	return nil
}

// IsConnected проверяет состояние соединения
func (c *Client) IsConnected() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.isConnected
}

// DeclareQueue объявляет очередь
func (c *Client) DeclareQueue(name string) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.isConnected {
		return fmt.Errorf("not connected to RabbitMQ")
	}

	_, err := c.channel.QueueDeclare(
		name,
		true, false, false, false, nil,
	)

	return err
}

// DeclareExchange объявляет exchange
func (c *Client) DeclareExchange(name, eType string) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.isConnected {
		return fmt.Errorf("not connected to RabbitMQ")
	}

	return c.channel.ExchangeDeclare(
		name, eType, true, false, false, false, nil,
	)
}

// BindQueue привязывает очередь к exchange
func (c *Client) BindQueue(queue, exchange, routingKey string) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.isConnected {
		return fmt.Errorf("not connected to RabbitMQ")
	}

	return c.channel.QueueBind(
		queue,
		routingKey,
		exchange,
		false,
		nil,
	)
}

// connectionManager управляет переподключениями
func (c *Client) connectionManager() {
	var reconnectDelay time.Duration = 1 * time.Second

	for {
		select {
		case <-c.done:
			return
		case <-c.reconnectChan:
			time.Sleep(reconnectDelay)

			if err := c.Connect(); err != nil {
				log.Printf("Reconnection failed: %v", err)
				reconnectDelay = min(reconnectDelay*2, 30*time.Second)
				c.reconnectChan <- true
			} else {
				reconnectDelay = 1 * time.Second
			}
		}
	}
}

// notifyReconnect уведомляет о необходимости переподключения
func (c *Client) notifyReconnect() {
	select {
	case c.reconnectChan <- true:
	default:
	}
}

// GetChannel возвращает AMQP channel (для advanced использования)
func (c *Client) GetChannel() (*amqp.Channel, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.isConnected || c.channel == nil {
		return nil, fmt.Errorf("not connected to RabbitMQ")
	}

	return c.channel, nil
}

func (c *Client) Consume(queue, consumer string) (<-chan amqp.Delivery, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.isConnected {
		return nil, fmt.Errorf("not connected to RabbitMQ")
	}

	msgs, err := c.channel.Consume(
		queue,
		consumer,
		false, // autoAck
		false, // exclusive
		false, // noLocal
		false, // noWait
		nil,   // args
	)
	if err != nil {
		return nil, fmt.Errorf("failed to start consuming: %w", err)
	}

	go func() {
		errChan := make(chan *amqp.Error)
		c.conn.NotifyClose(errChan)
		err := <-errChan
		if err != nil {
			log.Printf("Connection closed: %v", err)
			c.mu.Lock()
			c.isConnected = false
			c.mu.Unlock()
			c.notifyReconnect()
		}
	}()

	return msgs, nil
}

func (c *Client) Publish(exchange, routingKey string, body []byte) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// if !c.isConnected {
	// 	return fmt.Errorf("not connected to RabbitMQ")
	// }

	message := amqp.Publishing{
		ContentType: "application/octet-stream",
		Body:        body,
		Timestamp:   time.Now(),
	}

	err := c.channel.Publish(
		exchange,
		routingKey,
		false, // mandatory
		false, // immediate
		message,
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	return nil
}
