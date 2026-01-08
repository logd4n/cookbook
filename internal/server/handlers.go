package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"webtest/internal/database"
	"webtest/internal/models"
	. "webtest/internal/writeData"
	"webtest/pkg/colors"
)

// Обработчик для localhost:8080/
func mainHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, rootDir+"/templates/index_2026.html")

	//Получение IP-адреса клиента и вывод подлючения в консоль
	client, err := getClientIP(r)
	if err != nil {
		colors.SetColor(colors.Text_Red)
		log.Println("Не удалось получить IP-адрес клиента...")
		colors.ResetColor()
	}
	colors.SetColor(colors.Text_Purple)
	log.Printf("Выполнено подключение к \"%s\"! Client ip: [%s]\n\n", serverAddr+"/", client)
	colors.ResetColor()
}

// Обработчик для localhost:8080/add
func addHandler(w http.ResponseWriter, r *http.Request) {
	// Разрешаем CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	//Проверка метода OPTIONS
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	//Разрешаем только метод POST
	if r.Method != http.MethodPost {
		http.Error(w, "Метод "+r.Method+" не поддерживается!", http.StatusMethodNotAllowed)
		return
	}

	//Проверка заголовка запроса на содержание JSON
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Ожидались данные формата JSON!", http.StatusBadRequest)
	}

	//Десериализация тела запроса
	eat_data, err := Deserialization(r)
	if err != nil {
		colors.SetColor(colors.Text_Red)
		log.Print(err, "\n\n")
		colors.ResetColor()

		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//Запись данных в БД
	err = database.WriteDB((*models.Eat)(eat_data))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	client, _ := getClientIP(r) //Получение IP-адреса клиента
	colors.SetColor(colors.Text_Purple)
	log.Printf("Клиент [%s] записал данные на сервер!\n\n", client)
	colors.ResetColor()

	//Возвращаем статус OK
	w.WriteHeader(http.StatusOK)
}

// Обработчик для localhost:8080/search
func searchHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, rootDir+"/templates/search_2026.html")

	client, _ := getClientIP(r) //Получаем ip клиента
	colors.SetColor(colors.Text_Purple)
	log.Printf("Выполнено подключение к \"%s\"! Client ip: [%s]\n\n", serverAddr+"/search", client)
	colors.ResetColor()

	// Разрешаем CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
}

// Обработчик для localhost:8080/recipes
func getAllRecipesHandler(w http.ResponseWriter, r *http.Request) {
	// Разрешаем CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	//Разрешаем только метод GET
	if r.Method != http.MethodGet {
		http.Error(w, "Метод "+r.Method+" не поддерживается!", http.StatusMethodNotAllowed)
		return
	}

	//Отправляем запрос в БД
	data, err := database.GetAllRecipes()
	if err != nil {
		http.Error(w, "При выполнении запроса возникла ошибка!", http.StatusBadRequest)
		return
	}

	//Возвращаем структуру JSON
	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Обработчик для localhost:8080/recipes/
func getOneRecipeHandler(w http.ResponseWriter, r *http.Request) {
	// Разрешаем CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	//Разрешаем только метод GET
	if r.Method != http.MethodGet {
		http.Error(w, "Метод "+r.Method+" не поддерживается!", http.StatusMethodNotAllowed)
		return
	}

	// Сплитуем адрес по сепаратору /
	urlParts := strings.Split(r.URL.Path, "/") // [api recipes X] X --> id of recipe
	if len(urlParts) < 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	//Получаем id рецепта от результата сплита
	recipeID, err := strconv.Atoi(urlParts[3])
	if err != nil {
		http.Error(w, "Не удалось обработать URL!", http.StatusBadRequest)
		return
	}

	//Отправляем запрос в БД
	recipe, err := database.GetOneRecipe(recipeID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	//Возвращаем структуру JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recipe)
}

// Обработчик для localhost:8080/deleteRecipe/
func deleteRecipeHandler(w http.ResponseWriter, r *http.Request) {
	// Разрешаем CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	//Разрешаем только DELETE
	if r.Method != http.MethodDelete {
		http.Error(w, "Метод "+r.Method+" не поддерживается!", http.StatusMethodNotAllowed)
		return
	}

	//Сплит адреса по сепаратору "/"
	urlParts := strings.Split(r.URL.Path, "/") //[api deleteRecipe X] X --> id of recipe
	if len(urlParts) < 3 {
		http.Error(w, "Invalid URL!", http.StatusBadRequest)
		return
	}

	//Получаем ID рецепта
	recipeID, err := strconv.Atoi(urlParts[3])
	if err != nil {
		http.Error(w, "Не удалось обработать URL!", http.StatusBadRequest)
		return
	}

	//Отправляем запрос на удаление в БД
	err = database.DeleteRecipe(recipeID)
	if err != nil {
		http.Error(w, "Не удалось выполнить удаление!", http.StatusBadGateway)
		return
	}

	client, _ := getClientIP(r) //Получение ip клиента
	colors.SetColor(colors.Text_Purple)
	log.Printf("Клиент [%s] удалил данные по id [%d] на сервере!\n\n", client, recipeID)
	colors.ResetColor()

	//Возвращаем статус 200 ОК и сообщение
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Удаление выполнено успешно!"))
}

// Обработчик для localhost:8080/updateRecipe/
func updateRecipeHandler(w http.ResponseWriter, r *http.Request) {
	// Разрешаем CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	//Разрешаем только PUT
	if r.Method != http.MethodPut {
		http.Error(w, "Метод "+r.Method+" не поддерживается!", http.StatusMethodNotAllowed)
		return
	}

	//Десериализуем тело запроса
	data, err := Deserialization(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	//Отправляем запрос в БД
	err = database.UpdateRecipe(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	client, _ := getClientIP(r) //Получаем ip клиента
	colors.SetColor(colors.Text_Purple)
	log.Printf("Клиент [%s] обновил данные по id [%d] на сервере!\n\n", client, data.ID)
	colors.ResetColor()

	//Возвращаем статус 200 ОК и сообщение
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Редактирование выполнено успешно!"))
}
