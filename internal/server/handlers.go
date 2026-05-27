package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"webtest/internal/database"
	"webtest/internal/logger"
	"webtest/internal/models"
	. "webtest/internal/writeData"
	"webtest/pkg/colors"
)

var (
	InvMethodErr = errors.New("Метод не поддерживается:")
	ClientIPErr  = errors.New("Не удалось получить IP-адрес клиента...")
	JSON_Err     = errors.New("Ожидались данные формата JSON!")
	MarshallErr  = errors.New("Не удалось вернуть данные в формате JSON!")
	InvURL_Err   = errors.New("Invalid URL!")
	GetID_Err    = errors.New("Не удалось получить ID!")
)

// Обработчик для localhost:8080/
func mainHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, rootDir+"/templates/index.html")

	//Получение IP-адреса клиента и вывод подлючения в консоль
	client, err := getClientIP(r)
	if err != nil {
		colors.SetColor(colors.Text_Red)
		logger.LogPrint(ClientIPErr.Error(), models.Error)
		colors.ResetColor()
	}
	colors.SetColor(colors.Text_Purple)
	logger.LogPrint(fmt.Sprintf("Выполнено подключение к \"%s\"! Client ip: [%s]",
		serverAddr+"/",
		client),
		models.Info)
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
		http.Error(w,
			fmt.Sprintf("%v %v", InvMethodErr.Error(), r.Method),
			http.StatusMethodNotAllowed)

		logger.LogPrint(
			fmt.Sprintf("%v %v", InvMethodErr.Error(), r.Method),
			models.Error)
		return
	}

	//Проверка заголовка запроса на содержание JSON
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, JSON_Err.Error(), http.StatusBadRequest)

		logger.LogPrint(JSON_Err.Error(), models.Error)
	}

	//Десериализация тела запроса
	eat_data, err := Deserialization(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		colors.SetColor(colors.Text_Red)
		logger.LogPrint(err.Error(), models.Error)
		colors.ResetColor()
		return
	}

	//Запись данных в БД
	err = database.WriteDB((*models.Eat)(eat_data))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)

		logger.NewMessage(models.LogMessage{
			Level:   models.Error,
			Message: fmt.Sprintf("Не удалось записать данные в БД: %v", err.Error()),
		})
		return
	}

	client, err := getClientIP(r) //Получение IP-адреса клиента
	if err != nil {
		colors.SetColor(colors.Text_Red)
		logger.LogPrint(ClientIPErr.Error(), models.Error)
		colors.ResetColor()
	}

	colors.SetColor(colors.Text_Purple)
	logger.LogPrint(
		fmt.Sprintf("Клиент [%s] записал данные на сервер!\n\n", client),
		models.Info,
	)
	colors.ResetColor()

	//Возвращаем статус OK
	w.WriteHeader(http.StatusOK)
}

// Обработчик для localhost:8080/search
func searchHandler(w http.ResponseWriter, r *http.Request) {
	// Разрешаем CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	http.ServeFile(w, r, rootDir+"/templates/search.html")

	client, err := getClientIP(r) //Получаем ip клиента
	if err != nil {
		colors.SetColor(colors.Text_Red)
		logger.LogPrint(ClientIPErr.Error(), models.Error)
		colors.ResetColor()
	}

	colors.SetColor(colors.Text_Purple)
	logger.LogPrint(
		fmt.Sprintf("Выполнено подключение к \"%s\"! Client ip: [%s]\n\n", serverAddr+"/search", client),
		models.Info,
	)
	colors.ResetColor()
}

// Обработчик для localhost:8080/api/recipes
func getAllRecipesHandler(w http.ResponseWriter, r *http.Request) {
	// Разрешаем CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	//Разрешаем только метод GET
	if r.Method != http.MethodGet {
		colors.SetColor(colors.Text_Red)
		logger.LogPrint(
			fmt.Sprintf("%v %v", InvMethodErr.Error(), r.Method),
			models.Error,
		)
		http.Error(w,
			fmt.Sprintf("%v %v", InvMethodErr.Error(), r.Method),
			http.StatusMethodNotAllowed,
		)
		colors.ResetColor()
		return
	}

	//Отправляем запрос в БД
	data, err := database.GetAllRecipes()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		colors.SetColor(colors.Text_Red)
		logger.LogPrint(err.Error(), models.Error)
		colors.ResetColor()
		return
	}

	//Возвращаем структуру JSON
	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, MarshallErr.Error(), http.StatusInternalServerError)

		colors.SetColor(colors.Text_Red)
		logger.LogPrint(
			MarshallErr.Error(),
			models.Error,
		)
		colors.ResetColor()
	}
}

// Обработчик для localhost:8080/api/recipes/
func getOneRecipeHandler(w http.ResponseWriter, r *http.Request) {
	// Разрешаем CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	//Разрешаем только метод GET
	if r.Method != http.MethodGet {
		http.Error(w,
			fmt.Sprintf("%v %v", InvMethodErr, r.Method),
			http.StatusMethodNotAllowed,
		)

		colors.SetColor(colors.Text_Red)
		logger.LogPrint(
			fmt.Sprintf("%v %v", InvMethodErr, r.Method),
			models.Error,
		)
		colors.ResetColor()
		return
	}

	// Сплитуем адрес по сепаратору "/"
	urlParts := strings.Split(r.URL.Path, "/") // [api recipes X] X --> id of recipe
	if len(urlParts) < 3 {
		http.Error(w, InvURL_Err.Error(), http.StatusBadRequest)
		logger.LogPrint(
			InvURL_Err.Error(),
			models.Error,
		)
		return
	}

	//Получаем id рецепта от результата сплита
	recipeID, err := strconv.Atoi(urlParts[3])
	if err != nil {
		http.Error(w, GetID_Err.Error(), http.StatusBadRequest)
		logger.LogPrint(
			GetID_Err.Error(),
			models.Error,
		)
		return
	}

	//Отправляем запрос в БД
	recipe, err := database.GetOneRecipe(recipeID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		logger.LogPrint(
			err.Error(),
			models.Error,
		)
		return
	}

	//Возвращаем структуру JSON
	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(recipe); err != nil {
		http.Error(
			w,
			MarshallErr.Error(),
			http.StatusInternalServerError,
		)

		logger.LogPrint(
			MarshallErr.Error(),
			models.Error,
		)
	}
}

// Обработчик для localhost:8080/api/deleteRecipe/
func deleteRecipeHandler(w http.ResponseWriter, r *http.Request) {
	// Разрешаем CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	//Разрешаем только DELETE
	if r.Method != http.MethodDelete {
		http.Error(w,
			fmt.Sprintf("%v %v", InvMethodErr, r.Method),
			http.StatusMethodNotAllowed,
		)

		logger.LogPrint(
			fmt.Sprintf("%v %v", InvMethodErr, r.Method),
			models.Error,
		)
		return
	}

	//Сплит адреса по сепаратору "/"
	urlParts := strings.Split(r.URL.Path, "/") //[api deleteRecipe X] X --> id of recipe
	if len(urlParts) < 3 {
		http.Error(w, InvURL_Err.Error(), http.StatusBadRequest)

		logger.LogPrint(
			InvURL_Err.Error(),
			models.Error,
		)
		return
	}

	//Получаем ID рецепта
	recipeID, err := strconv.Atoi(urlParts[3])
	if err != nil {
		http.Error(w, GetID_Err.Error(), http.StatusBadRequest)

		logger.LogPrint(
			GetID_Err.Error(),
			models.Error,
		)
		return
	}

	//Отправляем запрос на удаление в БД
	err = database.DeleteRecipe(recipeID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)

		logger.LogPrint(
			err.Error(),
			models.Error,
		)
		return
	}

	client, err := getClientIP(r) //Получение ip клиента
	if err != nil {
		colors.SetColor(colors.Text_Red)
		logger.LogPrint(ClientIPErr.Error(), models.Error)
		colors.ResetColor()
	}

	colors.SetColor(colors.Text_Purple)
	logger.LogPrint(
		fmt.Sprintf("Клиент [%s] удалил данные по id [%d] на сервере!", client, recipeID),
		models.Info,
	)
	colors.ResetColor()

	//Возвращаем статус 200 ОК и сообщение
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Удаление выполнено успешно!"))
}

// Обработчик для localhost:8080/api/updateRecipe/
func updateRecipeHandler(w http.ResponseWriter, r *http.Request) {
	// Разрешаем CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	//Разрешаем только PUT
	if r.Method != http.MethodPut {
		http.Error(w,
			fmt.Sprintf("%v %v", InvMethodErr, r.Method),
			http.StatusMethodNotAllowed)

		logger.LogPrint(
			fmt.Sprintf("%v %v", InvMethodErr, r.Method),
			models.Error,
		)
		return
	}

	//Десериализуем тело запроса
	data, err := Deserialization(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		logger.LogPrint(
			err.Error(),
			models.Error,
		)
		return
	}

	//Отправляем запрос в БД
	err = database.UpdateRecipe(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)

		logger.LogPrint(
			err.Error(),
			models.Error,
		)
		return
	}

	client, err := getClientIP(r) //Получаем ip клиента
	if err != nil {
		colors.SetColor(colors.Text_Red)
		logger.LogPrint(ClientIPErr.Error(), models.Error)
		colors.ResetColor()
	}

	colors.SetColor(colors.Text_Purple)
	logger.LogPrint(
		fmt.Sprintf("Клиент [%s] обновил данные по id [%d] на сервере!", client, data.ID),
		models.Info,
	)
	colors.ResetColor()

	//Возвращаем статус 200 ОК и сообщение
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Редактирование выполнено успешно!"))
}
