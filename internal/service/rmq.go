package service

import (
	"fmt"
	"log"
	"otus/go-server-project/internal/transport/rabbitmq"

	"github.com/rabbitmq/amqp091-go"
)

const (
	DefaultQueueName    = "posts_queue"
	DefaultExchangeName = "posts_exchange"
	DefaultExchangeType = "fanout"
	ConsumerName        = ""
)

type Rmq struct {
	client *rabbitmq.Client
}

func NewRmq(cfg rabbitmq.RabbitMQConfig) *Rmq {
	rmqClient := *rabbitmq.NewRabbitMQClient(cfg)
	if err := rmqClient.Connect(); err != nil {
		log.Fatalf("Could not connect to RabbitMQ: %v", err)
	}

	if err := rmqClient.DeclareQueue(DefaultQueueName); err != nil {
		log.Fatalf("Could not declare queue: %v", err)
	}

	if err := rmqClient.DeclareExchange(DefaultExchangeName, DefaultExchangeType); err != nil {
		log.Fatalf("Could not setup exchange: %v", err)
	}
	if err := rmqClient.BindQueue(DefaultQueueName, DefaultExchangeName, ""); err != nil {
		log.Fatalf("Could not setup and bind queue: %v", err)
	}
	fmt.Println("RMQ Client connected: ", rmqClient.IsConnected())
	return &Rmq{client: &rmqClient}
}

func (r *Rmq) Consume(queue string) (<-chan amqp091.Delivery, error) {
	return r.client.Consume(queue, ConsumerName)
}

func (r *Rmq) Publish(exchange, routingKey string, body []byte) error {
	return r.client.Publish(exchange, routingKey, body)
}
