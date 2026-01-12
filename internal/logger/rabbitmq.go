package logger

import (
	"time"
	rabbitconfig "webtest/internal/rabbit-config"

	"github.com/rabbitmq/amqp091-go"
)

type RabbitMQClient struct {
	conn    *amqp091.Connection
	channel *amqp091.Channel
	config  rabbitconfig.RabbitMQConfig
}

func NewRabbitMQClient(config rabbitconfig.RabbitMQConfig) (*RabbitMQClient, error) {
	var conn *amqp091.Connection
	var err error

	conn, err = amqp091.Dial(config.URL)
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	err = channel.ExchangeDeclare(
		config.Exchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	err = channel.QueueBind(
		config.Queue,
		"#",
		config.Exchange,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &RabbitMQClient{
		conn:    conn,
		channel: channel,
		config:  config,
	}, nil
}

func (r *RabbitMQClient) PublishLog(message string) error {
	routingKey := "logs.test"

	return r.channel.Publish(
		r.config.Exchange,
		routingKey,
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        []byte(message),
			Timestamp:   time.Now(),
		},
	)
}

func (r *RabbitMQClient) Close() {
	r.channel.Close()
	r.conn.Close()
}
