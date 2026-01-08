package main

//version 2.2

import (
	"webtest/internal/database"
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

	//Запускаем сервер
	server.StartServer(&dbVersion)
}
