package server

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var (
	appVersion  string
	serverAddr  string
	rootDir     string // ./webtest
	logfilePath string
)

// Функция инициализации сервера
func initServer() {
	appVersion = os.Getenv("APP_VERSION")
	serverAddr = os.Getenv("SERVER_ADDR")

	workDir, err := os.Getwd() //Получаем путь ./webtest/cmd
	if err != nil {
		log.Fatal(err)
	}

	rootDir = workDir //filepath.Dir(workDir) // Убираем каталог выполнения exe-файла ./webtest
	logfilePath = filepath.Join(rootDir, "logs", "server.log")
}

func StartServer(dbVersion *string) {
	initServer() //Инициализируем сервер

	GetStartInfo(dbVersion)

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

	//Обработчик для favicon.ico
	http.Handle("/static/",
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("static"))))

	http.HandleFunc("/static/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/x-icon")
		http.ServeFile(w, r, filepath.Join(rootDir, "templates", "static", "favicon.ico"))
	})

	//Обработчики эндпоинтов
	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/add", addHandler)
	http.HandleFunc("/search", searchHandler)
	http.HandleFunc("/api/recipes", getAllRecipesHandler)
	http.HandleFunc("/api/recipes/", getOneRecipeHandler)
	http.HandleFunc("/api/deleteRecipe/", deleteRecipeHandler)
	http.HandleFunc("/api/updateRecipe/", updateRecipeHandler)

	http.ListenAndServe("0.0.0.0:8080", nil)
}
