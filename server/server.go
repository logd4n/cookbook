package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
	"webtest/colors"
	. "webtest/readBody"
)

var rootDir string // ./webtest

func initServer() {
	workDir, err := os.Getwd() // ./webtest/cmd
	if err != nil {
		log.Fatal(err)
	}

	rootDir = filepath.Dir(workDir) // ./webtest
}

func StartServer() {
	initServer()

	//Информация
	colors.SetColor(colors.Text_Yellow) //Установка желтого цвета
	fmt.Println("Start info:")
	//Вывод времени
	fmt.Println("Time:\t", time.Now().Format("02 January 2006, 15:04:05, MST"))
	//Вывод корневого каталога
	fmt.Println("Root direction:\t", rootDir)
	colors.ResetColor() //Сброс желтого цвета

	//Вывод сообщения о запуске сервера
	colors.SetColor(colors.Text_Green) //Установка зеленого цвета
	fmt.Print("\nСервер запущен...\n\n")
	colors.ResetColor() //Сброс зеленого цвета

	//Обработчик для CSS
	http.Handle("/static/css/",
		http.StripPrefix("/static/css/",
			http.FileServer(
				http.Dir(filepath.Join(rootDir, "templates", "static", "css")))))

	//Обработчик для JS
	http.Handle("/static/js/",
		http.StripPrefix("/static/js/",
			http.FileServer(
				http.Dir(filepath.Join(rootDir, "templates", "static", "js")))))

	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/add", addHandler)

	http.ListenAndServe("0.0.0.0:8080", nil)
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, rootDir+"/templates/index.html")
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	// Разрешаем CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Метод "+r.Method+" не поддерживается!", http.StatusMethodNotAllowed)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Ожидались данные формата JSON!", http.StatusBadRequest)
	}

	err := ReadBody(r, rootDir)
	if err != nil {
		colors.SetColor(colors.Text_Red)
		log.Print(err, "\n\n")
		colors.ResetColor()

		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
