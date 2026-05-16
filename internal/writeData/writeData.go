package writeData

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"webtest/internal/logger"
	"webtest/internal/models"
	"webtest/pkg/colors"
)

var (
	ErrorDeserialization = errors.New("Ошибка десериализации")
	ErrorReadBody        = errors.New("Ошибка чтения тела запроса!")
)

func Deserialization(r *http.Request) (*models.Eat, error) {
	var eat_data models.Eat

	//Читаем тело запроса
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, ErrorReadBody
	}
	colors.SetColor(colors.Text_Green)
	logger.LogPrint("Чтение тела запроса прошло успешно!", models.Info)
	colors.ResetColor()

	//Процесс десериализации
	err = json.Unmarshal(body, &eat_data)
	if err != nil {
		return nil, ErrorDeserialization
	}
	colors.SetColor(colors.Text_Green)
	logger.LogPrint("Десериализация прошла успешно!", models.Info)
	colors.ResetColor()

	colors.SetColor(colors.Text_Green)
	logger.LogPrint(fmt.Sprintf("Result: %v", eat_data), models.Info)
	colors.ResetColor()

	return &eat_data, nil
}
