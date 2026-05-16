package logger

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"
	"webtest/internal/models"

	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	ConnErr    = errors.New("Ошибка подключения")
	ConnTmtErr = errors.New("Превышено время ожидания RabbitMQ!")

	DeclareChErr    = errors.New("Не удалось создать канал!")
	DeclareExErr    = errors.New("Не удалось создать Exchange!")
	DeclareQueueErr = errors.New("Не удалось создать очередь!")
	QueueBindErr    = errors.New("Не удалось связать exchange и queue!")
	MarshallErr     = errors.New("Ошибка перевода структуры в byte[]: ")
	PublishErr      = errors.New("Не удалось опубликовать сообщение!")
)

func ConnectionAttempt() error {
	var conn *amqp.Connection
	var err error

	for i := 1; i <= 10; i++ {
		log.Printf("Попытка подключения к RabbitMQ №%d...\t\n", i)

		conn, err = amqp.Dial("amqp://guest:guest@rabbitmq:5672/")

		if err == nil {
			conn.Close()
			break
		}

		log.Printf("%v: \"%v\"\n", ConnErr.Error(), err.Error())
		time.Sleep(5 * time.Second)
	}

	if err != nil {
		log.Printf("%v Ошибка: %v\n", ConnTmtErr, err.Error())
		return err
	}

	log.Printf("Подключение к RabbitMQ выполнено успешно!")
	return nil
}

func NewMessage(message models.LogMessage) error {
	// 1.
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		log.Printf("%v (адрес:\"amqp://guest:guest@rabbitmq:5672/\"): %v\n",
			ConnErr.Error(), err.Error())
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Printf("%v: %v\n", DeclareChErr.Error(), err.Error())
		return err
	}
	defer ch.Close()

	//2.
	err = ch.ExchangeDeclare(
		"logs_exchange",
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Printf("%v: %v\n", DeclareExErr.Error(), err.Error())
		return err
	}

	//3.
	queue, err := ch.QueueDeclare(
		"q.database.saver",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Printf("%v: %v\n", DeclareQueueErr.Error(), err.Error())
		return err
	}

	//4.
	err = ch.QueueBind(
		queue.Name,
		"logs.#",
		"logs_exchange",
		false,
		nil,
	)
	if err != nil {
		log.Printf("%v: %v\n", QueueBindErr.Error(), err.Error())
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	//5.
	body, err := json.Marshal(message)
	if err != nil {
		return MarshallErr
	}

	err = ch.PublishWithContext(ctx, "logs_exchange",
		"logs."+string(message.Level),
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(body),
			Timestamp:   time.Now().UTC(),
		},
	)
	if err != nil {
		log.Printf("%v: %v\n", PublishErr.Error(), err.Error())
		return err
	}

	return nil
}
