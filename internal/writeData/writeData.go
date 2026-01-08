package writeData

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"webtest/internal/models"
	"webtest/pkg/colors"
)

func Deserialization(r *http.Request) (*models.Eat, error) {
	var eat_data models.Eat

	//Читаем тело запроса
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, errors.New("Ошибка чтения тела запроса!")
	}
	colors.SetColor(colors.Text_Green)
	log.Println("Чтение тела запроса прошло успешно!")
	colors.ResetColor()

	//Процесс десериализации
	err = json.Unmarshal(body, &eat_data)
	if err != nil {
		return nil, errors.New("Ошибка десериализации!")
	}
	colors.SetColor(colors.Text_Green)
	log.Println("Десериализация прошла успешно!")
	colors.ResetColor()

	colors.SetColor(colors.Text_Green)
	fmt.Println("\nResult: \n", eat_data)
	colors.ResetColor()

	return &eat_data, nil
}
