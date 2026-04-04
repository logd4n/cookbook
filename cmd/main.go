package main

//version 2.4.1

import (
	"log"
	"time"
	"webtest/internal/database"
	"webtest/internal/logger"
	"webtest/internal/server"
	"webtest/pkg/colors"
	"webtest/pkg/intro"
)

func init() {
	colors.SetColor(colors.Text_Red)
	intro.IntroLog()
	colors.ResetColor()
}

func main() {
	//Подключение к БД и вывод ее версии
	dataBase, dbVersion := database.ConnectDB()
	defer dataBase.Close()

	///////////////////////// <--- RabbitMQ TEST BEGIN
	time.Sleep(3 * time.Second) //Ждем RabbitMQ
	err := logger.NewMessage()
	if err != nil {
		return
	}
	log.Printf("---Сообщение отправлено---\n")

	/*
		time.Sleep(500 * time.Millisecond)
		err = logger.ReadMessage()
		if err != nil {
			return
		}
		log.Printf("---Сообщение прочтено---\n\n")
		///////////////////////// <--- RabbitMQ TEST END
	*/

	//Запускаем сервер
	log.Printf("Сервер запущен!")
	server.StartServer(&dbVersion)
}
