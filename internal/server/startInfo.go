package server

import (
	"fmt"
	"time"
	"webtest/pkg/colors"
)

func GetStartInfo(dbVersion *string) {
	//Информация при старте
	colors.SetColor(colors.Text_Yellow) //Установка желтого цвета
	fmt.Printf("\nStart info:\n")
	//Информация о сервере
	fmt.Printf("Server version:\t%s\n", appVersion)
	fmt.Printf("IP-address of server:\t%s\n", serverAddr)
	//Вывод времени
	fmt.Printf("Time:\t%s\n", time.Now().Format("02 January 2006, 15:04:05, MST"))
	//Вывод корневого каталога
	fmt.Printf("Root direction:\t%s\n", rootDir)
	//Вывод пути к файлу с логами
	fmt.Printf("Path to logs:\t%s", logfilePath)
	colors.SetColor(colors.Text_Red)
	fmt.Printf("DOESN'T WORK YET :(\n") //ПОКА ЧТО НЕ РАБОТАЕТ
	colors.SetColor(colors.Text_Yellow)
	//Вывод версии БД
	fmt.Printf("Database version:\t%s\n", *dbVersion)
	colors.ResetColor() //Сброс желтого цвета

	//Вывод сообщения о запуске сервера
	colors.SetColor(colors.Text_Blue) //Установка зеленого цвета
	fmt.Printf("\nСервер запущен...\n\n")
	colors.ResetColor() //Сброс зеленого цвета
}
