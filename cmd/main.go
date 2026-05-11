package main

//version 2.4.1

import (
	"webtest/internal/database"
	"webtest/internal/logger"
	"webtest/internal/models"
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
	err = logger.NewMessage(models.LogMessage{
		Level:   models.Info,
		Message: "Сервер запущен!",
	})
	if err != nil {
		return
	}

	//Запускаем сервер
	server.StartServer(&dbVersion)
}
