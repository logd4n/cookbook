package server

import (
	"fmt"
	"time"
	"webtest/pkg/colors"
)

func GetStartInfo(dbVersion *string) {
	//Информация при старте
	colors.SetColor(colors.Text_Yellow) //Установка желтого цвета
	fmt.Printf("Start info:\n")
	//Информация о сервере
	fmt.Printf("Server version - %s\n", appVersion)
	fmt.Printf("IP-address of server: %s\n", serverAddr)
	//Вывод времени
	fmt.Printf("Time:%s\t\n", time.Now().Format("02 January 2006, 15:04:05, MST"))
	//Вывод корневого каталога
	fmt.Printf("Root direction:%s\t\n", rootDir)
	//Вывод пути к файлу с логами
	fmt.Printf("Path to logs: %s\t ", logfilePath)
	colors.SetColor(colors.Text_Red)
	fmt.Printf("DOESN'T WORK YET :(\n") //ПОКА ЧТО НЕ РАБОТАЕТ
	colors.SetColor(colors.Text_Yellow)
	//Вывод версии БД
	fmt.Printf("Database version: %s\n", *dbVersion)
	colors.ResetColor() //Сброс желтого цвета

	//Вывод сообщения о запуске сервера
	colors.SetColor(colors.Text_Blue) //Установка зеленого цвета
	fmt.Printf("\nСервер запущен...\n\n")
	colors.ResetColor() //Сброс зеленого цвета
}
