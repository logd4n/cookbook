package logger

import (
	"context"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func NewMessage() error {
	// 1.
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		log.Printf("Не удалось подключиться к RabbitMQ (адрес:\"amqp://guest:guest@rabbitmq:5672/\") ")
		log.Printf("Error: %v", err.Error())
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Printf("Не удалось создать канал! ")
		log.Printf("Error: %v\n", err.Error())
		return err
	}
	defer ch.Close()

	//2.
	err = ch.ExchangeDeclare(
		"logs_exchange",
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Printf("Не удалось объявить exchange! ")
		log.Printf("Error: %v\n", err.Error())
		return err
	}

	//3.
	queue, err := ch.QueueDeclare(
		"test_queue",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Printf("Не удалось объявить queue! ")
		log.Printf("Error: %v\n", err.Error())
		return err
	}

	//4.
	err = ch.QueueBind(
		queue.Name,
		"test_key",
		"logs_exchange",
		false,
		nil,
	)
	if err != nil {
		log.Printf("Не удалось связать exchange и queue! ")
		log.Printf("Error: %v\n", err.Error())
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	//5.
	body := "Hello RabbitMQ!"
	err = ch.PublishWithContext(ctx, "logs_exchange",
		"test_key",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		},
	)
	if err != nil {
		log.Printf("Не удалось опубликовать сообщение! ")
		log.Printf("Error: %v\n", err.Error())
		return err
	}

	return nil
}

func ReadMessage() error {
	//1.
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		log.Printf("Не удалось подключиться к RabbitMQ (адрес:\"amqp://guest:guest@rabbitmq:5672/\") ")
		log.Printf("Error: %v", err.Error())
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Printf("Не удалось создать канал! ")
		log.Printf("Error: %v\n", err.Error())
		return err
	}
	defer ch.Close()

	//2.
	err = ch.ExchangeDeclare(
		"logs_exchange",
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Printf("Не удалось объявить exchange! ")
		log.Printf("Error: %v\n", err.Error())
		return err
	}

	//3.
	queue, err := ch.QueueDeclare(
		"test_queue",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Printf("Не удалось объявить queue! ")
		log.Printf("Error: %v\n", err.Error())
		return err
	}

	//4.
	err = ch.QueueBind(
		queue.Name,
		"test_key",
		"logs_exchange",
		false,
		nil,
	)
	if err != nil {
		log.Printf("Не удалось связать exchange и queue! ")
		log.Printf("Error: %v\n", err.Error())
		return err
	}

	//5.
	msgs, err := ch.Consume(
		queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Printf("Не удалось прочитать сообщение из очереди! ")
		log.Printf("Error: %v\n", err.Error())
		return err
	}

	go func() {
		for m := range msgs {
			log.Printf("MESSAGE: %s\n", m.Body)
			m.Ack(false)
		}
	}()

	time.Sleep(1 * time.Second)
	return nil
}
