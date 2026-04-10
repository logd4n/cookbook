package main

//version 2.4.1

import (
	"log"
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
	dataBase, dbVersion, err := database.ConnectDB()
	if err != nil {
		return
	}
	defer dataBase.Close()

	//Пытаемся подключиться к RabbitMQ
	err = logger.ConnectionAttempt()
	if err != nil {
		return
	}
	err = logger.NewMessage()
	if err != nil {
		return
	}
	log.Printf("---Сообщение отправлено---\n")

	//Запускаем сервер
	log.Printf("Сервер запущен!")
	server.StartServer(&dbVersion)
}
